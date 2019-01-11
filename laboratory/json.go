package laboratory

import (
	"fmt"
	"reflect"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/ugorji/go/codec"
)

type liquidJson struct {
	labBuild *LaboratoryBuilder
}

func (lj *liquidJson) ConvertExt(v interface{}) interface{} {
	panic("liquidJson.ConvertExt: not implemented")
}

func (lj *liquidJson) UpdateExt(dst interface{}, src interface{}) {
	if name, ok := src.(string); ok {
		if liquid, err := lj.labBuild.Inventory.NewComponent(name); err != nil {
			lj.labBuild.Fatal(err)
		} else {
			// nb dst is *always* a pointer (so *Liquid in this case)
			dstLiquid := dst.(*wtype.Liquid)
			*dstLiquid = *liquid
		}
	} else {
		lj.labBuild.Fatal(fmt.Errorf("Liquid specification in workflow should be a string (liquid name), but was: %T", src))
	}
}

type plateJson struct {
	labBuild *LaboratoryBuilder
}

func (pj *plateJson) ConvertExt(v interface{}) interface{} {
	panic("plateJson.ConvertExt: not implemented")
}

func (pj *plateJson) UpdateExt(dst interface{}, src interface{}) {
	if name, ok := src.(string); ok {
		if plate, err := pj.labBuild.Inventory.NewPlate(name); err != nil {
			pj.labBuild.Fatal(err)
		} else {
			dstPlate := dst.(*wtype.Plate)
			*dstPlate = *plate
		}
	} else {
		pj.labBuild.Fatal(fmt.Errorf("Plate specification in workflow should be a string (plate name), but was: %T", src))
	}
}

type tipboxJson struct {
	labBuild *LaboratoryBuilder
}

func (tj *tipboxJson) ConvertExt(v interface{}) interface{} {
	panic("tipboxJson.ConvertExt: not implemented")
}

func (tj *tipboxJson) UpdateExt(dst interface{}, src interface{}) {
	if name, ok := src.(string); ok {
		if tipbox, err := tj.labBuild.Inventory.NewTipbox(name); err != nil {
			tj.labBuild.Fatal(err)
		} else {
			dstTipBox := dst.(*wtype.LHTipbox)
			*dstTipBox = *tipbox
		}
	} else {
		tj.labBuild.Fatal(fmt.Errorf("Tipbox specification in workflow should be a string (tipbox name), but was: %T", src))
	}
}

func (labBuild *LaboratoryBuilder) RegisterJsonExtensions(jh *codec.JsonHandle) {
	if err := jh.SetInterfaceExt(reflect.TypeOf(wtype.Liquid{}), 0, &liquidJson{labBuild: labBuild}); err != nil {
		labBuild.Fatal(err)
	} else if err := jh.SetInterfaceExt(reflect.TypeOf(wtype.Plate{}), 0, &plateJson{labBuild: labBuild}); err != nil {
		labBuild.Fatal(err)
	} else if err := jh.SetInterfaceExt(reflect.TypeOf(wtype.LHTipbox{}), 0, &tipboxJson{labBuild: labBuild}); err != nil {
		labBuild.Fatal(err)
	}
}
