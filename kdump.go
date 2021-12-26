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

	type ApiResource struct {
		name       string
		shortNames []string
		namespaced bool
		kind       string
		verbs      []string
	}

	apiRsrcsStr := kubectl.ApiResources()
	_ /* schema */, apiResourcesRaw := stringutil.ParseStdOutTable(apiRsrcsStr)

	apiResources := funk.Map(apiResourcesRaw, func(in map[string]string) ApiResource {
		return ApiResource{}
	}).([]ApiResource)

	log.Printf("apiResources: %+v ...\n", apiResources)
}
