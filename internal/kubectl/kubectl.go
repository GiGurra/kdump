package kubectl

import (
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v2"
	"kdump/internal/shell"
	"kdump/internal/stringutil"
	"os/exec"
	"strings"
)

func NamespacesOrPanic() []string {
	return stringutil.MapStrArray(stringutil.SplitLines(runCommandOrPanic("get", "namespaces", "-o", "name")), removeK8sResourcePrefix)
}

func CurrentNamespaceOrPanic() string {
	return runCommandOrPanic("config", "view", "--minify", "--output", "jsonpath={..namespace}")
}

func CurrentContextOrPanic() string {
	return runCommandOrPanic("config", "current-context")
}

func ListNamespacedResourcesOfTypeOrPanic(namespace string, resourceType string) []string {
	rawString := runCommandOrPanic("-n", namespace, "get", resourceType, "-o", "name")
	itemsWithPrefixes := stringutil.RemoveEmptyLines(stringutil.SplitLines(rawString))
	return funk.Map(itemsWithPrefixes, func(in string) string {
		return stringutil.RemoveUpToAndIncluding(in, "/")
	}).([]string)
}

func DownloadNamespacedResourceOrPanic(namespace string, resourceType string, resourceName string, format string) string {
	rawString := runCommandOrPanic("-n", namespace, "get", resourceType, resourceName, "-o", format)
	return rawString
}

func DownloadGlobalResourceOrPanic(resourceType string, resourceName string, format string) string {
	rawString := runCommandOrPanic("get", resourceType, resourceName, "-o", format)
	return rawString
}

func DownloadEverythingOrPanic(types []*ApiResourceType) string {

	qualifiedTypenames := funk.Map(types, func(in *ApiResourceType) string {
		return in.QualifiedName
	}).([]string)

	return runCommandOrPanic("get", strings.Join(qualifiedTypenames, ","), "--all-namespaces", "-o", "yaml")
}

func DownloadEverythingInNamespaceOrPanic(types []*ApiResourceType, namespace string) string {

	qualifiedTypenames := funk.Map(types, func(in *ApiResourceType) string {
		return in.QualifiedName
	}).([]string)

	return runCommandOrPanic("get", strings.Join(qualifiedTypenames, ","), "-n", namespace, "-o", "yaml")
}

func DownloadEverythingOfTypeOrPanic(tpe *ApiResourceType) string {
	return runCommandOrPanic("get", tpe.QualifiedName, "-o", "yaml")
}

func DownloadEverythingOfTypeInNamespaceOrPanic(tpe *ApiResourceType, namespace string) string {
	return runCommandOrPanic("get", "-n", namespace, tpe.QualifiedName, "-o", "yaml")
}

func ParseK8sYamlOrPanic(in string) []*K8sResource {

	root := make(map[string]interface{})

	err := yaml.Unmarshal([]byte(in), &root)

	if err != nil {
		panic("Could not parse input yaml due to " + err.Error())
	}

	if hasKey(root, "items") {
		return parseMultiK8sYaml(root)
	} else {
		return []*K8sResource{parseSingleK8sYaml(root)}
	}
}

type K8sResource struct {
	Kind              string              `yaml:"kind"`
	ApiVersion        string              `yaml:"apiVersion"`
	MetaData          K8sResourceMetadata `yaml:"metadata"`
	SourceYaml        string
	QualifiedTypeName string
}

func (r *K8sResource) IsNamespaced() bool {
	return len(r.MetaData.Namespace) > 0
}

func (r *K8sResource) IsGlobal() bool {
	return !r.IsNamespaced()
}

func (r *K8sResource) IsSecret() bool {
	return strings.ToLower(r.Kind) == "secret"
}

type K8sResourceMetadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

func parseMultiK8sYaml(in map[string]interface{}) []*K8sResource {
	items := in["items"].([]interface{})
	return funk.Map(items, parseSingleK8sYaml).([]*K8sResource)
}

