package wunit

import (
	"errors"
	"strings"
)

var (
	noMatch = errors.New("no match")
)

// Manual implementation of siunit.peg
type SIPrefixedUnitGrammar struct {
	SIPrefixedUnit
}

func (a *SIPrefixedUnitGrammar) Parse(in string) error {
	var pos int

	pos = a.matchPrefix(in, pos, false)
	if pos < 0 {
		pos = a.matchUnit(in, 0, true)
		if pos < 0 {
			return noMatch
		}
		a.AddUnitPlusPrefixNode()
		return nil
	}

	pos = a.matchUnit(in, pos, false)
	if pos < 0 {
		pos = a.matchUnit(in, 0, true)
		if pos < 0 {
			return noMatch
		}

		a.AddUnitPlusPrefixNode()
		return nil
	}

	pos = a.matchPrefix(in, 0, true)
	a.matchUnit(in, pos, true)
	a.AddUnitPlusPrefixNode()
	return nil
}

func (a *SIPrefixedUnitGrammar) matchPrefix(in string, pos int, consume bool) int {
	if len(in) <= pos {
		return -1
	}

	found := func(x string) int {
		if consume {
			a.AddUnitPrefix(x)
		}
		return len(x)
	}

	if strings.HasPrefix(in[pos:], "da") {
		return found("da")
	}

	switch in[pos] {
	case 'y':
		return found("y")
	case 'z':
		return found("z")
	case 'a':
		return found("a")
	case 'f':
		return found("f")
	case 'p':
		return found("p")
	case 'n':
		return found("n")
	case 'u':
		return found("u")
	case 'm':
		return found("m")
	case 'c':
		return found("c")
	case 'd':
		return found("d")
	case 'h':
		return found("h")
	case 'k':
		return found("k")
	case 'M':
		return found("M")
	case 'G':
		return found("G")
	case 'T':
		return found("T")
	case 'P':
		return found("P")
	case 'E':
		return found("E")
	case 'Z':
		return found("Z")
	case 'Y':
		return found("Y")
	}

	return -1
}

func (a *SIPrefixedUnitGrammar) matchUnit(in string, pos int, consume bool) int {
	if len(in) <= pos {
		return -1
	}

	found := func(x string) int {
		if consume {
			a.AddUnit(x)
		}
		return len(x)
	}

	switch {
	case strings.HasPrefix(in[pos:], "rads"):
		return found("rads")
	case strings.HasPrefix(in[pos:], "radians"):
		return found("radians")
	case strings.HasPrefix(in[pos:], "degrees"):
		return found("degrees")
	case strings.HasPrefix(in[pos:], "Hz"):
		return found("Hz")
	case strings.HasPrefix(in[pos:], "rpm"):
		return found("rpm")
	}

	switch in[pos] {
	case 'h':
		return found("h")
	case 'H':
		return found("H")
	case 'M':
		return found("M")
	case 'm':
		return found("m")
	case 'l':
		return found("l")
	case 'L':
		return found("L")
	case 'g':
		return found("g")
	case 'V':
		return found("V")
	case 'J':
		return found("J")
	case 'A':
		return found("A")
	case 'C':
		return found("C")
	case 'N':
		return found("N")
	case 's':
		return found("s")
	case '%':
		return found("%")
	}

	return -1
}
