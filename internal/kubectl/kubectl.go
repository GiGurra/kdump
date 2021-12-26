package kubectl

import (
	"github.com/thoas/go-funk"
	"kdump/internal/shell"
	"kdump/internal/stringutil"
	"os/exec"
	"strings"
)

func Namespaces() []string {
	return stringutil.MapStrArray(stringutil.SplitLines(runCommand("get", "namespaces", "-o", "name")), removeK8sResourcePrefix)
}

func CurrentNamespace() string {
	return runCommand("config", "view", "--minify", "--output", "jsonpath={..namespace}")
}

func CurrentContext() string {
	return runCommand("config", "current-context")
}

func ListNamespacedResourcesOfType(namespace string, resourceType string) []string {
	// kubectl -n ${namespace} get ` + resourceName + " -o name"
	rawString := runCommand("-n", namespace, "get", resourceType, "-o", "name")
	return stringutil.RemoveEmptyLines(stringutil.SplitLines(rawString))
}

type ApiResourceType struct {
	Name       string
	ShortNames []string
	Namespaced bool
	Kind       string
	Verbs      []string
}

type ApiResourceTypesResponse struct {
	All        []ApiResourceType
	Accessible ApiResourceTypesAccessible
}

type ApiResourceTypesAccessible struct {
	All        []ApiResourceType
	Global     []ApiResourceType
	Namespaced []ApiResourceType
}

func ApiResourceTypes() ApiResourceTypesResponse {
	rawString := runCommand("api-resources", "-o", "wide")

	_ /* schema */, apiResourcesRaw := stringutil.ParseStdOutTable(rawString)

	allApiResources := funk.Map(apiResourcesRaw, func(in map[string]string) ApiResourceType {
		return ApiResourceType{
			Name:       stringutil.MapStrValOrElse(in, "NAME", ""),
			ShortNames: stringutil.CsvStr2arr(stringutil.MapStrValOrElse(in, "SHORTNAMES", "")),
			Namespaced: stringutil.Str2boolOrElse(stringutil.MapStrValOrElse(in, "NAMESPACED", ""), false),
			Kind:       stringutil.MapStrValOrElse(in, "KIND", ""),
			Verbs:      stringutil.WierdKubectlArray2arr(stringutil.MapStrValOrElse(in, "VERBS", "")),
		}
	}).([]ApiResourceType)

	accessibleApiResources := funk.Filter(allApiResources, func(r ApiResourceType) bool { return funk.ContainsString(r.Verbs, "get") }).([]ApiResourceType)
	globalResources := funk.Filter(accessibleApiResources, func(r ApiResourceType) bool { return !r.Namespaced }).([]ApiResourceType)
	namespacedResources := funk.Filter(accessibleApiResources, func(r ApiResourceType) bool { return r.Namespaced }).([]ApiResourceType)

	return ApiResourceTypesResponse{
		All: allApiResources,
		Accessible: ApiResourceTypesAccessible{
			All:        accessibleApiResources,
			Global:     globalResources,
			Namespaced: namespacedResources,
		},
	}
}

func runCommand(args ...string) string {
	if !shell.CommandExists("kubectl") {
		panic("kubectl not on path!")
	}

	fullCommand := "kubectl " + strings.Join(args, " ")

	cmd := exec.Command("kubectl", args...)

	outputBytes, err := cmd.Output()

	if err != nil {
		panic(`command "` + fullCommand + `" failed with error: ` + err.Error())
	}

	return strings.TrimSpace(string(outputBytes))
}

func removeK8sResourcePrefix(in string) string {
	return stringutil.RemoveUpToAndIncluding(in, "/")
}
