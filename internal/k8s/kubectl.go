package k8s

import (
	"context"
	"fmt"
	"github.com/GiGurra/cmder"
	"github.com/GiGurra/kdump/internal/errh"
	"github.com/samber/lo"
	"io"
	"log/slog"
	"os/exec"
	"strings"
	"time"
)

func commandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func CheckNecessaryCliAppsAvailable() {

	// Check we have
	slog.Info("Checking that kubectl is installed...")
	if !commandExists("kubectl") {
		panic("kubectl not on path!")
	}

	slog.Info("Checking that kubectl neat is installed...")
	RunCommand("kubectl", "neat", "--help")
}

func ListAvailableContexts() []string {
	allString := RunCommand("kubectl", "config", "get-contexts", "-o", "name")
	lines := strings.Split(allString, "\n")
	trimmedLines := lo.Map(lines, func(in string, _ int) string { return strings.TrimSpace(in) })
	return lo.Filter(trimmedLines, func(in string, _ int) bool { return len(in) > 0 })
}

func CurrentContext() string {
	return strings.TrimSpace(RunCommand("kubectl", "config", "current-context"))
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

	slog.Info("Checking what api resource types are available...")

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

	res := cmder.
		NewA(app, args...).
		WithTotalTimeout(30 * time.Minute).
		Run(context.Background())

	if res.Err != nil {
		panic(fmt.Sprintf(`command "%s" failed with error code %d: %s`, fullCommand, res.ExitCode, res.Combined))
	}

	return strings.TrimSpace(res.StdOut)
}

func PipeToCommand(input string, app string, args ...string) string {

	subProcess := exec.Command(app, args...)
	stdin := errh.Unwrap(subProcess.StdinPipe()) // stdin must be opened before .start/.output

	go func() {
		errh.Unwrap(io.WriteString(stdin, input))
		errh.Ignore(stdin.Close())
	}()

	outputBytes := errh.Unwrap(subProcess.Output()) // starts and awaits the result of the child process

	return strings.TrimSpace(string(outputBytes))
}
