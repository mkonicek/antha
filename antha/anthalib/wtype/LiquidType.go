package wtype

import (
	//"fmt"
	"fmt"
	"reflect"
	"strings"

	//"github.com/antha-lang/antha/antha/anthalib/wtype/liquidtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type PolicyName string

func (l PolicyName) String() string {
	return string(l)
}

func PolicyNameFromString(s string) PolicyName {
	return PolicyName(s)
}

type LiquidType LHPolicy // string // liqType

func (l LiquidType) String() (PolicyName, error) {
	return LiquidTypeName(l)
}

type liqType struct {
	Nm     string
	Policy *LHPolicy
	Rules  *map[string]LHPolicyRule // map using LHPolicyRule.Name as key, not a slice in order to be consistent with LHPolicyRuleSet
}

var DefaultLHPolicy = LTDefault

var (
	/*
		LTNIL             = LiquidType{Nm: "nil"}
		LTWater           = LiquidType{Nm: "water"}
		LTDefault         = LiquidType{Nm: "default"}
		LTCulture         = LiquidType{Nm: "culture"}
		LTProtoplasts     = LiquidType{Nm: "protoplasts"}
		LTDNA             = LiquidType{Nm: "dna"}
		LTDNAMIX          = LiquidType{Nm: "dna_mix"}
		LTProtein         = LiquidType{Nm: "protein"}
		LTMultiWater      = LiquidType{Nm: "multiwater"}
		LTLoad            = LiquidType{Nm: "load"}
		LTVISCOUS         = LiquidType{Nm: "viscous"}
		LTPEG             = LiquidType{Nm: "peg"}
		LTPAINT           = LiquidType{Nm: "paint"}
		LTNeedToMix       = LiquidType{Nm: "NeedToMix"}
		LTPostMix         = LiquidType{Nm: "PostMix"}
		LTload            = LiquidType{Nm: "load"}
		LTGlycerol        = LiquidType{Nm: "glycerol"}
		LTPLATEOUT        = LiquidType{Nm: "plateout"}
		LTDetergent       = LiquidType{Nm: "detergent"}
		LTCOLONY          = LiquidType{Nm: "colony"}
		LTNSrc            = LiquidType{Nm: "nitrogen_source"}
		InvalidPolicyName = LiquidType{Nm: "InvalidPolicyName"}
		LTSmartMix        = LiquidType{Nm: "SmartMix"}
		LTPreMix          = LiquidType{Nm: "PreMix"}
		LTDISPENSEABOVE   = LiquidType{Nm: "DispenseAboveLiquid"}
		LTMegaMix         = LiquidType{Nm: "MegaMix"}
		LTDoNotMix        = LiquidType{Nm: "DoNotMix"}
		LTDNACELLSMIX     = LiquidType{Nm: "dna_cells_mix"}
		LTloadwater       = LiquidType{Nm: "loadwater"}
	*/
	//
	InvalidPolicyName = LiquidType{}
	LTNIL             = LiquidType{}
	LTDefault         = LiquidType(MakeDefaultPolicy())
	LTSmartMix        = LiquidType(SmartMixPolicy())
	LTWater           = LiquidType(MakeWaterPolicy())
	LTMultiWater      = LiquidType(MakeMultiWaterPolicy())
	LTCulture         = LiquidType(MakeCulturePolicy())
	LTCultureReuse    = LiquidType(MakeCultureReusePolicy())
	LTGlycerol        = LiquidType(MakeGlycerolPolicy())
	LTSolvent         = LiquidType(MakeSolventPolicy())
	LtDefault         = LiquidType(MakeDefaultPolicy())
	LTDNA             = LiquidType(MakeDNAPolicy())
	LTDoNotMix        = LiquidType(MakeDefaultPolicy())
	LTNeedToMix       = LiquidType(MakeNeedToMixPolicy())
	LTPreMix          = LiquidType(PreMixPolicy())
	LTPostMix         = LiquidType(PostMixPolicy())
	LTMegaMix         = LiquidType(MegaMixPolicy())
	LTVISCOUS         = LiquidType(MakeViscousPolicy())
	LTPAINT           = LiquidType(MakePaintPolicy())

	// pols["lysate"] = MakeLysatePolicy()
	LTProtein            = LiquidType(MakeProteinPolicy())
	LTDetergent          = LiquidType(MakeDetergentPolicy())
	LTload               = LiquidType(MakeLoadPolicy())
	LTloadwater          = LiquidType(MakeLoadWaterPolicy())
	LTDISPENSEABOVE      = LiquidType(MakeDispenseAboveLiquidPolicy())
	LTDISPENSEABOVEMULTI = LiquidType(MakeDispenseAboveLiquidMultiPolicy())
	LTPEG                = LiquidType(MakePEGPolicy())
	LTProtoplasts        = LiquidType(MakeProtoplastPolicy())
	LTDNAMIX             = LiquidType(MakeDNAMixPolicy())
	LTDNAMIXMULTI        = LiquidType(MakeDNAMixMultiPolicy())
	LTDNACELLSMIX        = LiquidType(MakeDNACELLSMixPolicy())
	LTDNACELLSMIXMULTI   = LiquidType(MakeDNACELLSMixMultiPolicy())
	LTPLATEOUT           = LiquidType(MakePlateOutPolicy())
	LTCOLONY             = LiquidType(MakeColonyPolicy())
	LTCOLONYMIX          = LiquidType(MakeColonyMixPolicy())
	//      pols["lysate"] = MakeLysatePolicy()
	LTCSrc = LiquidType(MakeCarbonSourcePolicy())
	LTNSrc = LiquidType(MakeNitrogenSourcePolicy())

/*
	LTNIL = nil
	LTLoad, _            = GetPolicyByName("load")
	LTVISCOUS, _         = GetPolicyByName("viscous")
	LTPEG, _             = GetPolicyByName("peg")
	LTPAINT, _           = GetPolicyByName("paint")
	LTNeedToMix, _       = GetPolicyByName("NeedToMix")
	LTPostMix, _         = GetPolicyByName("PostMix")
	LTload, _            = GetPolicyByName("load")
	LTGlycerol, _        = GetPolicyByName("glycerol")
	LTPLATEOUT, _        = GetPolicyByName("plateout")
	LTDetergent, _       = GetPolicyByName("detergent")
	LTCOLONY, _          = GetPolicyByName("colony")
	LTNSrc, _            = GetPolicyByName("nitrogen_source")
	InvalidPolicyName, _ = GetPolicyByName("InvalidPolicyName")
	LTSmartMix, _        = GetPolicyByName("SmartMix")
	LTPreMix, _          = GetPolicyByName("PreMix")
	LTDISPENSEABOVE, _   = GetPolicyByName("DispenseAboveLiquid")
	LTMegaMix, _         = GetPolicyByName("MegaMix")
	LTDoNotMix, _        = GetPolicyByName("DoNotMix")
	LTDNACELLSMIX, _     = GetPolicyByName("dna_cells_mix")
	LTloadwater, _       = GetPolicyByName("loadwater")
*/
//

/*LTGlycerol
LTEthanol
LTDetergent
LTCulture
LTProtein
LTDNA
LTload


LTNeedToMix

LTPostMix
LTVISCOUS
LTPAINT

LTDISPENSEABOVEMULTI
LTPEG
LTProtoplasts
LTCulutureReuse
LTDNAMIX
LTDNAMIXMULTI
LTPLATEOUT
LTCOLONY
LTCOLONYMIX
LTDNACELLSMIX
LTDNACELLSMIXMULTI
LTMultiWater
LTCSrc
LTNSrc

LTSolvent
*/
)

func LiquidTypeFromString(s PolicyName) (LiquidType, error) {
	lt, err := GetPolicyByName(s)

	return LiquidType(lt), err
	/*
		match, number := liquidtype.LiquidTypeFromPolicyDOE(s.String())

		if match {
			panic("invalid policy string")
			return LiquidType(number), nil
		}

		switch s {
		case "water":
			return LTWater, nil
		case "":
			return LTWater, fmt.Errorf("no liquid policy specified so using default water policy")
		case "glycerol":
			return LTGlycerol, nil
		case "ethanol":
			return LTEthanol, nil
		case "detergent":
			return LTDetergent, nil
		case "culture":
			return LTCulture, nil
		case "culturereuse":
			return LTCulutureReuse, nil
		case "protein":
			return LTProtein, nil
		case "dna":
			return LTDNA, nil
		case "load":
			return LTload, nil
		case "DoNotMix":
			return LTDoNotMix, nil
		case "loadwater":
			return LTloadwater, nil
		case "NeedToMix":
			return LTNeedToMix, nil
		case "PreMix":
			return LTPreMix, nil
		case "PostMix":
			return LTPostMix, nil
		case "viscous":
			return LTVISCOUS, nil
		case "Paint":
			return LTPAINT, nil
		case "DispenseAboveLiquid":
			return LTDISPENSEABOVE, nil
		case "DispenseAboveLiquidMulti":
			return LTDISPENSEABOVEMULTI, nil
		case "PEG":
			return LTPEG, nil
		case "Protoplasts":
			return LTProtoplasts, nil
		case "dna_mix":
			return LTDNAMIX, nil
		case "dna_mix_multi":
			return LTDNAMIXMULTI, nil
		case "plateout":
			return LTPLATEOUT, nil
		case "colony":
			return LTCOLONY, nil
		case "colonymix":
			return LTCOLONYMIX, nil
		case "dna_cells_mix":
			return LTDNACELLSMIX, nil
		case "dna_cells_mix_multi":
			return LTDNACELLSMIXMULTI, nil
		case "multiwater":
			return LTMultiWater, nil
		case "carbon_source":
			return LTCSrc, nil
		case "nitrogen_source":
			return LTNSrc, nil
		case "MegaMix":
			return LTMegaMix, nil
		case "solvent":
			return LTSolvent, nil
		case "SmartMix":
			return LTSmartMix, nil
		case "default":
			return DefaultLHPolicy, nil
		default:
			panic("invalid policy string")
			return DefaultLHPolicy, fmt.Errorf("no liquid policy found for " + s.String() + " so using default policy")
		}
	*/
}

func LiquidTypeName(lt LiquidType) (PolicyName, error) {

	name, ok := lt[PolicyNameField]

	if !ok {
		return "", fmt.Errorf("Policy has no name")
	}

	str, ok := name.(string)

	if !ok {
		return PolicyName(fmt.Sprint(name)), fmt.Errorf("policy name %v is not string", name)
	}

	return PolicyName(str), nil
	/*
				match, str := liquidtype.StringFromLiquidTypeNumber(int(lt))

				if match {
					panic("invalid policy int")
					return PolicyName(str), nil
				}

		<<<<<<< HEAD
				switch lt {
				case LTWater:
					return "water", nil
				case LTGlycerol:
					return "glycerol", nil
				case LTEthanol:
					return "ethanol", nil
				case LTDetergent:
					return "detergent", nil
				case LTCulture:
					return "culture", nil
				case LTCulutureReuse:
					return "culturereuse", nil
				case LTProtein:
					return "protein", nil
				case LTDNA:
					return "dna", nil
				case LTload:
					return "load", nil
				case LTDoNotMix:
					return "DoNotMix", nil
				case LTloadwater:
					return "loadwater", nil
				case LTNeedToMix:
					return "NeedToMix", nil
				case LTPreMix:
					return "PreMix", nil
				case LTPostMix:
					return "PostMix", nil
				case LTPAINT:
					return "Paint", nil
				case LTDISPENSEABOVE:
					return "DispenseAboveLiquid", nil
				case LTProtoplasts:
					return "Protoplasts", nil
				case LTPEG:
					return "PEG", nil
				case LTDNAMIX:
					return "dna_mix", nil
				case LTPLATEOUT:
					return "plateout", nil
				case LTCOLONY:
					return "colony", nil
				case LTCOLONYMIX:
					return "colonymix", nil
				case LTDNACELLSMIX:
					return "dna_cells_mix", nil
				case LTMultiWater:
					return "multiwater", nil
				case LTCSrc:
					return "carbon_source", nil
				case LTNSrc:
					return "nitrogen_source", nil
				case LTMegaMix:
					return "MegaMix", nil
				case LTSmartMix:
					return "SmartMix", nil
				default:
					panic("invalid policy int")
					return "nil", fmt.Errorf("no policy %v", lt)
				}
	*/
}

func mergeSolubilities(c1, c2 *LHComponent) float64 {
	if c1.Smax < c2.Smax {
		return c1.Smax
	}

	return c2.Smax
}

// helper functions... will need extending eventually

func mergeTypes(c1, c2 *LHComponent) LiquidType {
	// couple of mixing rules: protein, dna etc. are basically
	// special water so we retain that characteristic whatever happens
	// ditto culture... otherwise we look for the majority
	// what we do for protein-dna mixtures I'm not sure! :)

	// nil type is overridden

	if reflect.DeepEqual(c1.Type, LiquidType(LTNIL)) {
		return c2.Type
	} else if reflect.DeepEqual(c2.Type, LTNIL) {
		return c1.Type
	}

	if reflect.DeepEqual(c1.Type, LTCulture) || reflect.DeepEqual(c2.Type, LTCulture) {
		return LiquidType(LTCulture)
	} else if reflect.DeepEqual(c1.Type, LTProtoplasts) || reflect.DeepEqual(c2.Type, LTProtoplasts) {
		return LiquidType(LTProtoplasts)
	} else if reflect.DeepEqual(c1.Type, LTDNA) || reflect.DeepEqual(c2.Type, LTDNA) || reflect.DeepEqual(c1.Type, LTDNAMIX) || reflect.DeepEqual(c2.Type, LTDNAMIX) {
		return LiquidType(LTDNA)
	} else if reflect.DeepEqual(c1.Type, LTProtein) || reflect.DeepEqual(c2.Type, LTProtein) {
		return LiquidType(LTProtein)
	}
	v1 := wunit.NewVolume(c1.Vol, c1.Vunit)
	v2 := wunit.NewVolume(c2.Vol, c2.Vunit)

	if v1.LessThan(&v2) {
		return c2.Type
	}

	return c1.Type
}

// merge two names... we have a lookup function to add here
func mergeNames(a, b string) string {
	tx := strings.Split(a, "+")
	tx2 := strings.Split(b, "+")

	tx3 := mergeTox(tx, tx2)

	tx3 = Normalize(tx3)

	return strings.Join(tx3, "+")
}

// very simple, just add elements of tx2 not already in tx
func mergeTox(tx, tx2 []string) []string {
	for _, v := range tx2 {
		ix := IndexOfStringArray(v, tx)

		if ix == -1 {
			tx = append(tx, v)
		}
	}

	return tx
}

func IndexOfStringArray(s string, a []string) int {
	ret := -1
	for i, v := range a {
		if v == s {
			ret = i
			break
		}
	}
	return ret
}

// TODO -- fill in some normalizations
// - water + salt = saline might be an eg
// but unlikely to be useful
func Normalize(tx []string) []string {
	if tx[0] == "" && len(tx) > 1 {
		return tx[1:len(tx)]
	}
	return tx
}
