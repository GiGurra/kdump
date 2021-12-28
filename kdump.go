package main

import (
	"fmt"
	"github.com/gigurra/kdump/config"
	"github.com/gigurra/kdump/internal/fileutil"
	"github.com/gigurra/kdump/internal/k8s"
	"github.com/gigurra/kdump/internal/k8s/kubectl"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
)

func validateNonEmpty(str string) string {
	if len(strings.TrimSpace(str)) == 0 {
		log.Fatal("empty string is not allowed")
	}
	return str
}

func main() {

	app := cli.NewApp()

	app.HideHelpCommand = true
	app.Usage = "Dump all kubernetes resources as yaml files to a dir"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "output-dir",
			Aliases:  []string{"o"},
			Usage:    "output directory to create",
			Required: true,
		},
		&cli.BoolFlag{
			Name:    "delete-previous-dir",
			Aliases: []string{"d"},
			Usage:   "if to delete previous output directory",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:  "include-secrets",
			Usage: "if to include secrets",
			Value: false,
		},
	}

	app.Action = func(c *cli.Context) error {
		appConfig := config.GetDefaultAppConfig()
		appConfig.OutputDir = validateNonEmpty(c.String("output-dir"))
		appConfig.IncludeSecrets = c.Bool("include-secrets")
		appConfig.DeletePrevDir = c.Bool("delete-previous-dir")
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
	outputDirRoot := ensureRootOutputDir(appConfig)

	log.Printf("Downloading all resources from current context")

	apiResourceTypes := kubectl.ApiResourceTypes()
	resourcesToDownload := appConfig.FilterIncludedResources(apiResourceTypes.Accessible.All)
	everything := kubectl.DownloadEverything(resourcesToDownload)

	log.Printf("Parsing %d bytes...\n", len(everything))

	k8sResources := k8s.ParseResources(everything)
	k8sResourcesByNamespace := k8s.GroupByNamespace(k8sResources)

	log.Printf("Storing resources in '%s'...\n", outputDirRoot)
	for namespace, resources := range k8sResourcesByNamespace {
		outDir := outputDirRoot
		if namespace != "" {
			outDir = outputDirRoot + "/" + namespace
			fileutil.CreateFolderIfNotExists(outDir, "could not create output dir: "+outDir)
		}
		for _, resource := range resources {
			filename := fileutil.SanitizePath(resource.MetaData.Name) + "." + fileutil.SanitizePath(resource.QualifiedTypeName) + ".yaml"
			if resource.IsSecret() {
				log.Printf("Ignoring secret storage (not yet implemented) for %s/%s: ", resource.MetaData.Namespace, resource.MetaData.Name)
			} else {
				fileutil.String2File(outDir+"/"+filename, resource.SourceYaml)
			}
		}
	}
}

func ensureRootOutputDir(appConfig config.AppConfig) string {

	out := appConfig.OutputDir

	if appConfig.DeletePrevDir {
		fileutil.Delete(out, fmt.Sprintf("removal of outputdir '%s' failed", out))
	}

	if fileutil.Exists(out, fmt.Sprintf("checking outputdir '%s' failed", out)) {
		log.Fatal("Bailing! output-dir already exists: " + out)
	}

	fileutil.CreateFolderIfNotExists(out, fmt.Sprintf("could not create folder '%s'", out))

	return out
}