func parseSingleK8sYaml(item interface{}) *K8sResource {

	m := item.(map[interface{}]interface{})

	delete(m, "lastRefresh")
	delete(m, "status")

	yamlString, marshallError := yaml.Marshal(m)
	if marshallError != nil {
		panic("Failed marshalling object to yaml, due to: " + marshallError.Error())
	}
	var out K8sResource
	if yaml.Unmarshal(yamlString, &out) != nil {
		panic("Failed reading back yaml: " + string(yamlString))
	}

	out.SourceYaml = string(yamlString)

	parts := strings.Split(out.ApiVersion, "/")
	if len(parts) == 1 {
		out.QualifiedTypeName = strings.ToLower(out.Kind)
	} else if len(parts) == 2 {
		out.QualifiedTypeName = strings.ToLower(out.Kind) + "." + parts[0]
	} else {
		panic("Unable to parse QualifiedTypeName from string: " + out.ApiVersion)
	}

	return &out
}

func hasKey(dict map[string]interface{}, key string) bool {
	if _, ok := dict[key]; ok {
		return true
	} else {
		return false
	}
}

type ApiVersion struct {
	Version string
	Name    string
}
type ApiResourceType struct {
	Name          string
	ShortNames    []string
	Namespaced    bool
	Kind          string
	Verbs         []string
	ApiVersion    ApiVersion
	QualifiedName string
}

type ApiResourceTypesResponse struct {
	All        []*ApiResourceType
	Accessible ApiResourceTypesAccessible
}

type ApiResourceTypesAccessible struct {
	All        []*ApiResourceType
	Global     []*ApiResourceType
	Namespaced []*ApiResourceType
}

func ApiResourceTypesOrPanic() ApiResourceTypesResponse {

	rawString := runCommandOrPanic("api-resources", "-o", "wide")

	_ /* schema */, apiResourcesRaw := stringutil.ParseStdOutTable(rawString)

	allApiResources := funk.Map(apiResourcesRaw, func(in map[string]string) *ApiResourceType {

		out := &ApiResourceType{
			Name:       stringutil.MapStrValOrElse(in, "NAME", ""),
			ShortNames: stringutil.CsvStr2arr(stringutil.MapStrValOrElse(in, "SHORTNAMES", "")),
			Namespaced: stringutil.Str2boolOrElse(stringutil.MapStrValOrElse(in, "NAMESPACED", ""), false),
			Kind:       stringutil.MapStrValOrElse(in, "KIND", ""),
			Verbs:      stringutil.WierdKubectlArray2arr(stringutil.MapStrValOrElse(in, "VERBS", "")),
		}

		apiVersionString := stringutil.MapStrValOrElse(in, "APIVERSION", "")
		parts := strings.Split(apiVersionString, "/")
		if len(parts) == 1 {
			out.ApiVersion = ApiVersion{parts[0], ""}
			out.QualifiedName = out.Name
		} else if len(parts) == 2 {
			out.ApiVersion = ApiVersion{parts[1], parts[0]}
			out.QualifiedName = out.Name + "." + out.ApiVersion.Name
		} else {
			panic("Unable to parse ApiVersion from string: " + apiVersionString)
		}

		return out

	}).([]*ApiResourceType)

	accessibleApiResources := funk.Filter(allApiResources, func(r *ApiResourceType) bool { return funk.ContainsString(r.Verbs, "get") }).([]*ApiResourceType)
	globalResources := funk.Filter(accessibleApiResources, func(r *ApiResourceType) bool { return !r.Namespaced }).([]*ApiResourceType)
	namespacedResources := funk.Filter(accessibleApiResources, func(r *ApiResourceType) bool { return r.Namespaced }).([]*ApiResourceType)

	return ApiResourceTypesResponse{
		All: allApiResources,
		Accessible: ApiResourceTypesAccessible{
			All:        accessibleApiResources,
			Global:     globalResources,
			Namespaced: namespacedResources,
		},
	}
}

func runCommandOrPanic(args ...string) string {
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
