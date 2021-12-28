package main

import (
	"fmt"
	"github.com/gigurra/kdump/config"
	"github.com/gigurra/kdump/internal/fileutil"
	"github.com/gigurra/kdump/internal/k8s"
	"github.com/gigurra/kdump/internal/k8s/kubectl"
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
