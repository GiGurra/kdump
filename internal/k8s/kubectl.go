package k8s

import (
	"github.com/samber/lo"
	"log"
	"os/exec"
	"strings"
)

func commandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func init() {

	// Check we have
	log.Printf("Checking that kubectl is installed...")
	if !commandExists("kubectl") {
		panic("kubectl not on path!")
	}

	log.Printf("Checking that kubectl neat is installed...")
	RunCommand("kubectl", "neat", "--help")
}

func DownloadEverything(types []*ApiResourceType) string {

	qualifiedTypeNames := lo.Map(types, func(in *ApiResourceType, _ int) string {
		return in.QualifiedName
	})

	return RunCommand("kubectl", "get", strings.Join(qualifiedTypeNames, ","), "--all-namespaces", "-o", "yaml")
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

func ApiResourceTypes() ApiResourceTypesResponse {

	log.Printf("Checking what api resource types are available...\n")

	rawString := RunCommand("kubectl", "api-resources", "-o", "wide")

	apiResourcesRaw := ParseStdOutTable(rawString)

	allApiResources := lo.Map(apiResourcesRaw, func(in map[string]string, _ int) *ApiResourceType {

		out := &ApiResourceType{
			Name:       MapStrValOrElse(in, "NAME", ""),
			ShortNames: CsvStr2arr(MapStrValOrElse(in, "SHORTNAMES", "")),
			Namespaced: Str2boolOrElse(MapStrValOrElse(in, "NAMESPACED", ""), false),
			Kind:       MapStrValOrElse(in, "KIND", ""),
			Verbs:      WierdKubectlArray2arr(MapStrValOrElse(in, "VERBS", "")),
		}

		apiVersionString := MapStrValOrElse(in, "APIVERSION", "")
		parts := strings.Split(apiVersionString, "/")
		if len(parts) == 1 {
			out.ApiVersion = ApiVersion{Version: parts[0]}
			out.QualifiedName = out.Name
		} else if len(parts) == 2 {
			out.ApiVersion = ApiVersion{Version: parts[1], Name: parts[0]}
			out.QualifiedName = out.Name + "." + out.ApiVersion.Name
		} else {
			panic("Unable to parse ApiVersion from string: " + apiVersionString)
		}

		return out

	})

	accessibleApiResources := lo.Filter(allApiResources, func(r *ApiResourceType, _ int) bool { return lo.Contains(r.Verbs, "get") })
	globalResources := lo.Filter(accessibleApiResources, func(r *ApiResourceType, _ int) bool { return !r.Namespaced })
	namespacedResources := lo.Filter(accessibleApiResources, func(r *ApiResourceType, _ int) bool { return r.Namespaced })

	return ApiResourceTypesResponse{
		All: allApiResources,
		Accessible: ApiResourceTypesAccessible{
			All:        accessibleApiResources,
			Global:     globalResources,
			Namespaced: namespacedResources,
		},
	}
}

func RunCommand(app string, args ...string) string {

	fullCommand := app + " " + strings.Join(args, " ")

	cmd := exec.Command(app, args...)

	outputBytes, err := cmd.Output()

	if err != nil {
		panic(`command "` + fullCommand + `" failed with error: ` + err.Error())
	}

	return strings.TrimSpace(string(outputBytes))
}
