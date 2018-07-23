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

// During code generation of elements, we append to the generated Go
// code a map that relates line numbers of Go back to the original
// line numbers of the element.an file. In the init() function of the
// generated elements, we call in here to RegisterLineMap in order to
// create a global mapping for every known element. This later allows
// us to rework panics and provide the original line numbers should a
// panic of some sort occur within an element.
func RegisterLineMap(goElementPath, anElementPath, elementName string, lineMap map[int]int) {
	em := &elementMap{
		anthaElementPath: anElementPath,
		elementName:      elementName,
		lineMap:          lineMap,
	}
	elementMaps[goElementPath] = em
}

// ElementStackTrace creates a stack trace, detecting whether or not
// the panic occured within an element. If the panic did not occur
// within an element, then the normal debug.Stack() is
// returned. Otherwise, we use the registered line maps to create a
// stack trace which refers back to the original elements, with the
// correct line numbers.
func ElementStackTrace() string {
	standard := string(debug.Stack())

	// This is a magic number :( It limits us to dealing with stack
	// traces that are 1000 frames deep. It is not expected this will
	// be a problem in practice!
	cs := make([]uintptr, 1000)

	// When a panic occurs, if a defer-with-recover has been
	// registered, the stack itself does not unwind. Instead, the
	// recover is invoked in a sub-frame:
	//
	// - Frame of defer-with-recover func
	// - Frame of function that panicked
	// - ... rest of call stack ...
	//
	// In the defer, if we call ElementStackTrace then that adds
	// another frame, and then we further have to call runtime.Callers
	// in order to generate the stack frame. At this point, the stack
	// will look like this:
	//
	// - Frame of runtime.Callers
	// - Frame of ElementStackTrace
	// - Frame of defer-with-recover func
	// - Frame of panic itself
	// - Frame of function that panicked
	// - ... rest of call stack ...
	//
	// Therefore, we have to skip over the first 5 items at the top of
	// the call stack in order to find the entry we're interested in.
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
