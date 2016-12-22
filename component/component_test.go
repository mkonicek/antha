package component

import (
	"testing"

	"github.com/antha-lang/antha/inject"
)

func TestMakeParams(t *testing.T) {
	type input struct {
		A int
	}
	type output struct {
		B int
	}

	comp := &Component{
		Constructor: func() interface{} {
			return &inject.CheckedRunner{
				In:  &input{},
				Out: &output{},
			}
		},
	}
	if err := UpdateParamTypes(comp); err != nil {
		t.Fatal(err)
	}
	params, err := comp.NewParams()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := params["A"].(*int); !ok {
		t.Errorf("expecting *int found %T", params["A"])
	}
	if _, ok := params["B"].(*int); !ok {
		t.Errorf("expecting *int found %T", params["B"])
	}
}
