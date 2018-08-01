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
	// - Frame of runtime.gopanic
	// - (optional) One or more frames of panic detail, eg
	//     runtime.panicdivide, runtime.panicmem, runtime.sigpanic
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
	// - Frame of runtime.gopanic
	// - (optional) One or more frames of panic detail, eg
	//     runtime.panicdivide, runtime.panicmem, runtime.sigpanic
	// - Frame of function that panicked
	// - ... rest of call stack ...
	//
	// Now we want to find if the the panic happened within an
	// element. To do that, we walk down until we find
	// runtime.gopanic. We then keep walking until we find the first
	// thing that doesn't start with runtime, and we inspect that!
	num := runtime.Callers(0, cs)
	if num < 0 {
		return standard
	}
	frames := runtime.CallersFrames(cs[:num])
	// skip over everything until we find the runtime.gopanic entry
	frame, more := frames.Next()
	for {
		if !more {
			return standard
		} else if frame.Function == "runtime.gopanic" {
			break
		} else {
			frame, more = frames.Next()
		}
	}
	// now keep going until we find the first thing that is _not_ runtime.:
	for {
		if !more {
			return standard
		} else if !strings.HasPrefix(frame.Function, "runtime.") {
			break
		} else {
			frame, more = frames.Next()
		}
	}

	// For each of the remaining frames, we need to detect whether they
	// are within an element or not. If the first of these remaining
	// frames is *not* within an element, then strs will remain empty,
	// and we will revert to using the standard stack trace.
	var strs []string
	for {
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
		if len(strs) == 0 { // we haven't been able to find any matching element, so use standard stack
			return standard
		}
		if !foundElement {
			strs = append(strs, fmt.Sprintf("- [Go] %s", frame.Function))
			strs = append(strs, fmt.Sprintf("       %s:%d", frame.File, frame.Line))
		}
		if !more {
			break
		}
		frame, more = frames.Next()
	}
	return strings.Join(strs, "\n")
}
