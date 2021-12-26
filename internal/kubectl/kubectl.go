package kubectl

import (
	"github.com/thoas/go-funk"
	"kdump/internal/shell"
	"kdump/internal/stringutil"
	"os/exec"
	"strings"
)

func removeK8sResourcePrefix(in string) string {
	return stringutil.RemoveUpToAndIncluding(in, "/")
}

func Namespaces() []string {
	arr := stringutil.SplitLines(runCommand("get", "namespaces", "-o", "name"))
	return funk.Map(arr, removeK8sResourcePrefix).([]string)
}

func CurrentNamespace() string {
	return runCommand("config", "view", "--minify", "--output", "jsonpath={..namespace}")
}

func CurrentContext() string {
	return runCommand("config", "current-context")
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
