// liquidhandling/lhtypes.Go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.Org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

// defines types for dealing with liquid handling requests
package wtype

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"strings"
	//"github.com/antha-lang/antha/microArch/logger"
	"github.com/antha-lang/antha/graph"
)

// structure describing a liquid component and its desired properties
type LHComponent struct {
	ID                 string
	BlockID            BlockID
	DaughterID         string
	ParentID           string
	Inst               string
	Order              int
	CName              string
	Type               LiquidType
	Vol                float64
	Conc               float64
	Vunit              string
	Cunit              string
	Tvol               float64
	Smax               float64
	Visc               float64
	StockConcentration float64
	Extra              map[string]interface{}
	Loc                string
	Destination        string
}

func (lhc *LHComponent) Generation() int {
	g, ok := lhc.Extra["Generation"]

	if ok {
		return g.(int)
	}

	return 0
}

func (lhc *LHComponent) SetGeneration(i int) {
	lhc.Extra["Generation"] = i
}

func (lhc *LHComponent) IsZero() bool {
	if lhc == nil || lhc.Type == LTNIL || lhc.CName == "" || lhc.Vol < 0.0000001 {
		return true
	}
	return false
}

func (lhc *LHComponent) SetVolume(v wunit.Volume) {
	lhc.Vol = v.RawValue()
	lhc.Vunit = v.Unit().PrefixedSymbol()
}

func (lhc *LHComponent) HasParent(s string) bool {
	return strings.Contains(lhc.ParentID, s)
}

func (lhc *LHComponent) HasDaughter(s string) bool {
	return strings.Contains(lhc.DaughterID, s)
}

func (lhc *LHComponent) Name() string {
	return lhc.CName
}

func (lhc *LHComponent) TypeName() string {
	return LiquidTypeName(lhc.Type)
}

func (lhc *LHComponent) Volume() wunit.Volume {
	return wunit.NewVolume(lhc.Vol, lhc.Vunit)
}

func (lhc *LHComponent) Remove(v wunit.Volume) {
	///TODO -- catch errors
	lhc.Vol -= v.ConvertToString(lhc.Vunit)

	if lhc.Vol < 0.0 {
		lhc.Vol = 0.0
	}
}

func (lhc *LHComponent) Dup() *LHComponent {
	c := NewLHComponent()
	c.ID = lhc.ID
	c.Order = lhc.Order
	c.CName = lhc.CName
	c.Type = lhc.Type
	c.Vol = lhc.Vol
	c.Conc = lhc.Conc
	c.Vunit = lhc.Vunit
	c.Tvol = lhc.Tvol
	c.Smax = lhc.Smax
	c.Visc = lhc.Visc
	c.StockConcentration = lhc.StockConcentration
	c.Extra = make(map[string]interface{}, len(lhc.Extra))
	for k, v := range lhc.Extra {
		c.Extra[k] = v
	}
	c.Loc = lhc.Loc
	c.Destination = lhc.Destination
	c.ParentID = lhc.ParentID
	c.DaughterID = lhc.DaughterID
	return c
}

func (cmp *LHComponent) SetSample(flag bool) bool {
	if cmp == nil {
		return false
	}

	if cmp.Extra == nil {
		cmp.Extra = make(map[string]interface{})
	}

	cmp.Extra["IsSample"] = flag

	return true
}

func (cmp *LHComponent) IsSample() bool {
	if cmp == nil {
		return false
	}

	f, ok := cmp.Extra["IsSample"]

	if !ok || !f.(bool) {
		return false
	}

	return true
}

func (cmp *LHComponent) HasAnyParent() bool {
	if cmp.ParentID != "" {
		return true
	}

	return false
}

func (cmp *LHComponent) AddParentComponent(cmp2 *LHComponent) {
	if cmp.ParentID != "" {
		cmp.ParentID += "_"
	}
	cmp.ParentID += cmp2.String() + "(" + cmp2.ParentID + ")"
	//cmp.ParentID += cmp2.ID + "(" + cmp2.ParentID + ")"
}

