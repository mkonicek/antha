package liquidhandling

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

// ObjectCategory enumerate the different types of objects which can appear on deck
type ObjectCategory int

const (
	Tipboxes  ObjectCategory = iota
	Inputs                   // plates containing input solutions
	Outputs                  // plates that will have output solutions added to them
	Tipwastes                // for disposing of tips - often device specific
	Wastes                   // for disponsing of waste liquids
	Washes                   // for washing
)

// this should require updating only very rarely, so probably OK to do it like this
var categoryNames = map[ObjectCategory]string{
	Tipboxes:  "tipboxes",
	Inputs:    "inputs",
	Outputs:   "outputs",
	Tipwastes: "tipwastes",
	Wastes:    "wastes",
	Washes:    "washes",
}

// String description of the category suitable for showing in error messages
func (c ObjectCategory) String() string {
	if ret, ok := categoryNames[c]; ok {
		return ret
	} else {
		panic("unknown category")
	}
}

type Addresses []string

// Map returns a set containing the addresses
func (a Addresses) Map() map[string]bool {
	ret := make(map[string]bool, len(a))
	for _, address := range a {
		ret[address] = true
	}
	return ret
}

// Filter returns a new address slice containing the addresses in b that also appear in a
func (a Addresses) Filter(b Addresses) Addresses {
	inA := a.Map()
	ret := make(Addresses, len(b))
	for _, address := range b {
		if inA[address] {
			ret = append(ret, address)
		}
	}

	return ret
}

// Remove return a new address slice containing the addresses in a which don't appear in b
func (a Addresses) Remove(b Addresses) Addresses {
	inB := b.Map()
	ret := make(Addresses, 0, len(b))
	for _, address := range a {
		if !inB[address] {
			ret = append(ret, address)
		}
	}
	return ret
}

// String list the addresses wrapped in inverted commas
func (a Addresses) String() string {
	return fmt.Sprintf("\"%s\"", strings.Join([]string(a), "\", \""))
}

// LayoutOpt describes the options available for layout of various types of objects
// on a liquidhandler deck
type LayoutOpt map[ObjectCategory]Addresses

func (lo LayoutOpt) Dup() LayoutOpt {
	ret := make(LayoutOpt, len(lo))
	for k, v := range lo {
		c := make(Addresses, len(v))
		copy(c, v)
		ret[k] = c
	}
	return ret
}

func joinCategories(a, b LayoutOpt) []ObjectCategory {
	seen := make(map[ObjectCategory]bool, len(a)+len(b))
	for cat := range a {
		seen[cat] = true
	}
	for cat := range b {
		seen[cat] = true
	}

	ret := make([]ObjectCategory, 0, len(seen))
	for cat := range seen {
		ret = append(ret, cat)
	}
	return ret
}

// Merge combine the driver reported layout preferences with the users'
func (lo LayoutOpt) Merge(user LayoutOpt) error {
	categories := joinCategories(lo, user)

	// user preferences should be a strict subset of this
	errs := make([]string, 0, len(categories))
	for _, category := range categories {
		if invalid := user[category].Remove(lo[category]); len(invalid) > 0 {
			errs = append(errs, fmt.Sprintf("cannot place %s at %s", category, invalid))
		}
	}
	if len(errs) > 0 {
		if len(errs) == 1 {
			return errors.Errorf("invalid user preferences: %s", errs[0])
		} else {
			return errors.Errorf("invalid user preferences:\n\t%s", strings.Join(errs, "\n\t"))
		}
	}

	// override driver preferences with user preferences, if specified
	for _, category := range categories {
		if u, ok := user[category]; ok && len(u) > 0 {
			lo[category] = u
		}
	}
	return nil
}
