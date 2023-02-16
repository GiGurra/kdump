package main

import (
	"fmt"
	"github.com/gigurra/go-util/cliUtil"
	"github.com/gigurra/go-util/crypt"
	"github.com/gigurra/go-util/fileutil"
	"github.com/gigurra/kdump/config"
	"github.com/gigurra/kdump/internal/k8s"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {

	app := cli.NewApp()

	app.HideHelpCommand = true
	app.Usage = "Dump all kubernetes resources as yaml files to a dir"
	app.Flags = cliUtil.FindAllFlags(&config.CliFlags)
	app.Version = Version
	app.Action = func(c *cli.Context) error {
		appConfig := config.GetDefaultAppConfig()
		appConfig.OutputDir = c.String(config.CliFlags.OutputDir.Name)
		appConfig.DeletePrevDir = c.Bool(config.CliFlags.DeletePrevDir.Name)
		appConfig.SecretsEncryptKey = c.String(config.CliFlags.EncryptKey.Name)
		appConfig.Validate()
		dumpCurrentContext(appConfig)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func dumpCurrentContext(appConfig config.AppConfig) {

	log.Printf("Running kdump version " + Version)
	log.Printf("Checking output dir..")
	rootOutputDir := ensureRootOutputDir(appConfig)

	log.Printf("Downloading all resources from current context")

	apiResourceTypes := k8s.ApiResourceTypes()
	resourcesToDownload := appConfig.FilterIncludedResources(apiResourceTypes.Accessible.All)

	log.Printf("Downloading all resources of %d types", len(resourcesToDownload))
	everythingRaw := k8s.DownloadEverything(resourcesToDownload)

	log.Printf("Running kubectl neat on everything.. (this takes about 3-4x the download time)")
	everything := k8s.PipeToCommand(everythingRaw, "kubectl", "neat")

	log.Printf("Parsing %d bytes...\n", len(everything))

	k8sResources := k8s.ParseResources(everything)
	k8sResourcesByNamespace := k8s.GroupByNamespace(k8sResources)

	totalResourceCount := len(k8sResources)
	iResource := 0

	log.Printf("Storing %d resources in '%s'...\n", totalResourceCount, rootOutputDir)
	for namespace, resources := range k8sResourcesByNamespace {
		outDir := rootOutputDir
		if namespace != "" {
			outDir = rootOutputDir + "/" + namespace
			fileutil.CreateFolderIfNotExists(outDir, "could not create output dir: "+outDir)
		}
		for _, resource := range resources {
			name := fileutil.SanitizePath(resource.MetaData.Name)
			typ := fileutil.SanitizePath(resource.QualifiedTypeName)
			filename := name + "." + typ + ".yaml"
			if resource.IsSecret() {
				filePath := outDir + "/" + filename + ".aes"
				fileutil.String2File(filePath, crypt.Encrypt(resource.SourceYaml, appConfig.SecretsEncryptKey))
			} else {
				filePath := outDir + "/" + filename
				fileutil.String2File(filePath, resource.SourceYaml)
			}
			iResource++
		}
	}
}

func ensureRootOutputDir(appConfig config.AppConfig) string {

	out := appConfig.OutputDir

	if appConfig.DeletePrevDir {
		fileutil.DeleteIfExists(out, fmt.Sprintf("removal of outputdir '%s' failed", out))
	}

	if fileutil.Exists(out, fmt.Sprintf("checking outputdir '%s' failed", out)) {
		log.Fatal("Bailing! output-dir already exists: " + out)
	}

	fileutil.CreateFolderIfNotExists(out, fmt.Sprintf("could not create folder '%s'", out))

	return out
}
