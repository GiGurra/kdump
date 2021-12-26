package main

import (
	"fmt"
	"github.com/thoas/go-funk"
	"kdump/internal/fileutil"
	"kdump/internal/kubectl"
	"kdump/internal/stringutil"
	"log"
	"os"
)

func main() {

	currentContext := kubectl.CurrentContext()
	currentNamespace := kubectl.CurrentNamespace()
	namespaces := kubectl.Namespaces()

	fmt.Println(currentContext)
	fmt.Println(currentNamespace)
	fmt.Println(namespaces)

	dumpCurrentContext("cx", true)

}

func dumpCurrentContext(outputDir string, allowOverwrite bool) {

	if allowOverwrite {
		err := os.RemoveAll(outputDir)
		if err != nil {
			panic(fmt.Sprintf("removal of outputdir '%s' failed with err %v", outputDir, err))
		}
	}

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

	accessibleApiResources := funk.Filter(allApiResources, isAccessible).([]ApiResource)

	for _, resource := range accessibleApiResources {
		log.Printf("resource: %+v \n", resource)
	}
}

type ApiResource struct {
	name       string
	shortNames []string
	namespaced bool
	kind       string
	verbs      []string
}

func isAccessible(r ApiResource) bool {
	return funk.ContainsString(r.verbs, "get")
}
