package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"

	"github.com/pkg/errors"
	"strings"
)

// TipFactory store descriptions of all tipboxes and tipwastes compatible with the device
type TipFactory struct {
	Tipboxes  map[string]*wtype.LHTipbox
	Tipwastes map[string]*wtype.LHTipwaste
}

// NewTipFactory initialise a new tip factory
func NewTipFactory(tipboxes []*wtype.LHTipbox, tipwastes []*wtype.LHTipwaste) *TipFactory {
	ret := &TipFactory{
		Tipboxes:  make(map[string]*wtype.LHTipbox, len(tipboxes)),
		Tipwastes: make(map[string]*wtype.LHTipwaste, len(tipwastes)),
	}

	for _, tb := range tipboxes {
		ret.Tipboxes[tb.Type] = tb
	}
	for _, tw := range tipwastes {
		ret.Tipwastes[tw.Type] = tw
	}
	return ret
}

// NewTipbox creates a new tipbox of the given type, returning an error if no such tipbox is known
func (tf *TipFactory) NewTipbox(name string) (*wtype.LHTipbox, error) {
	if tb, ok := tf.Tipboxes[name]; !ok {
		types := make([]string, 0, len(tf.Tipboxes))
		for n := range tf.Tipboxes {
			types = append(types, n)
		}
		return nil, errors.Errorf("cannot create tipbox: unknown tip type %q, valid types are %s", name, strings.Join(types, ", "))
	} else {
		return tb.Dup(), nil
	}
}

// Dup return a copy of the factory
func (tf *TipFactory) Dup() *TipFactory {
	ret := &TipFactory{
		Tipboxes:  make(map[string]*wtype.LHTipbox, len(tf.Tipboxes)),
		Tipwastes: make(map[string]*wtype.LHTipwaste, len(tf.Tipwastes)),
	}

	for k, tb := range tf.Tipboxes {
		ret.Tipboxes[k] = tb.Dup()
	}
	for k, tw := range tf.Tipwastes {
		ret.Tipwastes[k] = tw.Dup()
	}
	return ret
}

// ConstrainTipboxTypes removes any tipbox types which are not present in names
func (tf *TipFactory) ConstrainTipboxTypes(names []string) {
	valid := make(map[string]bool, len(names))
	for _, name := range names {
		valid[name] = true
	}

	tipboxes := make(map[string]*wtype.LHTipbox, len(names))
	for name, tipbox := range tf.Tipboxes {
		if valid[name] {
			tipboxes[name] = tipbox
		}
	}

	tf.Tipboxes = tipboxes
}

// NewTipwaste creates a new tipwaste of the given type, returning an error if it isn't known
func (tf *TipFactory) NewTipwaste(name string) (*wtype.LHTipwaste, error) {
	if tw, ok := tf.Tipwastes[name]; !ok {
		return nil, errors.Errorf("cannot create tipwaste: unknown name %s", name)
	} else {
		return tw.Dup(), nil
	}
}
