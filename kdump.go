package main

import (
	"fmt"
	"kdump/internal/fileutil"
	"kdump/internal/kubectl"
	"log"
)

func main() {
	currentContext := kubectl.CurrentContext()
	dumpCurrentContext("test/" + currentContext)
}

func dumpCurrentContext(outputDir string) {

	fileutil.PanicIfCantDelete(outputDir, fmt.Sprintf("removal of outputdir '%s' failed", outputDir))
	fileutil.PanicIfExists(outputDir, fmt.Sprintf("output folder '%s' already exists!", outputDir), fmt.Sprintf("output folder '%s' inaccessible!", outputDir))
	fileutil.CreateFolderOrPanic(outputDir, fmt.Sprintf("could not create folder '%s'", outputDir))

	log.Printf("Downloading all resources from current context to dir %s ...\n", outputDir)

	namespaces := kubectl.Namespaces()
	log.Printf("Namespaces: %v ...\n", namespaces)

	apiResourceTypes := kubectl.ApiResourceTypes()

	log.Printf("\n")

	log.Printf("global resources: \n")
	for _, resource := range apiResourceTypes.Accessible.Global {
		log.Printf("  resource: %+v \n", resource)
	}
	log.Printf("\n")

	log.Printf("namespaced resources: \n")
	for _, resource := range apiResourceTypes.Accessible.Namespaced {
		log.Printf("  resource: %+v \n", resource)
	}
	log.Printf("\n")

	for _, namespace := range namespaces {

		fileutil.CreateFolderOrPanic(outputDir+"/"+namespace, "could not create output dir for namespace "+namespace)

		for _, namespaceResourceType := range apiResourceTypes.Accessible.Namespaced {
			items := kubectl.ListNamespacedResourcesOfType(namespace, namespaceResourceType.Name)
			if len(items) == 0 {
				continue
			}
			itemDir := outputDir + "/" + namespace + "/" + namespaceResourceType.Name
			fileutil.CreateFolderOrPanic(itemDir, "could not create output dir for namespace resource "+namespaceResourceType.Name)
			for _, item := range items {
				log.Printf("Storing item %v in folder %v", item, itemDir)
			}
		}
	}

}
