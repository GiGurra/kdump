package stringutil

import (
	"bufio"
	"fmt"
	"github.com/thoas/go-funk"
	"strings"
	"unicode"
)

func SplitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}

func RemoveUpToAndIncluding(fullString string, key string) string {
	idx := strings.Index(fullString, key)
	if idx >= 0 {
		keyLen := len(key)
		return fullString[idx+keyLen:]
	} else {
		return fullString
	}
}

func MapStrArray(arr []string, mapFn func(string) string) []string {
	return funk.Map(arr, mapFn).([]string)
}

func ParseStdOutTable(table string) string {
	lines := SplitLines(table)
	headingLine := lines[0]
	//dataLines := lines[1:]

	beginIndices := make([]int, 0)
	endIndices := make([]int, 0)

	fmt.Printf("len(lines): %d \n", len(lines))
	fmt.Printf("headingLine: %v \n", headingLine)

	prevIsSpace := true
	for i, r := range headingLine {
		if prevIsSpace && !unicode.IsSpace(r) {
			beginIndices = append(beginIndices, i)
		}
		if !prevIsSpace && unicode.IsSpace(r) {
			endIndices = append(endIndices, i)
		}
		prevIsSpace = unicode.IsSpace(r)
	}

	if len(endIndices) < len(beginIndices) {
		endIndices = append(endIndices, len(headingLine))
	}

	fmt.Printf("beginIndices: %v \n", beginIndices)
	fmt.Printf("endIndices: %v \n", endIndices)
	//	fmt.Printf("dataLines: %v \n", dataLines)

	return ""
}
