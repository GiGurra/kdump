package shell

import "os/exec"

func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