func (cmp *LHComponent) AddDaughterComponent(cmp2 *LHComponent) {
	if cmp.DaughterID != "" {
		cmp.DaughterID += "_"
	}

	//cmp.DaughterID += cmp2.String()

	cmp.DaughterID += cmp2.ID
}

func (cmp *LHComponent) MixPreserveTvol(cmp2 *LHComponent) {
	cmp.Mix(cmp2)
	if cmp2.Vol == 0.00 && cmp2.Tvol > 0.00 {
		vcmp := wunit.NewVolume(cmp.Vol, cmp.Vunit)
		vcmp2 := wunit.NewVolume(cmp2.Tvol, cmp2.Vunit)
		vcmp.SetValue(vcmp2.ConvertToString(cmp.Vunit))
		cmp.Vol = vcmp.RawValue() // same units
		cmp.Tvol = vcmp2.ConvertToString(cmp.Vunit)
	} else if cmp.Tvol > 0.00 {
		cmp.Vol = cmp.Tvol
	}
}

func (cmp *LHComponent) Mix(cmp2 *LHComponent) {
	//wasEmpty := cmp.IsZero()
	cmp.Smax = mergeSolubilities(cmp, cmp2)
	// determine type of final
	cmp.Type = mergeTypes(cmp, cmp2)
	// add cmp2 to cmp

	vcmp := wunit.NewVolume(cmp.Vol, cmp.Vunit)
	vcmp2 := wunit.NewVolume(cmp2.Vol, cmp2.Vunit)
	vcmp.Add(vcmp2)
	cmp.Vol = vcmp.RawValue() // same units
	cmp.CName = mergeNames(cmp.CName, cmp2.CName)

	// allow trace back
	//logger.Track(fmt.Sprintf("MIX %s %s %s", cmp.ID, cmp2.ID, vcmp.ToString()))

	// add parent IDs... this is recursive

	/*
		if !wasEmpty {
			cmp.AddParentComponent(cmp)
		}
	*/
	cmp.AddParentComponent(cmp2)
	//	cmp.ID = "component-" + GetUUID()
	cmp.ID = GetUUID()
	cmp2.AddDaughterComponent(cmp)
}

// @implement Liquid
// @deprecate Liquid

func (lhc *LHComponent) GetSmax() float64 {
	return lhc.Smax
}

func (lhc *LHComponent) GetVisc() float64 {
	return lhc.Visc
}

func (lhc *LHComponent) GetExtra() map[string]interface{} {
	return lhc.Extra
}

func (lhc *LHComponent) GetConc() float64 {
	return lhc.Conc
}

func (lhc *LHComponent) GetCunit() string {
	return lhc.Cunit
}

// new
func (lhc *LHComponent) Concentration() (conc wunit.Concentration) {
	conc = wunit.NewConcentration(lhc.Conc, lhc.Cunit)
	return conc
}

func (lhc *LHComponent) GetVunit() string {
	return lhc.Vunit
}

func (lhc *LHComponent) GetType() string {
	return LiquidTypeName(lhc.Type)
}

func NewLHComponent() *LHComponent {
	var lhc LHComponent
	//lhc.ID = "component-" + GetUUID()
	lhc.ID = GetUUID()
	lhc.Vunit = "ul"
	lhc.Extra = make(map[string]interface{})
	return &lhc
}

func (cmp *LHComponent) String() string {
	id := cmp.ID

	l := cmp.Loc

	v := fmt.Sprintf("%-6.3f:%s", cmp.Vol, cmp.Vunit)

	if l == "" {
		l = "NOPLATE:NOWELL"
	}

	return id + ":" + l + ":" + v
}

func (cmp *LHComponent) ParentTree() graph.StringGraph {
	g := graph.StringGraph{Nodes: make([]string, 0, 3), Outs: make(map[string][]string)}
	parseTree(cmp.ID+"("+cmp.ParentID+")", &g)
	return g
}

// graphviz format
func (cmp *LHComponent) ParentTreeString() string {
	g := cmp.ParentTree()
	s := graph.Print(graph.PrintOpt{Graph: &g})
	return s
}

