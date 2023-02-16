package kubectl

import (
	"bufio"
	"github.com/samber/lo"
	"github.com/thoas/go-funk"
	"golang.org/x/exp/constraints"
	"log"
	"strconv"
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

func TrimSpaces(lines []string) []string {
	return funk.Map(lines, func(in string) string { return strings.TrimSpace(in) }).([]string)
}

func NonEmptyLines(source string) []string {
	return RemoveEmptyLines(SplitLines(source))
}

func RemoveEmptyLines(lines []string) []string {
	return funk.FilterString(TrimSpaces(lines), func(in string) bool { return len(in) > 0 })
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
	start      int
	end        int
	headingEnd int
}

func figureOutLayout(table string) []StdOutTableColumn {
	lines := SplitLines(table)
	headingLine := lines[0]
	startIndices := make([]int, 0)
	headingEndIndices := make([]int, 0)

	log.Println("Headings:: " + headingLine)

	/////////////////////////////////////////////////////////
	// Build knowledge of where data can be in the table :S

	prevIsSpace := true
	for i, r := range headingLine {
		if prevIsSpace && !unicode.IsSpace(r) {
			startIndices = append(startIndices, i)
		}
		if !prevIsSpace && unicode.IsSpace(r) {
			headingEndIndices = append(headingEndIndices, i)
		}
		prevIsSpace = unicode.IsSpace(r)
	}

	if len(headingEndIndices) < len(startIndices) {
		headingEndIndices = append(headingEndIndices, len(headingLine))
	}

	headings := lo.Map(lo.Zip2(startIndices, headingEndIndices), func(indcs lo.Tuple2[int, int], _ int) StdOutTableColumn {
		return StdOutTableColumn{name: headingLine[indcs.A:indcs.B], start: indcs.A, end: -1, headingEnd: indcs.B}
	})
	for i := range headings {
		if i+1 < len(headings) {
			headings[i].end = headings[i+1].start
		}
	}
	headings[len(headings)-1].end = lo.Max(lo.Map(lines, func(item string, _ int) int {
		return len(item)
	}))

	return headings
}
func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func ParseStdOutTable(table string) []map[string]string {

	layout := figureOutLayout(table)
	dataLines := lo.Filter(SplitLines(table)[1:], func(in string, _ int) bool {
		return len(strings.TrimSpace(in)) > 0
	})

	/////////////////////////////////////////////////////////
	// Extract the data
	log.Printf("  cols: %+v", layout)

	result := make([]map[string]string, 0)

	for _, dataLine := range dataLines {
		lineResult := make(map[string]string, 0)
		log.Println("Checking line: " + dataLine)
		for _, heading := range layout {
			log.Println("  Checking heading: " + heading.name)
			endIndex := min(heading.end, len(dataLine))
			if heading.start < endIndex {
				lineResult[heading.name] = strings.TrimSpace(dataLine[heading.start:endIndex])
			} else {
				lineResult[heading.name] = ""
			}
		}

		result = append(result, lineResult)

	}

	return result
}

func MapStrValOrElse(dict map[string]string, key string, fallback string) string {
	if val, ok := dict[key]; ok {
		return val
	} else {
		return fallback
	}
}

func Str2boolOrElse(str string, fallback bool) bool {
	if val, err := strconv.ParseBool(str); err == nil {
		return val
	} else {
		return fallback
	}
}

func CsvStr2arrSep(str string, sep string) []string {
	return MapStrArray(strings.Split(str, sep), func(in string) string {
		return strings.TrimSpace(in)
	})
}

func CsvStr2arr(str string) []string {
	return CsvStr2arrSep(str, ",")
}

func WierdKubectlArray2arr(strIn string) []string {
	if strings.HasPrefix(strIn, "[") && strings.HasSuffix(strIn, "]") {
		return CsvStr2arrSep(strIn[1:(len(strIn)-1)], " ")
	} else {
		return CsvStr2arrSep(strIn, ",")
	}
}
