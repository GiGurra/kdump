package stringutil

import (
	"bufio"
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
	name      string
	byteIndex int
}

func ParseStdOutTable(table string) ([]StdOutTableColumn, []map[string]string) {
	lines := SplitLines(table)
	headingLine := lines[0]
	dataLines := funk.FilterString(lines[1:], func(in string) bool {
		return len(strings.TrimSpace(in)) > 0
	})

	beginIndices := make([]int, 0)

	prevIsSpace := true
	for i, r := range headingLine {
		if prevIsSpace && !unicode.IsSpace(r) {
			beginIndices = append(beginIndices, i)
		}
		prevIsSpace = unicode.IsSpace(r)
	}

	headings := make([]StdOutTableColumn, 0)
	for i, beginIndex := range beginIndices {
		endIndex := len(headingLine)
		if i+1 < len(beginIndices) {
			endIndex = beginIndices[i+1]
		}
		name := strings.TrimSpace(headingLine[beginIndex:endIndex])
		headings = append(headings, StdOutTableColumn{name, beginIndex})
	}

	lineValues := make([]map[string]string, 0)

	for _, dataLine := range dataLines {
		lineValue := make(map[string]string, 0)
		for iHeading, heading := range headings {
			endIndex := 0
			if iHeading+1 < len(headings) {
				endIndex = headings[iHeading+1].byteIndex
			} else {
				endIndex = len(dataLine)
			}
			lineValue[heading.name] = strings.TrimSpace(dataLine[heading.byteIndex:endIndex])
		}
		lineValues = append(lineValues, lineValue)

	}

	return headings, lineValues
}
