package wtype

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

const (
	// PolicyNameField returns the map entry of a liquid policy corresponding to the name of the policy.
	PolicyNameField string = "POLICYNAME"

	// PolicyNameField returns the map entry of a liquid policy corresponding to the name of the policy.
	PolicyDescriptionField string = "DESCRIPTION"

	// LiquidClass is the name of the liquid type checked at instruction generation.
	// Currently this is analogous to the POLICYNAME
	LiquidClass string = "LIQUIDCLASS"

	// This value will be used for aspirating, dispensing and mixing.
	DefaultPipetteSpeed string = "DEFAULTPIPETTESPEED"
)

func GetPolicyConsequents() AParamSet {
	return MakePolicyItems()
}

func MakeInstructionParameters() AParamSet {
	typemap := maketypemap()

	// TODO -- make parameter/instruction relation explicit
	return AParamSet{
		"BLOWOUT":         AParam{Name: "BLOWOUT", Type: typemap[Float64]},
		"CHANNEL":         AParam{Name: "CHANNEL", Type: typemap[Float64]},
		"CYCLES":          AParam{Name: "CYCLES", Type: typemap[Float64]},
		"DRIVE":           AParam{Name: "DRIVE", Type: typemap[String]},
		"FROMPLATETYPE":   AParam{Name: "FROMPLATETYPE", Type: typemap[String]},
		"HEAD":            AParam{Name: "HEAD", Type: typemap[Float64]},
		"INSTRUCTIONTYPE": AParam{Name: "INSTRUCTIONTYPE", Type: typemap[String]},
		LiquidClass:       AParam{Name: LiquidClass, Type: typemap[String]},
		"LLF":             AParam{Name: "LLF", Type: typemap[String]}, // actually bool but no checks on that yet
		"MULTI":           AParam{Name: "MULTI", Type: typemap[Float64]},
		"NEWADAPTOR":      AParam{Name: "NEWADAPTOR", Type: typemap[String]},
		"NEWSTATE":        AParam{Name: "NEWSTATE", Type: typemap[String]},
		"OFFSETX":         AParam{Name: "OFFSETX", Type: typemap[Float64]},
		"OFFSETY":         AParam{Name: "OFFSETY", Type: typemap[Float64]},
		"OFFSETZ":         AParam{Name: "OFFSETZ", Type: typemap[Float64]},
		"OLDADAPTOR":      AParam{Name: "OLDADAPTOR", Type: typemap[String]},
		"OLDSTATE":        AParam{Name: "OLDSTATE", Type: typemap[String]},
		"OVERSTROKE":      AParam{Name: "OVERSTROKE", Type: typemap[String]}, // bool
		"PARAMS":          AParam{Name: "PARAMS", Type: typemap[Float64]},
		"PLATE":           AParam{Name: "PLATE", Type: typemap[String]},
		"PLATETYPE":       AParam{Name: "PLATETYPE", Type: typemap[String]},
		"PLATFORM":        AParam{Name: "PLATFORM", Type: typemap[String]},
		"PLT":             AParam{Name: "PLT", Type: typemap[String]},
		"POS":             AParam{Name: "POS", Type: typemap[String]},
		"POSFROM":         AParam{Name: "POSFROM", Type: typemap[String]},
		"POSTO":           AParam{Name: "POSTO", Type: typemap[String]},
		"REFERENCE":       AParam{Name: "REFERENCE", Type: typemap[Float64]},
		"SPEED":           AParam{Name: "SPEED", Type: typemap[Float64]},
		"TIME":            AParam{Name: "TIME", Type: typemap[Float64]},
		"TIPTYPE":         AParam{Name: "TIPTYPE", Type: typemap[String]},
		"TOPLATETYPE":     AParam{Name: "TOPLATETYPE", Type: typemap[String]},
		"VOLUME":          AParam{Name: "VOLUME", Type: typemap[Float64]},
		"VOLUNT":          AParam{Name: "VOLUNT", Type: typemap[Float64]},
		"WELL":            AParam{Name: "WELL", Type: typemap[Float64]},
		"WELLFROM":        AParam{Name: "WELLFROM", Type: typemap[String]},
		"WELLFROMVOLUME":  AParam{Name: "WELLFROMVOLUME", Type: typemap[Float64]},
		"WELLTO":          AParam{Name: "WELLTO", Type: typemap[String]},
		"WELLTOVOLUME":    AParam{Name: "WELLTOVOLUME", Type: typemap[Float64]},
		"WELLVOLUME":      AParam{Name: "WELLVOLUME", Type: typemap[Float64]},
		"WHAT":            AParam{Name: "WHAT", Type: typemap[String]},
	}
}

