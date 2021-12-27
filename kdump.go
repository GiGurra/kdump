package main

import (
	"fmt"
	"kdump/config"
	"kdump/internal/fileutil"
	"kdump/internal/kubectl"
	"log"
	"syscall"
)

func main() {
	// TODO: Parse real config
	appConfig := config.GetDefaultAppConfig()
	dumpCurrentContext(appConfig)
}

func dumpCurrentContext(appConfig config.AppConfig) {

	currentK8sContext := kubectl.CurrentContextOrPanic()
	outputDir := appConfig.GetOutDir(currentK8sContext)

	fileutil.PanicIfCantDelete(outputDir, fmt.Sprintf("removal of outputdir '%s' failed", outputDir))
	fileutil.PanicIfExists(outputDir, fmt.Sprintf("output folder '%s' already exists!", outputDir), fmt.Sprintf("output folder '%s' inaccessible!", outputDir))
	fileutil.CreateFolderOrPanic(outputDir, fmt.Sprintf("could not create folder '%s'", outputDir))

	log.Printf("Downloading all resources from current context to dir %s ...\n", outputDir)

	namespaces := kubectl.NamespacesOrPanic()
	apiResourceTypes := kubectl.ApiResourceTypesOrPanic()

	everything := kubectl.DownloadEverythingInNamespaceOrPanic(apiResourceTypes.Accessible.Namespaced[0:1], "default")
	kubectl.ParseK8sYamlOrPanic(everything)
	fileutil.String2FileOrPanic(outputDir+"/default.yaml", everything)

	syscall.Exit(1)

	for _, namespace := range namespaces {

		fileutil.CreateFolderOrPanic(outputDir+"/"+namespace, "could not create output dir for namespace "+namespace)

		for _, namespaceResourceType := range apiResourceTypes.Accessible.Namespaced {
			if appConfig.IncludeResource(namespaceResourceType) {
				dumpNamespacedResourcesOrPanic(outputDir, namespace, namespaceResourceType.QualifiedName)
			}
		}
	}
}

func dumpNamespacedResourcesOrPanic(
	outputDir string,
	namespace string,
	resourceTypeName string,
) {
	resourceNames := kubectl.ListNamespacedResourcesOfTypeOrPanic(namespace, resourceTypeName)
	if len(resourceNames) > 0 {
		if resourceTypeName == "secrets" {
			dumpSecretsOrPanic(outputDir, namespace, resourceTypeName, resourceNames)
		} else {
			dumpRegularNamespacedResourcesOrPanic(outputDir, namespace, resourceTypeName, resourceNames)
		}
	}
}

func dumpRegularGlobalResources(
	outputDir string,
	resourceTypeName string,
	resourceNames []string,
) {
	itemDir := outputDir + "/" + resourceTypeName
	fileutil.CreateFolderOrPanic(itemDir, "could not create output dir for global resource "+resourceTypeName)
	for _, item := range resourceNames {
		resource := kubectl.DownloadGlobalResourceOrPanic(resourceTypeName, item, "yaml")
		log.Printf("Storing item %v in folder %v", item, itemDir)
		fileutil.String2FileOrPanic(itemDir+"/"+fileutil.ReplaceInvalidChars(item)+".yaml", resource)
	}
}

func dumpRegularNamespacedResourcesOrPanic(
	outputDir string,
	namespace string,
	resourceTypeName string,
	resourceNames []string,
) {
	itemDir := outputDir + "/" + namespace + "/" + resourceTypeName
	fileutil.CreateFolderOrPanic(itemDir, "could not create output dir for namespace resource "+resourceTypeName)
	for _, item := range resourceNames {
		resource := kubectl.DownloadNamespacedResourceOrPanic(namespace, resourceTypeName, item, "yaml")
		log.Printf("Storing item %v in folder %v", item, itemDir)
		fileutil.String2FileOrPanic(itemDir+"/"+fileutil.ReplaceInvalidChars(item)+".yaml", resource)
	}
}

func dumpSecretsOrPanic(
	outputDir string,
	namespace string,
	resourceTypeName string,
	resourceNames []string,
) {
	log.Printf("ignoring storage of secrets. Not yet implemented \n")
}
