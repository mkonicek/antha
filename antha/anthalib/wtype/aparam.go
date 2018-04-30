package wtype

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

const PolicyNameField string = "PolicyName"

func GetPolicyConsequents() AParamSet {
	return MakePolicyItems()
}

func MakeInstructionParameters() AParamSet {
	typemap := maketypemap()

	params := make(AParamSet, 30)

	// TODO -- make parameter/instruction relation explicit
	params["BLOWOUT"] = AParam{Name: "BLOWOUT", Type: typemap["float64"]}
	params["CHANNEL"] = AParam{Name: "CHANNEL", Type: typemap["float64"]}
	params["CYCLES"] = AParam{Name: "CYCLES", Type: typemap["float64"]}
	params["DRIVE"] = AParam{Name: "DRIVE", Type: typemap["string"]}
	params["FROMPLATETYPE"] = AParam{Name: "FROMPLATETYPE", Type: typemap["string"]}
	params["HEAD"] = AParam{Name: "HEAD", Type: typemap["float64"]}
	params["INSTRUCTIONTYPE"] = AParam{Name: "INSTRUCTIONTYPE", Type: typemap["string"]}
	params["LIQUIDCLASS"] = AParam{Name: "LIQUIDCLASS", Type: typemap["string"]}
	params["LLF"] = AParam{Name: "LLF", Type: typemap["string"]} // actually bool but no checks on that yet
	params["MULTI"] = AParam{Name: "MULTI", Type: typemap["float64"]}
	params["NEWADAPTOR"] = AParam{Name: "NEWADAPTOR", Type: typemap["string"]}
	params["NEWSTATE"] = AParam{Name: "NEWSTATE", Type: typemap["string"]}
	params["OFFSETX"] = AParam{Name: "OFFSETX", Type: typemap["float64"]}
	params["OFFSETY"] = AParam{Name: "OFFSETY", Type: typemap["float64"]}
	params["OFFSETZ"] = AParam{Name: "OFFSETZ", Type: typemap["float64"]}
	params["OLDADAPTOR"] = AParam{Name: "OLDADAPTOR", Type: typemap["string"]}
	params["OLDSTATE"] = AParam{Name: "OLDSTATE", Type: typemap["string"]}
	params["OVERSTROKE"] = AParam{Name: "OVERSTROKE", Type: typemap["string"]} // bool
	params["PARAMS"] = AParam{Name: "PARAMS", Type: typemap["float64"]}
	params["PLATE"] = AParam{Name: "PLATE", Type: typemap["string"]}
	params["PLATETYPE"] = AParam{Name: "PLATETYPE", Type: typemap["string"]}
	params["PLATFORM"] = AParam{Name: "PLATFORM", Type: typemap["string"]}
	params["PLT"] = AParam{Name: "PLT", Type: typemap["string"]}
	params["POS"] = AParam{Name: "POS", Type: typemap["string"]}
	params["POSFROM"] = AParam{Name: "POSFROM", Type: typemap["string"]}
	params["POSTO"] = AParam{Name: "POSTO", Type: typemap["string"]}
	params["REFERENCE"] = AParam{Name: "REFERENCE", Type: typemap["float64"]}
	params["SPEED"] = AParam{Name: "SPEED", Type: typemap["float64"]}
	params["TIME"] = AParam{Name: "TIME", Type: typemap["float64"]}
	params["TIPTYPE"] = AParam{Name: "TIPTYPE", Type: typemap["string"]}
	params["TOPLATETYPE"] = AParam{Name: "TOPLATETYPE", Type: typemap["string"]}
	params["VOLUME"] = AParam{Name: "VOLUME", Type: typemap["float64"]}
	params["VOLUNT"] = AParam{Name: "VOLUNT", Type: typemap["float64"]}
	params["WELL"] = AParam{Name: "WELL", Type: typemap["float64"]}
	params["WELLFROM"] = AParam{Name: "WELLFROM", Type: typemap["string"]}
	params["WELLFROMVOLUME"] = AParam{Name: "WELLFROMVOLUME", Type: typemap["float64"]}
	params["WELLTO"] = AParam{Name: "WELLTO", Type: typemap["string"]}
	params["WELLTOVOLUME"] = AParam{Name: "WELLTOVOLUME", Type: typemap["float64"]}
	params["WELLVOLUME"] = AParam{Name: "WELLVOLUME", Type: typemap["float64"]}
	params["WHAT"] = AParam{Name: "WHAT", Type: typemap["string"]}
	return params
}

