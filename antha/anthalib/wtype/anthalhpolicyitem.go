package wtype

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

func GetPolicyConsequents() AnthaLHPolicyItemSet {
	return MakePolicyItems()
}

func MakePolicyItems() AnthaLHPolicyItemSet {
	typemap = maketypemap()
	alhpis := make(AnthaLHPolicyItemSet, 30)
	alhpis["ASPENTRYSPEED"] = AnthaLHPolicyItem{Name: "ASPENTRYSPEED", Type: typemap["float64"], Desc: "allows slow moves into liquids"}
	alhpis["ASPSPEED"] = AnthaLHPolicyItem{Name: "ASPSPEED", Type: typemap["float64"], Desc: "aspirate pipetting rate"}
	alhpis["ASPZOFFSET"] = AnthaLHPolicyItem{Name: "ASPZOFFSET", Type: typemap["float64"], Desc: "mm above well bottom when aspirating"}
	alhpis["ASP_WAIT"] = AnthaLHPolicyItem{Name: "ASP_WAIT", Type: typemap["float64"], Desc: "wait time in seconds post aspirate"}
	alhpis["BLOWOUTOFFSET"] = AnthaLHPolicyItem{Name: "BLOWOUTOFFSET", Type: typemap["float64"], Desc: "mm above BLOWOUTREFERENCE"}
	alhpis["BLOWOUTREFERENCE"] = AnthaLHPolicyItem{Name: "BLOWOUTREFERENCE", Type: typemap["int"], Desc: "where to be when blowing out: 0 well bottom"}
	alhpis["BLOWOUTVOLUME"] = AnthaLHPolicyItem{Name: "BLOWOUTVOLUME", Type: typemap["float64"], Desc: "how much to blow out"}
	alhpis["CAN_MULTI"] = AnthaLHPolicyItem{Name: "CAN_MULTI", Type: typemap["bool"], Desc: "is multichannel operation allowed?"}
	alhpis["DSPENTRYSPEED"] = AnthaLHPolicyItem{Name: "DSPENTRYSPEED", Type: typemap["float64"], Desc: "allows slow moves into liquids"}
	alhpis["DSPREFERENCE"] = AnthaLHPolicyItem{Name: "DSPREFERENCE", Type: typemap["int"], Desc: "where to be when dispensing: 0 well bottom"}
	alhpis["DSPSPEED"] = AnthaLHPolicyItem{Name: "DSPSPEED", Type: typemap["float64"], Desc: "dispense pipetting rate"}
	alhpis["DSPZOFFSET"] = AnthaLHPolicyItem{Name: "DSPZOFFSET", Type: typemap["float64"], Desc: "mm above DSPREFERENCE"}
	alhpis["DSP_WAIT"] = AnthaLHPolicyItem{Name: "DSP_WAIT", Type: typemap["float64"], Desc: "wait time in seconds post dispense"}
	alhpis["EXTRA_ASP_VOLUME"] = AnthaLHPolicyItem{Name: "EXTRA_ASP_VOLUME", Type: typemap["Volume"], Desc: "additional volume to take up when aspirating"}
	alhpis["EXTRA_DISP_VOLUME"] = AnthaLHPolicyItem{Name: "EXTRA_DISP_VOLUME", Type: typemap["Volume"], Desc: "additional volume to dispense"}
	alhpis["JUSTBLOWOUT"] = AnthaLHPolicyItem{Name: "JUSTBLOWOUT", Type: typemap["bool"], Desc: "shortcut to get single transfer"}
	alhpis["POST_MIX"] = AnthaLHPolicyItem{Name: "POST_MIX", Type: typemap["int"], Desc: "number of mix cycles to do after dispense"}
	alhpis["POST_MIX_RATE"] = AnthaLHPolicyItem{Name: "POST_MIX_RATE", Type: typemap["float64"], Desc: "pipetting rate when post mixing"}
	alhpis["POST_MIX_VOL"] = AnthaLHPolicyItem{Name: "POST_MIX_VOL", Type: typemap["float64"], Desc: "volume to post mix (ul)"}
	alhpis["POST_MIX_X"] = AnthaLHPolicyItem{Name: "POST_MIX_X", Type: typemap["float64"], Desc: "x offset from centre of well (mm) when post-mixing"}
	alhpis["POST_MIX_Y"] = AnthaLHPolicyItem{Name: "POST_MIX_Y", Type: typemap["float64"], Desc: "y offset from centre of well (mm) when post-mixing"}
	alhpis["POST_MIX_Z"] = AnthaLHPolicyItem{Name: "POST_MIX_Z", Type: typemap["float64"], Desc: "z offset from centre of well (mm) when post-mixing"}
	alhpis["PRE_MIX"] = AnthaLHPolicyItem{Name: "PRE_MIX", Type: typemap["int"], Desc: "number of mix cycles to do before aspirating"}
	alhpis["PRE_MIX_RATE"] = AnthaLHPolicyItem{Name: "PRE_MIX_RATE", Type: typemap["float64"], Desc: "pipetting rate when pre mixing"}
	alhpis["PRE_MIX_VOL"] = AnthaLHPolicyItem{Name: "PRE_MIX_VOL", Type: typemap["float64"], Desc: "volume to pre mix (ul)"}
	alhpis["PRE_MIX_X"] = AnthaLHPolicyItem{Name: "PRE_MIX_X", Type: typemap["float64"], Desc: "x offset from centre of well (mm) when pre-mixing"}
	alhpis["PRE_MIX_Y"] = AnthaLHPolicyItem{Name: "PRE_MIX_Y", Type: typemap["float64"], Desc: "y offset from centre of well (mm) when pre-mixing"}
	alhpis["PRE_MIX_Z"] = AnthaLHPolicyItem{Name: "PRE_MIX_Z", Type: typemap["float64"], Desc: "z offset from centre of well (mm) when pre-mixing"}
	alhpis["TIP_REUSE_LIMIT"] = AnthaLHPolicyItem{Name: "TIP_REUSE_LIMIT", Type: typemap["int"], Desc: "number of times tips can be reused for asp/dsp cycles"}
	alhpis["TOUCHOFF"] = AnthaLHPolicyItem{Name: "TOUCHOFF", Type: typemap["bool"], Desc: "whether to move to TOUCHOFFSET after dispense"}
	alhpis["TOUCHOFFSET"] = AnthaLHPolicyItem{Name: "TOUCHOFFSET", Type: typemap["float64"], Desc: "mm above wb to touch off at"}
	return alhpis
}

