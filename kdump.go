package main

import (
	"fmt"
	"github.com/thoas/go-funk"
	"kdump/internal/fileutil"
	"kdump/internal/kubectl"
	"kdump/internal/stringutil"
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

	apiRsrcsStr := kubectl.ApiResources()
	_ /* schema */, apiResourcesRaw := stringutil.ParseStdOutTable(apiRsrcsStr)

	allApiResources := funk.Map(apiResourcesRaw, func(in map[string]string) ApiResource {
		return ApiResource{
			name:       stringutil.MapStrValOrElse(in, "NAME", ""),
			shortNames: stringutil.CsvStr2arr(stringutil.MapStrValOrElse(in, "SHORTNAMES", "")),
			namespaced: stringutil.Str2boolOrElse(stringutil.MapStrValOrElse(in, "NAMESPACED", ""), false),
			kind:       stringutil.MapStrValOrElse(in, "KIND", ""),
			verbs:      stringutil.WierdKubectlArray2arr(stringutil.MapStrValOrElse(in, "VERBS", "")),
		}
	}).([]ApiResource)

	accessibleApiResources := funk.Filter(allApiResources, func(r ApiResource) bool { return funk.ContainsString(r.verbs, "get") }).([]ApiResource)
	globalResources := funk.Filter(accessibleApiResources, func(r ApiResource) bool { return !r.namespaced }).([]ApiResource)
	namespacedResources := funk.Filter(accessibleApiResources, func(r ApiResource) bool { return r.namespaced }).([]ApiResource)
	log.Printf("\n")

	log.Printf("global resources: \n")
	for _, resource := range globalResources {
		log.Printf("  resource: %+v \n", resource)
	}
	log.Printf("\n")

	log.Printf("namespaced resources: \n")
	for _, resource := range namespacedResources {
		log.Printf("  resource: %+v \n", resource)
	}
	log.Printf("\n")

}

type ApiResource struct {
	name       string
	shortNames []string
	namespaced bool
	kind       string
	verbs      []string
}
