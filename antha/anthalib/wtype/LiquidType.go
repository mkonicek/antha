package wtype

import "github.com/antha-lang/antha/antha/anthalib/wunit"

// PolicyName represents the name of a liquid handling policy
// used to look up the details of that policy.
type PolicyName string

func (l PolicyName) String() string {
	return string(l)
}

func PolicyNameFromString(s string) PolicyName {
	return PolicyName(s)
}

// LiquidType represents the type of a Liquid
type LiquidType string

func (l LiquidType) String() (PolicyName, error) {
	return LiquidTypeName(l)
}

// DefaultLHPolicy which will be used if no policy is set.
const DefaultLHPolicy = LTDefault

// Valid default LiquidTypes
const (
	LTNIL                   LiquidType = "nil"
	LTWater                 LiquidType = "water"
	LTDefault               LiquidType = "default"
	LTCulture               LiquidType = "culture"
	LTProtoplasts           LiquidType = "protoplasts"
	LTDNA                   LiquidType = "dna"
	LTDNAMIX                LiquidType = "dna_mix"
	LTProtein               LiquidType = "protein"
	LTMultiWater            LiquidType = "multiwater"
	LTLoad                  LiquidType = "load"
	LTVISCOUS               LiquidType = "viscous"
	LTPEG                   LiquidType = "peg"
	LTPAINT                 LiquidType = "paint"
	LTNeedToMix             LiquidType = "NeedToMix"
	LTPostMix               LiquidType = "PostMix"
	LTload                  LiquidType = "load"
	LTGlycerol              LiquidType = "glycerol"
	LTPLATEOUT              LiquidType = "plateout"
	LTDetergent             LiquidType = "detergent"
	LTCOLONY                LiquidType = "colony"
	LTNSrc                  LiquidType = "nitrogen_source"
	InvalidPolicyName       LiquidType = "InvalidPolicyName"
	LTEthanol               LiquidType = "ethanol"
	LTDoNotMix              LiquidType = "DoNotMix"
	LTloadwater             LiquidType = "loadwater"
	LTPreMix                LiquidType = "PreMix"
	LTDISPENSEABOVE         LiquidType = "DispenseAboveLiquid"
	LTDISPENSEABOVEMULTI    LiquidType = "DispenseAboveLiquidMulti"
	LTCulutureReuse         LiquidType = "culturereuse"
	LTDNAMIXMULTI           LiquidType = "dna_mix_multi"
	LTCOLONYMIX             LiquidType = "colonymix"
	LTDNACELLSMIX           LiquidType = "dna_cells_mix"
	LTDNACELLSMIXMULTI      LiquidType = "dna_cells_mix_multi"
	LTCSrc                  LiquidType = "carbon_source"
	LTMegaMix               LiquidType = "MegaMix"
	LTSolvent               LiquidType = "solvent"
	LTSmartMix              LiquidType = "SmartMix"
	LTSingleChannel         LiquidType = "SingleChannel"
	LTSmartMixSingleChannel LiquidType = "SmartMixSingleChannel"
	LTLiquidLevel           LiquidType = "LiquidLevel"
	LTSmartMixLiquidLevel   LiquidType = "SmartMixLiquidLevel"
)

// LiquidTypeFromString returns a LiquidType from a PolicyName
// If the PolicyName is invalid and the DoNotPermitCustomPolicies option is used as an argument then an error is returned.
// By default, custom policyNames may be added and the validity of these will be checked later when robot instructions are generated, rather than in the element.
func LiquidTypeFromString(s PolicyName, options ...PolicyOption) (LiquidType, error) {
	_, err := GetPolicyByName(PolicyName(s))
	for _, option := range options {
		if string(option) == string(DoNotPermitCustomPolicies) {
			return LiquidType(s), err
		}
	}
	return LiquidType(s), nil
}

// LiquidTypeName returns a PolicyName from a LiquidType
func LiquidTypeName(lt LiquidType) (PolicyName, error) {
	_, err := GetPolicyByName(PolicyName(lt))
	return PolicyName(lt), err
}

func mergeSolubilities(c1, c2 *Liquid) float64 {
	if c1.Smax < c2.Smax {
		return c1.Smax
	}

	return c2.Smax
}

// helper functions... will need extending eventually

func mergeTypes(c1, c2 *Liquid) LiquidType {
	// couple of mixing rules: protein, dna etc. are basically
	// special water so we retain that characteristic whatever happens
	// ditto culture... otherwise we look for the majority
	// what we do for protein-dna mixtures I'm not sure! :)

	// nil type is overridden

	if c1.Type == LTNIL {
		return c2.Type
	} else if c2.Type == LTNIL {
		return c1.Type
	}

	if c1.Type == LTCulture || c2.Type == LTCulture {
		return LTCulture
	} else if c1.Type == LTProtoplasts || c2.Type == LTProtoplasts {
		return LTProtoplasts
	} else if c1.Type == LTDNA || c2.Type == LTDNA || c1.Type == LTDNAMIX || c2.Type == LTDNAMIX {
		return LTDNA
	} else if c1.Type == LTProtein || c2.Type == LTProtein {
		return LTProtein
	}
	v1 := wunit.NewVolume(c1.Vol, c1.Vunit)
	v2 := wunit.NewVolume(c2.Vol, c2.Vunit)

	if v1.LessThan(&v2) {
		return c2.Type
	}

	return c1.Type
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
		return tx[1:]
	}
	return tx
}
