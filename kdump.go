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
	namespaces := kubectl.NamespacesOrPanic()
	apiResourceTypes := kubectl.ApiResourceTypesOrPanic()
	resourcesToDownload := appConfig.FilterIncludeResources(apiResourceTypes.Accessible.All)
	everything := kubectl.DownloadEverythingOrPanic(resourcesToDownload)

	log.Printf("Parsing %d bytes...\n", len(everything))

	parsed := kubectl.ParseK8sYamlOrPanic(everything)

	log.Printf("Deleting old data in '%s'...\n", outputDirRoot)

	fileutil.PanicIfCantDelete(outputDirRoot, fmt.Sprintf("removal of outputdir '%s' failed", outputDirRoot))
	fileutil.PanicIfExists(outputDirRoot, fmt.Sprintf("output folder '%s' already exists!", outputDirRoot), fmt.Sprintf("output folder '%s' inaccessible!", outputDirRoot))
	fileutil.CreateFolderOrPanic(outputDirRoot, fmt.Sprintf("could not create folder '%s'", outputDirRoot))

	log.Printf("Storing resources in '%s'...\n", outputDirRoot)
	for _, ns := range namespaces {
		nsOutputDir := outputDirRoot + "/" + ns
		fileutil.CreateFolderOrPanic(nsOutputDir, "could not create folder: "+nsOutputDir)
	}

	for _, resource := range parsed {
		filename := fileutil.SanitizePath(resource.MetaData.Name) + "." + fileutil.SanitizePath(resource.QualifiedTypeName) + ".yaml"
		if resource.IsSecret() {
			log.Printf("Ignoring secret storage (not yet implemented) for %s/%s: ", resource.MetaData.Namespace, resource.MetaData.Name)
		} else if resource.IsNamespaced() {
			nsOutputDir := outputDirRoot + "/" + resource.MetaData.Namespace
			fileutil.String2FileOrPanic(nsOutputDir+"/"+filename, resource.SourceYaml)
		} else {
			fileutil.String2FileOrPanic(outputDirRoot+"/"+filename, resource.SourceYaml)
		}
	}
}
