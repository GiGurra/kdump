package main

import (
	"fmt"
	"github.com/thoas/go-funk"
	"kdump/config"
	"kdump/internal/fileutil"
	"kdump/internal/kubectl"
	"log"
)

func main() {
	// TODO: Parse real config
	appConfig := config.GetDefaultAppConfig()
	dumpCurrentContext(appConfig)
}

func dumpCurrentContext(appConfig config.AppConfig) {

	currentK8sContext := kubectl.CurrentContextOrPanic()
	outputDir := appConfig.GetOutDir(currentK8sContext)

	log.Printf("Downloading all resources from current context to dir %s ...\n", outputDir)

	namespaces := kubectl.NamespacesOrPanic()
	apiResourceTypes := kubectl.ApiResourceTypesOrPanic()
	resourcesToDownload := funk.Filter(apiResourceTypes.Accessible.All, func(r *kubectl.ApiResourceType) bool {
		return appConfig.IncludeResource(r)
	}).([]*kubectl.ApiResourceType)
	everything := kubectl.DownloadEverythingOrPanic(resourcesToDownload)

	log.Printf("Parsing %d bytes...\n", len(everything))

	parsed := kubectl.ParseK8sYamlOrPanic(everything)

	log.Printf("Deleting old folders...\n")

	fileutil.PanicIfCantDelete(outputDir, fmt.Sprintf("removal of outputdir '%s' failed", outputDir))
	fileutil.PanicIfExists(outputDir, fmt.Sprintf("output folder '%s' already exists!", outputDir), fmt.Sprintf("output folder '%s' inaccessible!", outputDir))
	fileutil.CreateFolderOrPanic(outputDir, fmt.Sprintf("could not create folder '%s'", outputDir))

	log.Printf("Storing resources...\n")
	for _, ns := range namespaces {
		nsOutputDir := outputDir + "/" + ns
		fileutil.CreateFolderOrPanic(nsOutputDir, "could not create folder: "+nsOutputDir)
	}

	for _, resource := range parsed {
		filename := fileutil.SanitizePath(resource.MetaData.Name) + "." + fileutil.SanitizePath(resource.QualifiedTypeName) + ".yaml"
		if resource.IsSecret() {
			log.Printf("Ignoring secret storage (not yet implemented) for %s/%s: ", resource.MetaData.Namespace, resource.MetaData.Name)
		} else if resource.IsNamespaced() {
			nsOutputDir := outputDir + "/" + resource.MetaData.Namespace
			fileutil.String2FileOrPanic(nsOutputDir+"/"+filename, resource.SourceYaml)
		} else {
			fileutil.String2FileOrPanic(outputDir+"/"+filename, resource.SourceYaml)
		}
	}
}
