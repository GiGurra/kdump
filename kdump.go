package main

import (
	"fmt"
	"github.com/gigurra/kdump/config"
	"github.com/gigurra/kdump/internal/fileutil"
	"github.com/gigurra/kdump/internal/kubectl"
	"log"
)

func main() {
	// TODO: Parse real config
	appConfig := config.GetDefaultAppConfig()
	dumpCurrentContext(appConfig)
}

func dumpCurrentContext(appConfig config.AppConfig) {

	log.Printf("Downloading all resources from current context")

	outputDirRoot := appConfig.OutputDir
	apiResourceTypes := kubectl.ApiResourceTypes()
	resourcesToDownload := appConfig.FilterIncludedResources(apiResourceTypes.Accessible.All)
	everything := kubectl.DownloadEverything(resourcesToDownload)

	log.Printf("Parsing %d bytes...\n", len(everything))

	parsed := kubectl.ParseK8sYaml(everything)

	log.Printf("Deleting old data in '%s'...\n", outputDirRoot)

	fileutil.Delete(outputDirRoot, fmt.Sprintf("removal of outputdir '%s' failed", outputDirRoot))
	fileutil.CreateFolder(outputDirRoot, fmt.Sprintf("could not create folder '%s'", outputDirRoot))

	log.Printf("Storing resources in '%s'...\n", outputDirRoot)
	for _, resource := range parsed {
		filename := fileutil.SanitizePath(resource.MetaData.Name) + "." + fileutil.SanitizePath(resource.QualifiedTypeName) + ".yaml"
		if resource.IsSecret() {
			log.Printf("Ignoring secret storage (not yet implemented) for %s/%s: ", resource.MetaData.Namespace, resource.MetaData.Name)
		} else if resource.IsNamespaced() {
			nsOutputDir := outputDirRoot + "/" + resource.MetaData.Namespace
			if !fileutil.Exists(nsOutputDir, "could not determine if outputdir exists: "+nsOutputDir) {
				fileutil.CreateFolder(nsOutputDir, "could not create outputdir: "+nsOutputDir)
			}
			fileutil.String2File(nsOutputDir+"/"+filename, resource.SourceYaml)
		} else {
			fileutil.String2File(outputDirRoot+"/"+filename, resource.SourceYaml)
		}
	}
}
