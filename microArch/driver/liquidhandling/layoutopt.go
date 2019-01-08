package liquidhandling

import (
	"fmt"
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
// has not allowed
func (lo *LayoutOpt) ApplyUserPreferences(user *LayoutOpt) (*LayoutOpt, error) {

	// override driver preferences with user preferences, if specified
	override := func(driver, user Addresses) (Addresses, error) {
		// return an error if the user requests an address that's not OK
		valid := driver.Map()
		invalid := make(Addresses, 0, len(user))
		for _, u := range user {
			if !valid[u] {
				invalid = append(invalid, u)
			}
		}
		if len(invalid) > 0 {
			return nil, errors.New(invalid.String())
		}

		if len(user) > 0 {
			return user.Dup(), nil
		} else {
			return driver.Dup(), nil
		}
	}

	if tipboxes, err := override(lo.Tipboxes, user.Tipboxes); err != nil {
		return nil, errors.WithMessage(err, "invalid user preferences: cannot place tipboxes at")
	} else if inputs, err := override(lo.Inputs, user.Inputs); err != nil {
		return nil, errors.WithMessage(err, "invalid user preferences: cannot place inputs at")
	} else if outputs, err := override(lo.Outputs, user.Outputs); err != nil {
		return nil, errors.WithMessage(err, "invalid user preferences: cannot place outputs at")
	} else if tipwastes, err := override(lo.Tipwastes, user.Tipwastes); err != nil {
		return nil, errors.WithMessage(err, "invalid user preferences: cannot place tipwastes at")
	} else if wastes, err := override(lo.Wastes, user.Wastes); err != nil {
		return nil, errors.WithMessage(err, "invalid user preferences: cannot place wastes at")
	} else if washes, err := override(lo.Washes, user.Washes); err != nil {
		return nil, errors.WithMessage(err, "invalid user preferences: cannot place washes at")
	} else {
		return &LayoutOpt{
			Tipboxes:  tipboxes,
			Inputs:    inputs,
			Outputs:   outputs,
			Tipwastes: tipwastes,
			Wastes:    wastes,
			Washes:    washes,
		}, nil
	}
}