//   a(b_c_d)_e()_f(g_h)
//   nodes: [abcdefgh]
//   outs : a:[] b[a] c[a] d[a] e[] f[] g[f] h[f]
//

func parseTree(p string, g *graph.StringGraph) []string {
	newnodes := make([]string, 0, 3)
	if p[0] == '(' {
		// strip brackets
		p = p[1 : len(p)-1]
	}

	if len(p) == 0 {
		// empty bracket pair
		return newnodes
	}

	bc := 0

	splits := make([]int, 0, len(p))

	for i, c := range p {
		if c == '(' {
			bc += 1
			continue
		} else if c == ')' {
			bc -= 1
			continue
		}

		if bc == 0 && c == '_' {
			splits = append(splits, i)
		}
		//   a(b()_c()_d())_e()_f(g()_h())
		//                 s   s
	}

	splits = append(splits, len(p))

	// carve up

	b := 0

	for _, e := range splits {
		tok := p[b:e]
		lb := strings.Index(tok, "(")
		node := tok[:lb]
		if !wutil.StrInStrArray(node, g.Nodes) {
			g.Nodes = append(g.Nodes, node)
			g.Outs[node] = make([]string, 0, 3)
			newnodes = append(newnodes, node)
		}

		childnodes := parseTree(tok[lb:], g)

		for _, child := range childnodes {
			g.Outs[child] = append(g.Outs[child], node)
		}
		b = e + 1
	}

	return newnodes
}

func (lhc *LHComponent) AddVolumeRule(minvol, maxvol float64, pol LHPolicy) error {
	lhpr, err := lhc.GetPolicies()

	if err != nil {
		return err
	}

	rulenum := len(lhpr.Rules)

	name := fmt.Sprintf("UserRule%d", rulenum+1)

	rule := NewLHPolicyRule(name)
	err = rule.AddNumericConditionOn("VOLUME", minvol, maxvol)
	if err != nil {
		return err
	}
	lhpr.AddRule(rule, pol)

	err = rule.AddCategoryConditionOn("INSTANCE", lhc.ID)
	if err != nil {
		return err
	}

	err = lhc.SetPolicies(lhpr)

	return err
}

func (lhc *LHComponent) AddPolicy(pol LHPolicy) error {
	lhpr, err := lhc.GetPolicies()

	if err != nil {
		return err
	}

	rulenum := len(lhpr.Rules)

	name := fmt.Sprintf("UserRule%d", rulenum+1)

	rule := NewLHPolicyRule(name)
	err = rule.AddCategoryConditionOn("INSTANCE", lhc.ID)
	if err != nil {
		return err
	}
	lhpr.AddRule(rule, pol)

	err = lhc.SetPolicies(lhpr)

	return err

}

// in future this will be deprecated... should not let user completely reset policies
func (lhc *LHComponent) SetPolicies(rs *LHPolicyRuleSet) error {
	buf, err := json.Marshal(rs)

	if err == nil {
		lhc.Extra["Policies"] = string(buf)
	}

	return err
}

func (lhc *LHComponent) GetPolicies() (*LHPolicyRuleSet, error) {
	var rs LHPolicyRuleSet
	var err error

	if lhc.Extra == nil {
		return NewLHPolicyRuleSet(), err
	}

	ent, ok := lhc.Extra["Policies"]

	if !ok {
		return NewLHPolicyRuleSet(), err
	}

	s, ok := ent.(string)

	if !ok {
		err = fmt.Errorf("Wrong type for policies entry (%v)", ent)
		return &rs, err
	}

	ba := []byte(s)

	err = json.Unmarshal(ba, &rs)
	return &rs, err
}

func (lhc *LHComponent) IsValuable() bool {
	if lhc.Extra == nil {
		return false
	}

	v, ok := lhc.Extra["valuable"]

	if !ok {
		return false
	}

	b, ok := v.(bool)

	if !ok {
		return false
	}

	return b
}

func (lhc *LHComponent) SetValue(b bool) {
	if lhc.Extra == nil {
		lhc.Extra = make(map[string]interface{})
	}

	lhc.Extra["valuable"] = b
}
