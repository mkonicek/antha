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
		if liquid, err := lj.labBuild.effects.Inventory.Components.NewComponent(name); err != nil {
			lj.labBuild.RecordError(err, true)
		} else {
			// nb dst is *always* a pointer (so *Liquid in this case)
			dstLiquid := dst.(*wtype.Liquid)
			*dstLiquid = *liquid
		}
	} else {
		lj.labBuild.RecordError(fmt.Errorf("Liquid specification in workflow should be a string (liquid name), but was: %T", src), true)
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
		if plate, err := pj.labBuild.effects.Inventory.Plates.NewPlate(wtype.PlateTypeName(name)); err != nil {
			pj.labBuild.RecordError(err, true)
		} else {
			dstPlate := dst.(*wtype.Plate)
			*dstPlate = *plate
		}
	} else {
		pj.labBuild.RecordError(fmt.Errorf("Plate specification in workflow should be a string (plate name), but was: %T", src), true)
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
		if tipbox, err := tj.labBuild.effects.Inventory.TipBoxes.NewTipbox(name); err != nil {
			tj.labBuild.RecordError(err, true)
		} else {
			dstTipBox := dst.(*wtype.LHTipbox)
			*dstTipBox = *tipbox
		}
	} else {
		tj.labBuild.RecordError(fmt.Errorf("Tipbox specification in workflow should be a string (tipbox name), but was: %T", src), true)
	}
}

func (labBuild *LaboratoryBuilder) RegisterJsonExtensions(jh *codec.JsonHandle) error {
	if err := jh.SetInterfaceExt(reflect.TypeOf(wtype.Liquid{}), 0, &liquidJson{labBuild: labBuild}); err != nil {
		labBuild.RecordError(err, true)
	} else if err := jh.SetInterfaceExt(reflect.TypeOf(wtype.Plate{}), 0, &plateJson{labBuild: labBuild}); err != nil {
		labBuild.RecordError(err, true)
	} else if err := jh.SetInterfaceExt(reflect.TypeOf(wtype.LHTipbox{}), 0, &tipboxJson{labBuild: labBuild}); err != nil {
		labBuild.RecordError(err, true)
	}
	return labBuild.Errors()
}
