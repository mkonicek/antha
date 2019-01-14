package liquidhandling

import (
	"fmt"
	"github.com/antha-lang/antha/utils"
	"github.com/pkg/errors"
	"strings"
)

// Addresses an ordered list of names of positions within a liquidhandling device, as defined by the plugin for that device
type Addresses []string

// Map returns a set containing the addresses
func (a Addresses) Map() map[string]bool {
	ret := make(map[string]bool, len(a))
	for _, address := range a {
		ret[address] = true
	}
	return ret
}

// Dup return a copy of the list of addresses
func (a Addresses) Dup() Addresses {
	if a == nil {
		return nil
	}
	ret := make(Addresses, len(a))
	copy(ret, a)
	return ret
}

// String list the addresses wrapped in inverted commas
func (a Addresses) String() string {
	return fmt.Sprintf(`"%s"`, strings.Join([]string(a), `", "`))
}

// LayoutOpt describes the options available for layout of various types of objects
// on a liquidhandler deck, listed in priority order from highest to lowest
type LayoutOpt struct {
	Tipboxes  Addresses // locations where boxes of tips can be placed
	Inputs    Addresses // plates containing input solutions
	Outputs   Addresses // plates that will have output solutions added to them
	Tipwastes Addresses // for disposing of tips - often device specific
	Wastes    Addresses // for disponsing of waste liquids
	Washes    Addresses // for washing
}

// Dup return a copy of the layout options
func (lo *LayoutOpt) Dup() *LayoutOpt {
	return &LayoutOpt{
		Tipboxes:  lo.Tipboxes,
		Inputs:    lo.Inputs,
		Outputs:   lo.Outputs,
		Tipwastes: lo.Tipwastes,
		Wastes:    lo.Wastes,
		Washes:    lo.Washes,
	}
}

// ApplyUserPreferences combine the user supplied preferences with the driver plugin supplied rules
// An error is returned if the user requests an object be placed in a location that the driver
// has not allowed, i.e. the user preferences must be a subset of the driver rules,
// but can specify any order of the acceptable addresses.
//
// Nb. (HJK)
// If the user preferences for any category are empty, they are replaced by the driver's defaults.
// This is because currently the UI cannot access the set of driver defaults (and they will soon
// vary with the specific device instance not just device type), so a blank entry signifies defaults.
// This has the side effect that it is not possible for the user to specify that there are no positions
// for objects of a particular type.
// This seems like an unlikely use case, so it was decided to proceed as is until there is UX in place
// to correctly set these preferences
func (lo *LayoutOpt) ApplyUserPreferences(user *LayoutOpt) (*LayoutOpt, error) {

	// override driver preferences with user preferences, if specified
	override := func(target *Addresses, driver, user Addresses, name string) error {
		// return an error if the user requests an address that's not OK
		valid := driver.Map()
		invalid := make(Addresses, 0, len(user))
		for _, u := range user {
			if !valid[u] {
				invalid = append(invalid, u)
			}
		}
		if len(invalid) > 0 {
			return errors.Errorf("cannot place %s at: %s", name, invalid.String())
		}

		if len(user) > 0 {
			*target = user.Dup()
		} else {
			*target = driver.Dup()
		}
		return nil
	}

	ret := &LayoutOpt{}
	errs := utils.ErrorSlice{
		override(&ret.Tipboxes, lo.Tipboxes, user.Tipboxes, "tipboxes"),
		override(&ret.Inputs, lo.Inputs, user.Inputs, "input plates"),
		override(&ret.Outputs, lo.Outputs, user.Outputs, "output plates"),
		override(&ret.Tipwastes, lo.Tipwastes, user.Tipwastes, "tipwastes"),
		override(&ret.Wastes, lo.Wastes, user.Wastes, "liquid wastes"),
		override(&ret.Washes, lo.Washes, user.Washes, "tip washes"),
	}
	return ret, errors.WithMessage(errs.Pack(), "invalid layout preferences")

}
