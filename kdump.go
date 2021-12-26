package main

import (
	"fmt"
	"github.com/thoas/go-funk"
	"kdump/internal/fileutil"
	"kdump/internal/kubectl"
	"kdump/internal/stringutil"
	"log"
	"os"
	"strconv"
	"strings"
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
		return ApiResource{
			name:       getMapStrValOrEmpty(in, "NAME"),
			shortNames: csvStr2arr(getMapStrValOrEmpty(in, "SHORTNAMES")),
			namespaced: str2bool(getMapStrValOrEmpty(in, "NAMESPACED")),
			kind:       getMapStrValOrEmpty(in, "KIND"),
			verbs:      wierdKubectlArray2arr(getMapStrValOrEmpty(in, "VERBS")),
		}
	}).([]ApiResource)

	for _, resource := range apiResources {
		log.Printf("resource: %+v \n", resource)
	}
}

func getMapStrValOrEmpty(dict map[string]string, key string) string {
	if val, ok := dict[key]; ok {
		return val
	} else {
		return ""
	}
}

func str2bool(str string) bool {
	if val, err := strconv.ParseBool(str); err == nil {
		return val
	} else {
		return true
	}
}

func csvStr2arrSep(str string, sep string) []string {
	return stringutil.MapStrArray(strings.Split(str, sep), func(in string) string {
		return strings.TrimSpace(in)
	})
}

func csvStr2arr(str string) []string {
	return csvStr2arrSep(str, ",")
}

func wierdKubectlArray2arr(strIn string) []string {
	return csvStr2arrSep(strIn[1:(len(strIn)-1)], " ")
}
