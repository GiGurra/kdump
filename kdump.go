package main

import (
	"fmt"
	"github.com/gigurra/kdump/config"
	"github.com/gigurra/kdump/internal/coll"
	"github.com/gigurra/kdump/internal/fileutil"
	"github.com/gigurra/kdump/internal/kubectl"
	"log"
	"strconv"
	"syscall"
)

func main() {
	// TODO: Parse real config
	appConfig := config.GetDefaultAppConfig()
	dumpCurrentContext(appConfig)
}

func dumpCurrentContext(appConfig config.AppConfig) {

	log.Printf("Downloading all resources from current context")

	testArray := []int{1, 1, 2, 3}

	mapped := coll.GroupBy(testArray, func(in int) string {
		return strconv.Itoa(in + 1)
	})

	log.Printf("mapped: %+v", mapped)

	syscall.Exit(0)

	outputDirRoot := appConfig.OutputDir
	apiResourceTypes := kubectl.ApiResourceTypes()
	resourcesToDownload := appConfig.FilterIncludedResources(apiResourceTypes.Accessible.All)
	everything := kubectl.DownloadEverything(resourcesToDownload)

	log.Printf("Parsing %d bytes...\n", len(everything))

	k8sResources := kubectl.ParseK8sYaml(everything)

	log.Printf("Deleting old data in '%s'...\n", outputDirRoot)

	fileutil.Delete(outputDirRoot, fmt.Sprintf("removal of outputdir '%s' failed", outputDirRoot))
	fileutil.CreateFolder(outputDirRoot, fmt.Sprintf("could not create folder '%s'", outputDirRoot))

	log.Printf("Storing resources in '%s'...\n", outputDirRoot)
	for _, resource := range k8sResources {
		filename := fileutil.SanitizePath(resource.MetaData.Name) + "." + fileutil.SanitizePath(resource.QualifiedTypeName) + ".yaml"
		if resource.IsSecret() {
			log.Printf("Ignoring secret storage (not yet implemented) for %s/%s: ", resource.MetaData.Namespace, resource.MetaData.Name)
		} else if resource.IsNamespaced() {
			nsOutputDir := outputDirRoot + "/" + resource.MetaData.Namespace
			fileutil.CreateFolderIfMissing(nsOutputDir, "could not create output dir: "+nsOutputDir)
			fileutil.String2File(nsOutputDir+"/"+filename, resource.SourceYaml)
		} else {
			fileutil.String2File(outputDirRoot+"/"+filename, resource.SourceYaml)
		}
	}
}
