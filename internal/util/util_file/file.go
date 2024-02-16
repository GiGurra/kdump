package util_file

import (
	"github.com/samber/lo"
	"log"
	"os"
	"regexp"
	"strings"
)

func DeleteIfExists(path string, notDeterminableMsg string) {
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
	Bytes2File(path, []byte(data))
}

func File2String(path string) string {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic("Failed reading file '" + path + "' due to " + err.Error())
	}
	return string(bytes)
}

func ListFolder(path string) []os.DirEntry {
	info, err := os.ReadDir(path)
	if err != nil {
		panic("Failed listing files on path '" + path + "' due to " + err.Error())
	}
	return info
}

func ListFilesInFolder(path string) []os.DirEntry {
	all := ListFolder(path)
	return lo.Filter(all, func(i os.DirEntry, _ int) bool {
		return !i.IsDir()
	})
}

func ListFilesInFolderWithSuffix(path string, suffix string) []os.DirEntry {
	all := ListFilesInFolder(path)
	return lo.Filter(all, func(i os.DirEntry, _ int) bool {
		return strings.HasSuffix(strings.ToLower(i.Name()), suffix)
	})
}

func Bytes2File(path string, bytes []byte) {
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