func MakePolicyItems() AParamSet {
	typemap = maketypemap()
	alhpis := make(AParamSet, 30)
	alhpis["ASPENTRYSPEED"] = AParam{Name: "ASPENTRYSPEED", Type: typemap["float64"], Desc: "allows slow moves into liquids"}
	alhpis["ASPREFERENCE"] = AParam{Name: "ASPREFERENCE", Type: typemap["int"], Desc: "where to be when aspirating: 0 well bottom, 1 well top, 2 liquid level (if known)"}
	alhpis["ASPSPEED"] = AParam{Name: "ASPSPEED", Type: typemap["float64"], Desc: "aspirate pipetting rate"}
	alhpis["ASPXOFFSET"] = AParam{Name: "ASPXOFFSET", Type: typemap["float64"], Desc: "mm east of well target when aspirating"}
	alhpis["ASPYOFFSET"] = AParam{Name: "ASPYOFFSET", Type: typemap["float64"], Desc: "mm south of well target  when aspirating"}
	alhpis["ASPZOFFSET"] = AParam{Name: "ASPZOFFSET", Type: typemap["float64"], Desc: "mm above ASPREFERENCE when aspirating"}
	alhpis["ASP_WAIT"] = AParam{Name: "ASP_WAIT", Type: typemap["float64"], Desc: "wait time in seconds post aspirate"}
	alhpis["BLOWOUTOFFSET"] = AParam{Name: "BLOWOUTOFFSET", Type: typemap["float64"], Desc: "mm above BLOWOUTREFERENCE"}
	alhpis["BLOWOUTREFERENCE"] = AParam{Name: "BLOWOUTREFERENCE", Type: typemap["int"], Desc: "where to be when blowing out: 0 well bottom"}
	alhpis["BLOWOUTVOLUME"] = AParam{Name: "BLOWOUTVOLUME", Type: typemap["float64"], Desc: "how much to blow out"}
	alhpis["BLOWOUTVOLUMEUNIT"] = AParam{Name: "BLOWOUTVOLUMEUNIT", Type: typemap["string"], Desc: "volume unit for blowout volume"}
	alhpis["CAN_MULTI"] = AParam{Name: "CAN_MULTI", Type: typemap["bool"], Desc: "is multichannel operation allowed?"}
	alhpis["DSPENTRYSPEED"] = AParam{Name: "DSPENTRYSPEED", Type: typemap["float64"], Desc: "allows slow moves into liquids"}
	alhpis["DSPREFERENCE"] = AParam{Name: "DSPREFERENCE", Type: typemap["int"], Desc: "where to be when dispensing: 0 well bottom, 1 well top, 2 liquid level (if known)"}
	alhpis["DSPSPEED"] = AParam{Name: "DSPSPEED", Type: typemap["float64"], Desc: "dispense pipetting rate"}
	alhpis["DSPXOFFSET"] = AParam{Name: "DSPXOFFSET", Type: typemap["float64"], Desc: "mm east of well target when dispensing"}
	alhpis["DSPYOFFSET"] = AParam{Name: "DSPYOFFSET", Type: typemap["float64"], Desc: "mm south of well target  when dispensing"}
	alhpis["DSPZOFFSET"] = AParam{Name: "DSPZOFFSET", Type: typemap["float64"], Desc: "mm above DSPREFERENCE"}
	alhpis["DSP_WAIT"] = AParam{Name: "DSP_WAIT", Type: typemap["float64"], Desc: "wait time in seconds post dispense"}
	alhpis["EXTRA_ASP_VOLUME"] = AParam{Name: "EXTRA_ASP_VOLUME", Type: typemap["Volume"], Desc: "additional volume to take up when aspirating"}
	alhpis["EXTRA_DISP_VOLUME"] = AParam{Name: "EXTRA_DISP_VOLUME", Type: typemap["Volume"], Desc: "additional volume to dispense"}
	alhpis["JUSTBLOWOUT"] = AParam{Name: "JUSTBLOWOUT", Type: typemap["bool"], Desc: "shortcut to get single transfer"}
	alhpis["OFFSETZADJUST"] = AParam{Name: "OFFSETZADJUST", Type: typemap["float64"], Desc: "Added to z offset"}
	alhpis["POST_MIX"] = AParam{Name: "POST_MIX", Type: typemap["int"], Desc: "number of mix cycles to do after dispense"}
	alhpis["MIX_VOLUME_OVERRIDE_TIP_MAX"] = AParam{Name: "MIX_VOLUME_OVERRIDE_TIP_MAX", Type: typemap["bool"], Desc: "Default to using the maximum volume for the current tip type if the specified post mix volume is too high"}
	alhpis["POST_MIX_RATE"] = AParam{Name: "POST_MIX_RATE", Type: typemap["float64"], Desc: "pipetting rate when post mixing"}
	alhpis["POST_MIX_VOLUME"] = AParam{Name: "POST_MIX_VOLUME", Type: typemap["float64"], Desc: "volume to post mix (ul)"}
	alhpis["POST_MIX_X"] = AParam{Name: "POST_MIX_X", Type: typemap["float64"], Desc: "x offset from centre of well (mm) when post-mixing"}
	alhpis["POST_MIX_Y"] = AParam{Name: "POST_MIX_Y", Type: typemap["float64"], Desc: "y offset from centre of well (mm) when post-mixing"}
	alhpis["POST_MIX_Z"] = AParam{Name: "POST_MIX_Z", Type: typemap["float64"], Desc: "z offset from centre of well (mm) when post-mixing"}
	alhpis["PRE_MIX"] = AParam{Name: "PRE_MIX", Type: typemap["int"], Desc: "number of mix cycles to do before aspirating"}
	alhpis["PRE_MIX_RATE"] = AParam{Name: "PRE_MIX_RATE", Type: typemap["float64"], Desc: "pipetting rate when pre mixing"}
	alhpis["PRE_MIX_VOLUME"] = AParam{Name: "PRE_MIX_VOLUME", Type: typemap["float64"], Desc: "volume to pre mix (ul)"}
	alhpis["PRE_MIX_X"] = AParam{Name: "PRE_MIX_X", Type: typemap["float64"], Desc: "x offset from centre of well (mm) when pre-mixing"}
	alhpis["PRE_MIX_Y"] = AParam{Name: "PRE_MIX_Y", Type: typemap["float64"], Desc: "y offset from centre of well (mm) when pre-mixing"}
	alhpis["PRE_MIX_Z"] = AParam{Name: "PRE_MIX_Z", Type: typemap["float64"], Desc: "z offset from centre of well (mm) when pre-mixing"}
	alhpis["RESET_OVERRIDE"] = AParam{Name: "RESET_OVERRIDE", Type: typemap["bool"], Desc: "Do not generate reset commands"}
	alhpis["TIP_REUSE_LIMIT"] = AParam{Name: "TIP_REUSE_LIMIT", Type: typemap["int"], Desc: "number of times tips can be reused for asp/dsp cycles"}
	alhpis["TOUCHOFF"] = AParam{Name: "TOUCHOFF", Type: typemap["bool"], Desc: "whether to move to TOUCHOFFSET after dispense"}
	alhpis["TOUCHOFFSET"] = AParam{Name: "TOUCHOFFSET", Type: typemap["float64"], Desc: "mm above wb to touch off at"}
	alhpis["DEFAULTPIPETTESPEED"] = AParam{Name: "DEFAULTPIPETTESPEED", Type: typemap["float64"], Desc: "Default pipette speed in ml/min"}
	alhpis["PTZOFFSET"] = AParam{Name: "PTZOFFSET", Type: typemap["float64"], Desc: "Z offset for pistons to zero"}
	alhpis["PTZREFERENCE"] = AParam{Name: "PTZREFERENCE", Type: typemap["int"], Desc: "Well reference for piston to zero: 0 = well bottom, 1 = well top, 2 = liquid level"}
	alhpis["CAN_SDD"] = AParam{Name: "CAN_SDD", Type: typemap["bool"], Desc: "Is it permissible just to do a one-shot transfer"}
	alhpis["ASPREFERENCE"] = AParam{Name: "ASPREFERENCE", Type: typemap["int"], Desc: "Reference point for aspirate: 0 = well bottom, 1 = well top, 2 = liquid level"}
	alhpis["MANUALPTZ"] = AParam{Name: "MANUALPTZ", Type: typemap["bool"], Desc: "Is explicit piston reset required? "}
	alhpis["DONT_BE_DIRTY"] = AParam{Name: "DONT_BE_DIRTY", Type: typemap["bool"], Desc: "Don't switch this off"}
	alhpis["NO_AIR_DISPENSE"] = AParam{Name: "NO_AIR_DISPENSE", Type: typemap["bool"], Desc: "Prevent dispensing anywhere other than under liquid?"}
	alhpis["CAN_MSA"] = AParam{Name: "CAN_MSA", Type: typemap["bool"], Desc: "Permissible to aspirate from multiple sources? -- currently non functional"}
	alhpis["DESCRIPTION"] = AParam{Name: "DESCRIPTION", Type: typemap["string"], Desc: "Summary of LHPolicy to present to the user"}
	alhpis["LLFBELOWSURFACE"] = AParam{Name: "LLFBELOWSURFACE", Type: typemap["float64"], Desc: "Distance below surface for Liquid Level Following (LLF) when aspirating"}
	alhpis["LLFABOVESURFACE"] = AParam{Name: "LLFABOVESURFACE", Type: typemap["float64"], Desc: "Distance below surface for Liquid Level Following (LLF) when dispensing"}
	alhpis[PolicyNameField] = AParam{Name: PolicyNameField, Type: typemap["string"], Desc: "Name of the Liquid Policy"}

	return alhpis
}

