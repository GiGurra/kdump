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

func overrideStrIfNonEmpty(prev *string, override string) {
	if len(strings.TrimSpace(override)) > 0 {
		*prev = override
	}
}

func overrideBool(prev *bool, override bool) {
	*prev = override
}

func main() {

	app := cli.NewApp()

	app.HideHelpCommand = true
	app.Usage = "Dump all kubernetes resources as yaml files to a dir"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "output-dir",
			Aliases: []string{"o"},
			Usage:   "output directory to create",
			Value:   "test",
		},
		&cli.BoolFlag{
			Name:  "include-secrets",
			Usage: "if to include secrets",
			Value: false,
		},
	}

	app.Action = func(c *cli.Context) error {
		appConfig := config.GetDefaultAppConfig()
		overrideStrIfNonEmpty(&appConfig.OutputDir, c.String("output-dir"))
		overrideBool(&appConfig.IncludeSecrets, c.Bool("include-secrets"))
		//log.Printf("Config: \n%s\n", util.OrPanic(yaml.Marshal(appConfig)))
		dumpCurrentContext(appConfig)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func dumpCurrentContext(appConfig config.AppConfig) {

	log.Printf("Downloading all resources from current context")

	outputDirRoot := appConfig.OutputDir
	apiResourceTypes := kubectl.ApiResourceTypes()
	resourcesToDownload := appConfig.FilterIncludedResources(apiResourceTypes.Accessible.All)
	everything := kubectl.DownloadEverything(resourcesToDownload)

	log.Printf("Parsing %d bytes...\n", len(everything))

	k8sResources := k8s.ParseResources(everything)
	k8sResourcesByNamespace := k8s.GroupByNamespace(k8sResources)

	log.Printf("Deleting old data in '%s'...\n", outputDirRoot)

	fileutil.Delete(outputDirRoot, fmt.Sprintf("removal of outputdir '%s' failed", outputDirRoot))
	fileutil.CreateFolder(outputDirRoot, fmt.Sprintf("could not create folder '%s'", outputDirRoot))

	log.Printf("Storing resources in '%s'...\n", outputDirRoot)
	for namespace, resources := range k8sResourcesByNamespace {
		outDir := outputDirRoot
		if namespace != "" {
			outDir = outputDirRoot + "/" + namespace
			fileutil.CreateFolder(outDir, "could not create output dir: "+outDir)
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