func MakePolicyItems() AParamSet {
	typemap := maketypemap()
	return AParamSet{
		"ASPENTRYSPEED":               AParam{Name: "ASPENTRYSPEED", Type: typemap[Float64], Desc: "allows slow moves into liquids"},
		"ASPREFERENCE":                AParam{Name: "ASPREFERENCE", Type: typemap[Int], Desc: "where to be when aspirating: 0 well bottom, 1 well top, 2 liquid level (if known)"},
		"ASPSPEED":                    AParam{Name: "ASPSPEED", Type: typemap[Float64], Desc: "aspirate pipetting rate"},
		"ASPXOFFSET":                  AParam{Name: "ASPXOFFSET", Type: typemap[Float64], Desc: "mm east of well target when aspirating"},
		"ASPYOFFSET":                  AParam{Name: "ASPYOFFSET", Type: typemap[Float64], Desc: "mm south of well target  when aspirating"},
		"ASPZOFFSET":                  AParam{Name: "ASPZOFFSET", Type: typemap[Float64], Desc: "mm above ASPREFERENCE when aspirating"},
		"ASP_WAIT":                    AParam{Name: "ASP_WAIT", Type: typemap[Float64], Desc: "wait time in seconds post aspirate"},
		"BLOWOUTOFFSET":               AParam{Name: "BLOWOUTOFFSET", Type: typemap[Float64], Desc: "mm above BLOWOUTREFERENCE"},
		"BLOWOUTREFERENCE":            AParam{Name: "BLOWOUTREFERENCE", Type: typemap[Int], Desc: "where to be when blowing out: 0 well bottom"},
		"BLOWOUTVOLUME":               AParam{Name: "BLOWOUTVOLUME", Type: typemap[Float64], Desc: "how much to blow out"},
		"BLOWOUTVOLUMEUNIT":           AParam{Name: "BLOWOUTVOLUMEUNIT", Type: typemap[String], Desc: "volume unit for blowout volume"},
		"CAN_MULTI":                   AParam{Name: "CAN_MULTI", Type: typemap[Bool], Desc: "is multichannel operation allowed?"},
		"DSPENTRYSPEED":               AParam{Name: "DSPENTRYSPEED", Type: typemap[Float64], Desc: "allows slow moves into liquids"},
		"DSPREFERENCE":                AParam{Name: "DSPREFERENCE", Type: typemap[Int], Desc: "where to be when dispensing: 0 well bottom, 1 well top, 2 liquid level (if known)"},
		"DSPSPEED":                    AParam{Name: "DSPSPEED", Type: typemap[Float64], Desc: "dispense pipetting rate"},
		"DSPXOFFSET":                  AParam{Name: "DSPXOFFSET", Type: typemap[Float64], Desc: "mm east of well target when dispensing"},
		"DSPYOFFSET":                  AParam{Name: "DSPYOFFSET", Type: typemap[Float64], Desc: "mm south of well target  when dispensing"},
		"DSPZOFFSET":                  AParam{Name: "DSPZOFFSET", Type: typemap[Float64], Desc: "mm above DSPREFERENCE"},
		"DSP_WAIT":                    AParam{Name: "DSP_WAIT", Type: typemap[Float64], Desc: "wait time in seconds post dispense"},
		"EXTRA_ASP_VOLUME":            AParam{Name: "EXTRA_ASP_VOLUME", Type: typemap[Volume], Desc: "additional volume to take up when aspirating"},
		"EXTRA_DISP_VOLUME":           AParam{Name: "EXTRA_DISP_VOLUME", Type: typemap[Volume], Desc: "additional volume to dispense"},
		"JUSTBLOWOUT":                 AParam{Name: "JUSTBLOWOUT", Type: typemap[Bool], Desc: "shortcut to get single transfer"},
		"OFFSETZADJUST":               AParam{Name: "OFFSETZADJUST", Type: typemap[Float64], Desc: "Added to z offset"},
		"OVERRIDEPIPETTESPEED":        AParam{Name: "OVERRIDEPIPETTESPEED", Type: typemap[Bool], Desc: "If true, out of range values will be set to the nearest acceptable value. If false, out of range values will generate errors"},
		"POST_MIX":                    AParam{Name: "POST_MIX", Type: typemap[Int], Desc: "number of mix cycles to do after dispense"},
		"MIX_VOLUME_OVERRIDE_TIP_MAX": AParam{Name: "MIX_VOLUME_OVERRIDE_TIP_MAX", Type: typemap[Bool], Desc: "Default to using the maximum volume for the current tip type if the specified post mix volume is too high"},
		"POST_MIX_RATE":               AParam{Name: "POST_MIX_RATE", Type: typemap[Float64], Desc: "pipetting rate when post mixing"},
		"POST_MIX_VOLUME":             AParam{Name: "POST_MIX_VOLUME", Type: typemap[Float64], Desc: "volume to post mix (ul)"},
		"POST_MIX_X":                  AParam{Name: "POST_MIX_X", Type: typemap[Float64], Desc: "x offset from centre of well (mm) when post-mixing"},
		"POST_MIX_Y":                  AParam{Name: "POST_MIX_Y", Type: typemap[Float64], Desc: "y offset from centre of well (mm) when post-mixing"},
		"POST_MIX_Z":                  AParam{Name: "POST_MIX_Z", Type: typemap[Float64], Desc: "z offset from centre of well (mm) when post-mixing"},
		"PRE_MIX":                     AParam{Name: "PRE_MIX", Type: typemap[Int], Desc: "number of mix cycles to do before aspirating"},
		"PRE_MIX_RATE":                AParam{Name: "PRE_MIX_RATE", Type: typemap[Float64], Desc: "pipetting rate when pre mixing"},
		"PRE_MIX_VOLUME":              AParam{Name: "PRE_MIX_VOLUME", Type: typemap[Float64], Desc: "volume to pre mix (ul)"},
		"PRE_MIX_X":                   AParam{Name: "PRE_MIX_X", Type: typemap[Float64], Desc: "x offset from centre of well (mm) when pre-mixing"},
		"PRE_MIX_Y":                   AParam{Name: "PRE_MIX_Y", Type: typemap[Float64], Desc: "y offset from centre of well (mm) when pre-mixing"},
		"PRE_MIX_Z":                   AParam{Name: "PRE_MIX_Z", Type: typemap[Float64], Desc: "z offset from centre of well (mm) when pre-mixing"},
		"RESET_OVERRIDE":              AParam{Name: "RESET_OVERRIDE", Type: typemap[Bool], Desc: "Do not generate reset commands"},
		"TIP_REUSE_LIMIT":             AParam{Name: "TIP_REUSE_LIMIT", Type: typemap[Int], Desc: "number of times tips can be reused for asp/dsp cycles"},
		"TOUCHOFF":                    AParam{Name: "TOUCHOFF", Type: typemap[Bool], Desc: "whether to move to TOUCHOFFSET after dispense"},
		"TOUCHOFFSET":                 AParam{Name: "TOUCHOFFSET", Type: typemap[Float64], Desc: "mm above wb to touch off at"},
		DefaultPipetteSpeed:           AParam{Name: DefaultPipetteSpeed, Type: typemap[Float64], Desc: "Default pipette speed in ml/min"},
		"DEFAULTZSPEED":               AParam{Name: "DEFAULTZSPEED", Type: typemap[Float64], Desc: "Default z movement speed in mm/s"},
		"PTZOFFSET":                   AParam{Name: "PTZOFFSET", Type: typemap[Float64], Desc: "Z offset for pistons to zero"},
		"PTZREFERENCE":                AParam{Name: "PTZREFERENCE", Type: typemap[Int], Desc: "Well reference for piston to zero: 0 = well bottom, 1 = well top, 2 = liquid level"},
		"CAN_SDD":                     AParam{Name: "CAN_SDD", Type: typemap[Bool], Desc: "Is it permissible just to do a one-shot transfer"},
		"MANUALPTZ":                   AParam{Name: "MANUALPTZ", Type: typemap[Bool], Desc: "Is explicit piston reset required? "},
		"DONT_BE_DIRTY":               AParam{Name: "DONT_BE_DIRTY", Type: typemap[Bool], Desc: "Don't switch this off"},
		"NO_AIR_DISPENSE":             AParam{Name: "NO_AIR_DISPENSE", Type: typemap[Bool], Desc: "Prevent dispensing anywhere other than under liquid?"},
		"CAN_MSA":                     AParam{Name: "CAN_MSA", Type: typemap[Bool], Desc: "Permissible to aspirate from multiple sources? -- currently non functional"},
		"DESCRIPTION":                 AParam{Name: "DESCRIPTION", Type: typemap[String], Desc: "Summary of LHPolicy to present to the user"},
		PolicyNameField:               AParam{Name: PolicyNameField, Type: typemap[String], Desc: "Name of the Liquid Policy"},
	}
}

func GetLHPolicyOptions() AParamSet {
	tm := maketypemap()

	return AParamSet{
		"USE_DRIVER_TIP_TRACKING": AParam{Name: "USE_DRIVER_TIP_TRACKING", Type: tm[Bool], Desc: "If driver has the option to use its own tip tracking, do so"},
	}
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

type Kind uint8

const (
	Float64 Kind = iota
	Int
	String
	Volume
	Bool
	_Max_Kind
)

func maketypemap() map[Kind]reflect.Type {
	// prototypical types for map
	var f float64
	var i int
	var s string
	var v wunit.Volume
	var b bool

	ret := map[Kind]reflect.Type{
		Float64: reflect.TypeOf(f),
		Int:     reflect.TypeOf(i),
		String:  reflect.TypeOf(s),
		Volume:  reflect.TypeOf(v),
		Bool:    reflect.TypeOf(b),
	}

	if len(ret) != int(_Max_Kind) {
		panic("We forgot to add full support for a type!")
	}

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
