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

type StdOutTableColumn struct {
	name       string
	byteIndex  int
	maxByteLen int
}

func ParseStdOutTable(table string) string {
	lines := SplitLines(table)
	headingLine := lines[0]
	//dataLines := lines[1:]

	beginIndices := make([]int, 0)

	fmt.Printf("len(lines): %d \n", len(lines))
	fmt.Printf("headingLine: %v \n", headingLine)

	prevIsSpace := true
	for i, r := range headingLine {
		if prevIsSpace && !unicode.IsSpace(r) {
			beginIndices = append(beginIndices, i)
		}
		prevIsSpace = unicode.IsSpace(r)
	}

	headings := make([]StdOutTableColumn, 0)
	for i, _ := range beginIndices {
		beginIndex := beginIndices[i]
		endIndex := len(headingLine)
		if i+1 < len(beginIndices) {
			endIndex = beginIndices[i+1] - 1
		}
		name := strings.TrimSpace(headingLine[beginIndex:endIndex])
		headings = append(headings, StdOutTableColumn{name, beginIndex, endIndex - beginIndex})
	}

	fmt.Printf("beginIndices: %v \n", beginIndices)
	fmt.Printf("headings: %+v \n", headings)
	//	fmt.Printf("dataLines: %v \n", dataLines)

	return ""
}