// defines possible liquid handling policy consequences
type AnthaLHPolicyItem struct {
	Name string
	Type reflect.Type
	Desc string
}

func (alhpi AnthaLHPolicyItem) TypeName() string {
	return alhpi.Type.Name()
}

var typemap map[string]reflect.Type

func maketypemap() map[string]reflect.Type {
	// prototypical types for map
	var f float64
	var i int
	var s string
	var v wunit.Volume
	var b bool

	ret := make(map[string]reflect.Type, 4)
	ret["float64"] = reflect.TypeOf(f)
	ret["int"] = reflect.TypeOf(i)
	ret["string"] = reflect.TypeOf(s)
	ret["Volume"] = reflect.TypeOf(v)
	ret["bool"] = reflect.TypeOf(b)

	return ret
}

type AnthaLHPolicyItemSet map[string]AnthaLHPolicyItem

func (alhpis AnthaLHPolicyItemSet) TypeList() string {
	ks := make([]string, 0, len(alhpis))

	for k, _ := range alhpis {
		ks = append(ks, k)
	}

	sort.Strings(ks)

	s := ""

	for _, k := range ks {
		alhpi := alhpis[k]
		s += fmt.Sprintf("%s,%s,%s\n", k, alhpi.TypeName(), alhpi.Desc)
	}

	return s
}

func (alhpis AnthaLHPolicyItemSet) CodeForIt() string {
	ks := make([]string, 0, len(alhpis))

	for k, _ := range alhpis {
		ks = append(ks, k)
	}

	sort.Strings(ks)

	s := ""

	s += "func MakePolicyItems()AnthaLHPolicyItemSet{\n"
	s += "ahlpis:=make(AnthaLHPolicyItemSet, 30)\n"
	for _, k := range ks {
		alhpi := alhpis[k]
		s += fmt.Sprintf("alhpis[\"%s\"] = AnthaLHPolicyItem{Name:\"%s\",Type:typemap[\"%s\"],Desc:\"%s\"}\n", k, k, alhpi.TypeName(), alhpi.Desc)
	}
	s += "return ahlpis\n"
	s += "}\n"

	return s
}
