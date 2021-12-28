package fileutil

import (
	"log"
	"os"
	"regexp"
)

func Delete(path string, notDeterminableMsg string) {
	err := os.RemoveAll(path)
	if err != nil {
		panic(notDeterminableMsg)
	}
}

func Exists(path string, notDeterminableMsg string) bool {

	existingFolder, err := os.Stat(path)

	if existingFolder != nil {
		return true
	}

	if err != nil && !os.IsNotExist(err) {
		panic(notDeterminableMsg)
	}

	return false
}

func CreateFolderIfNotExists(path string, notPossibleMsg string) {

	err := os.MkdirAll(path, 0755)
	if err != nil {
		panic(notPossibleMsg + ", reason: " + err.Error())
	}
}

func CreateFolderIfMissing(path string, notPossibleMsg string) {

	if !Exists(path, notPossibleMsg) {
		err := os.MkdirAll(path, 0755)
		if err != nil && !os.IsExist(err) {
			panic(notPossibleMsg)
		}
	}
}

func String2File(path string, data string) {
	bytes := []byte(data)
	err := os.WriteFile(path, bytes, 0644)
	if err != nil {
		panic("Failed writing to file '" + path + "' due to " + err.Error())
	}
}

func SanitizePath(filename string) string {

	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9\\-_.]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(filename, "_")

	return processedString
}
