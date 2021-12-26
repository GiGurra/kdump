package kubectl

import (
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

func ApiResources() string {
	return runCommand("api-resources", "-o", "wide")
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
