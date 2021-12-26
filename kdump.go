package main

import (
	"fmt"
	"kdump/internal/fileutil"
	"kdump/internal/kubectl"
	"log"
)

func main() {

	currentContext := kubectl.CurrentContext()
	currentNamespace := kubectl.CurrentNamespace()
	namespaces := kubectl.Namespaces()

	fmt.Println(currentContext)
	fmt.Println(currentNamespace)
	fmt.Println(namespaces)

	dumpCurrentContext("cx")

}

func dumpCurrentContext(outputDir string) {

	fileutil.PanicIfExists(outputDir, fmt.Sprintf("output folder '%s' already exists!", outputDir), fmt.Sprintf("output folder '%s' inaccessible!", outputDir))

	log.Printf("Downloading all resources from current context to dir %s ...\n", outputDir)
}
