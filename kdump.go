package main

import (
	"fmt"
	"github.com/gigurra/go-util/cliUtil"
	"github.com/gigurra/go-util/crypt"
	"github.com/gigurra/go-util/fileutil"
	"github.com/gigurra/go-util/shell"
	"github.com/gigurra/kdump/config"
	"github.com/gigurra/kdump/internal/k8s"
	"github.com/gigurra/kdump/internal/k8s/kubectl"
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

	apiResourceTypes := kubectl.ApiResourceTypes()
	resourcesToDownload := appConfig.FilterIncludedResources(apiResourceTypes.Accessible.All)
	everything := kubectl.DownloadEverything(resourcesToDownload)

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
				trgFilePath := outDir + "/" + filename + ".aes"
				log.Printf("Processing (%d / %d) %s", iResource+1, totalResourceCount, trgFilePath)
				fileutil.String2File(trgFilePath, crypt.Encrypt(resource.SourceYaml, appConfig.SecretsEncryptKey))
			} else {
				filePath := outDir + "/" + filename
				log.Printf("Processing (%d / %d) %s", iResource+1, totalResourceCount, filePath)
				fileutil.String2File(filePath, resource.SourceYaml)
				neatifiedYaml := shell.RunCommand("kubectl", "neat", "-f", filePath)
				os.Remove(filePath)
				fileutil.String2File(filePath, neatifiedYaml)
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
