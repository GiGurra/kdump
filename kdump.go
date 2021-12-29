package main

import (
	"fmt"
	"github.com/gigurra/kdump/config"
	"github.com/gigurra/kdump/internal/cliUtil"
	"github.com/gigurra/kdump/internal/crypt"
	"github.com/gigurra/kdump/internal/fileutil"
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

	log.Printf("Checking output dir..")
	rootOutputDir := ensureRootOutputDir(appConfig)

	log.Printf("Downloading all resources from current context")

	apiResourceTypes := kubectl.ApiResourceTypes()
	resourcesToDownload := appConfig.FilterIncludedResources(apiResourceTypes.Accessible.All)
	everything := kubectl.DownloadEverything(resourcesToDownload)

	log.Printf("Parsing %d bytes...\n", len(everything))

	k8sResources := k8s.ParseResources(everything)
	k8sResourcesByNamespace := k8s.GroupByNamespace(k8sResources)

	log.Printf("Storing resources in '%s'...\n", rootOutputDir)
	for namespace, resources := range k8sResourcesByNamespace {
		outDir := rootOutputDir
		if namespace != "" {
			outDir = rootOutputDir + "/" + namespace
			fileutil.CreateFolderIfNotExists(outDir, "could not create output dir: "+outDir)
		}
		for _, resource := range resources {
			filename := fileutil.SanitizePath(resource.MetaData.Name) + "." + fileutil.SanitizePath(resource.QualifiedTypeName) + ".yaml"
			if resource.IsSecret() {
				fileutil.String2File(outDir+"/"+filename+".aes", crypt.Encrypt(resource.SourceYaml, appConfig.SecretsEncryptKey))
			} else {
				fileutil.String2File(outDir+"/"+filename, resource.SourceYaml)
			}
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
