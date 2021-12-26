package fileutil

import "os"

func PanicIfExists(path string, existsMsg string, notDeterminableMsg string) {

	existingFolder, err := os.Stat(path)

	if existingFolder != nil {
		panic(existsMsg)
	}

	if err != nil && !os.IsNotExist(err) {
		panic(notDeterminableMsg)
	}
}
