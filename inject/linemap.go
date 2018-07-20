package inject

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

var elementMaps map[string]*elementMap

func init() {
	elementMaps = make(map[string]*elementMap)
}

type elementMap struct {
	anthaElementPath string
	elementName      string
	lineMap          map[int]int
}

func RegisterLineMap(goElementPath, anElementPath, elementName string, lineMap map[int]int) {
	em := &elementMap{
		anthaElementPath: anElementPath,
		elementName:      elementName,
		lineMap:          lineMap,
	}
	elementMaps[goElementPath] = em
}

func ElementStackTrace() string {
	standard := string(debug.Stack())

	cs := make([]uintptr, 1000)
	// 5 is a magic number that allows us to skip back from where we
	// are here to the top of the meaningful stack that refers to the
	// panic. This is an unfortunate high-coupling to where
	// ElementStackTrace is being called from.
	num := runtime.Callers(5, cs)
	if num < 0 {
		return standard
	}
	frames := runtime.CallersFrames(cs[:num])
	// if we can find the first frame, then we treat it entirely as an
	// antha stack trace:
	first := true

	var strs []string
	for {
		frame, more := frames.Next()
		foundElement := false
		for suffix, em := range elementMaps {
			if strings.HasSuffix(frame.File, suffix) {
				foundElement = true
				lineStr := "(unknown line)"
				if line, foundLine := em.lineMap[frame.Line]; foundLine {
					lineStr = fmt.Sprint(line)
				}
				strs = append(strs, fmt.Sprintf("- [Element %s] %s:%s", em.elementName, em.anthaElementPath, lineStr))
				strs = append(strs, fmt.Sprintf("       [Go] %s", frame.Function))
				strs = append(strs, fmt.Sprintf("            %s:%d", frame.File, frame.Line))
				break
			}
		}
		if first && len(strs) == 0 { // was the first frame, and we failed to match an element, so abandon, and use standard
			return standard
		}
		first = false
		if !foundElement {
			strs = append(strs, fmt.Sprintf("- [Go] %s", frame.Function))
			strs = append(strs, fmt.Sprintf("       %s:%d", frame.File, frame.Line))
		}
		if !more {
			break
		}
	}
	return "Panic in Element\n" + strings.Join(strs, "\n")
}
