package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"

	"github.com/pkg/errors"
	"strings"
)

// TipFactory store descriptions of all tipboxes and tipwastes compatible with the device
type TipFactory struct {
	Tipboxes          map[string]*wtype.LHTipbox
	Tipwastes         map[string]*wtype.LHTipwaste
	TipboxesByTipType map[string]string // maps from the name of the type type to the name of the tipbox
}

// NewTipFactory initialise a new tip factory
func NewTipFactory(tipboxes []*wtype.LHTipbox, tipwastes []*wtype.LHTipwaste) *TipFactory {
	ret := &TipFactory{
		Tipboxes:          make(map[string]*wtype.LHTipbox, len(tipboxes)),
		Tipwastes:         make(map[string]*wtype.LHTipwaste, len(tipwastes)),
		TipboxesByTipType: make(map[string]string, len(tipboxes)),
	}

	for _, tb := range tipboxes {
		ret.Tipboxes[tb.Type] = tb
		ret.TipboxesByTipType[tb.Tiptype.Type] = tb.Type
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
		return nil, errors.Errorf("cannot create tipbox: unknown tipbox type %q, valid types are %s", name, strings.Join(types, ", "))
	} else {
		return tb.Dup(), nil
	}
}

// NewTipboxByTipType create a new tipbox which contains tips of the given type
func (tf *TipFactory) NewTipboxByTipType(ttype string) (*wtype.LHTipbox, error) {
	if tbtype, ok := tf.TipboxesByTipType[ttype]; !ok {
		types := make([]string, 0, len(tf.TipboxesByTipType))
		for n := range tf.TipboxesByTipType {
			types = append(types, n)
		}
		return nil, errors.Errorf("cannot create tipbox: unknown tip type %q, valid types are %s", ttype, strings.Join(types, ", "))
	} else {
		return tf.NewTipbox(tbtype)
	}
}

// Dup return a copy of the factory
func (tf *TipFactory) Dup() *TipFactory {
	ret := &TipFactory{
		Tipboxes:          make(map[string]*wtype.LHTipbox, len(tf.Tipboxes)),
		Tipwastes:         make(map[string]*wtype.LHTipwaste, len(tf.Tipwastes)),
		TipboxesByTipType: make(map[string]string, len(tf.Tipboxes)),
	}

	for k, tb := range tf.Tipboxes {
		ret.Tipboxes[k] = tb.Dup()
		ret.TipboxesByTipType[tb.Tiptype.Type] = k
	}
	for k, tw := range tf.Tipwastes {
		ret.Tipwastes[k] = tw.Dup()
	}
	return ret
}

// ConstrainTipTypes removes any tip types which are not present in names
func (tf *TipFactory) ConstrainTipTypes(tipTypes []string) {

	tipboxes := make(map[string]*wtype.LHTipbox, len(tipTypes))
	tip2tb := make(map[string]string, len(tipTypes))

	for _, tipType := range tipTypes {
		if tbType, ok := tf.TipboxesByTipType[tipType]; !ok {
			continue
		} else if tb, ok := tf.Tipboxes[tbType]; ok {
			tipboxes[tbType] = tb
			tip2tb[tipType] = tbType
		}
	}

	tf.Tipboxes = tipboxes
	tf.TipboxesByTipType = tip2tb
}

// NewTipwaste creates a new tipwaste of the given type, returning an error if it isn't known
func (tf *TipFactory) NewTipwaste(name string) (*wtype.LHTipwaste, error) {
	if tw, ok := tf.Tipwastes[name]; !ok {
		return nil, errors.Errorf("cannot create tipwaste: unknown name %s", name)
	} else {
		return tw.Dup(), nil
	}
}
