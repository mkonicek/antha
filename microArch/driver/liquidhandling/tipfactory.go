package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"

	"github.com/pkg/errors"
)

type TipFactory struct {
	tipboxes  map[string]*wtype.LHTipbox
	tips      []*wtype.LHTip
	tipwastes map[string]*wtype.LHTipwaste
}

func NewTipFactory(tipboxes []*wtype.LHTipbox, tipwastes []*wtype.LHTipwaste) *TipFactory {
	ret := &TipFactory{
		tipboxes:  make(map[string]*wtype.LHTipbox, len(tipboxes)),
		tipwastes: make(map[string]*wtype.LHTipwaste, len(tipwastes)),
	}

	for _, tb := range tipboxes {
		ret.tipboxes[tb.Type] = tb
	}
	for _, tw := range tipwastes {
		ret.tipwastes[tw.Type] = tw
	}
	ret.updateTips()
	return ret
}

func (tf *TipFactory) NewTipbox(name string) (*wtype.LHTipbox, error) {
	if tb, ok := tf.tipboxes[name]; !ok {
		return nil, errors.Errorf("cannot create tipbox: unknown name %s", name)
	} else {
		return tb.Dup(), nil
	}
}

// Dup return a copy of the factory
func (tf *TipFactory) Dup() *TipFactory {
	ret := &TipFactory{
		tipboxes:  make(map[string]*wtype.LHTipbox, len(tf.tipboxes)),
		tipwastes: make(map[string]*wtype.LHTipwaste, len(tf.tipwastes)),
	}

	for _, tb := range tf.tipboxes {
		ret.tipboxes[tb.Type] = tb
	}
	for _, tw := range tf.tipwastes {
		ret.tipwastes[tw.Type] = tw
	}
	ret.updateTips()
	return ret
}

func (tf *TipFactory) TipboxTypes() []string {
	ret := make([]string, 0, len(tf.tipboxes))
	for n := range tf.tipboxes {
		ret = append(ret, n)
	}
	return ret
}

// Tips returns a list of all the tips available
func (tf *TipFactory) Tips() []*wtype.LHTip {
	return tf.tips
}

func (tf *TipFactory) updateTips() {
	tf.tips = make([]*wtype.LHTip, 0, len(tf.tipboxes))
	for _, tb := range tf.tipboxes {
		tf.tips = append(tf.tips, tb.Tiptype.Dup())
	}
}

// ConstrainTipboxTypes removes any tipbox types which are not present in names
func (tf *TipFactory) ConstrainTipboxTypes(names []string) {
	valid := make(map[string]bool, len(names))
	for _, name := range names {
		valid[name] = true
	}

	tipboxes := make(map[string]*wtype.LHTipbox, len(names))
	for name, tipbox := range tf.tipboxes {
		if valid[name] {
			tipboxes[name] = tipbox
		}
	}

	tf.tipboxes = tipboxes
	tf.updateTips()
}

func (tf *TipFactory) NewTipwaste(name string) (*wtype.LHTipwaste, error) {
	if tw, ok := tf.tipwastes[name]; !ok {
		return nil, errors.Errorf("cannot create tipwaste: unknown name %s", name)
	} else {
		return tw.Dup(), nil
	}
}

func (tf *TipFactory) TipwasteTypes() []string {
	ret := make([]string, 0, len(tf.tipwastes))
	for n := range tf.tipwastes {
		ret = append(ret, n)
	}
	return ret
}

func (tf *TipFactory) Tipwastes() []*wtype.LHTipwaste {
	ret := make([]*wtype.LHTipwaste, 0, len(tf.tipwastes))
	for _, tw := range tf.tipwastes {
		ret = append(ret, tw.Dup())
	}
	return ret
}