func GetLHPolicyOptions() AParamSet {
	ps := make(AParamSet, 5)
	tm := maketypemap()

	ps["USE_DRIVER_TIP_TRACKING"] = AParam{Name: "USE_DRIVER_TIP_TRACKING", Type: tm["bool"], Desc: "If driver has the option to use its own tip tracking, do so"}

	ps["USE_LLF"] = AParam{Name: "USE_LLF", Type: tm["bool"], Desc: "Use Liquid-level following if plate has a model for liquid height-volume relations and the driver can use it."}

	return ps
}

// a typed parameter, with description
type AParam struct {
	Name string
	Type reflect.Type
	Desc string
}

func (alhpi AParam) TypeName() string {
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

type AParamSet map[string]AParam

func (alhpis AParamSet) OrderedList() []string {
	ks := make([]string, 0, len(alhpis))

	for k := range alhpis {
		ks = append(ks, k)
	}

	sort.Strings(ks)

	return ks
}

func (alhpis AParamSet) TypeList() string {
	ks := alhpis.OrderedList()

	s := ""

	for _, k := range ks {
		alhpi := alhpis[k]
		s += fmt.Sprintf("%s,%s,%s\n", k, alhpi.TypeName(), alhpi.Desc)
	}

	return s
}

func (alhpis AParamSet) CodeForIt() string {
	ks := make([]string, 0, len(alhpis))

	for k := range alhpis {
		ks = append(ks, k)
	}

	sort.Strings(ks)

	s := ""

	s += "func MakePolicyItems()AParamSet{\n"
	s += "ahlpis:=make(AParamSet, 30)\n"
	for _, k := range ks {
		alhpi := alhpis[k]
		s += fmt.Sprintf("alhpis[\"%s\"] = AParam{Name:\"%s\",Type:typemap[\"%s\"],Desc:\"%s\"}\n", k, k, alhpi.TypeName(), alhpi.Desc)
	}
	s += "return ahlpis\n"
	s += "}\n"

	return s
}
