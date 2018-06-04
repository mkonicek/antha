package liquidhandling

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/cache/plateCache"
	"github.com/antha-lang/antha/inventory/testinventory"
)

type Condition interface {
	ApplyTo(*wtype.LHPolicyRule)
}

type CategoryCondition struct {
	Attribute string
	Value     string
}

func (self *CategoryCondition) ApplyTo(rule *wtype.LHPolicyRule) {
	rule.AddCategoryConditionOn(self.Attribute, self.Value)
}

type NumericCondition struct { //nolint
	Attribute string
	Low       float64
	High      float64
}

func (self *NumericCondition) ApplyTo(rule *wtype.LHPolicyRule) {
	rule.AddNumericConditionOn(self.Attribute, self.Low, self.High)
}

type Rule struct {
	Name       string
	Conditions []Condition
	Policy     map[string]interface{}
}

func (self *Rule) AddToPolicy(pol *wtype.LHPolicyRuleSet) {
	rule := wtype.NewLHPolicyRule(self.Name)
	for _, c := range self.Conditions {
		c.ApplyTo(&rule)
	}

	policy := make(wtype.LHPolicy, len(self.Policy))
	for k, v := range self.Policy {
		policy[k] = v
	}

	pol.AddRule(rule, policy)
}

type InstructionAssertion struct {
	Instruction int
	Values      map[string]interface{}
}

func (self *InstructionAssertion) Assert(t *testing.T, ris []RobotInstruction, name string) {
	if self.Instruction < 0 || self.Instruction >= len(ris) {
		t.Errorf("%s: test error: assertion on instruction %d, but only %d instructions", name, self.Instruction, len(ris))
		return
	}
	ins := ris[self.Instruction]

	for param, e := range self.Values {
		if g := ins.GetParameter(param); !reflect.DeepEqual(e, g) {
			t.Errorf("%s: instruction %d parameter %s: expected %v, got %v", name, self.Instruction, param, e, g)
		}
	}

}

type PolicyTest struct {
	Name                 string
	Rules                []*Rule
	Instruction          RobotInstruction
	Robot                *LHProperties
	ExpectedInstructions string
	Assertions           []*InstructionAssertion
	Error                string
}

func stringInstructions(inss []RobotInstruction) string {
	s := make([]string, len(inss))
	for i, ins := range inss {
		s[i] = InstructionTypeName(ins)
	}
	return "[" + strings.Join(s, ",") + "]"
}

func (self *PolicyTest) Run(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())
	ctx = plateCache.NewContext(ctx)
	if self.Robot == nil {
		robot, err := makeTestGilsonWithPlates(ctx)
		if err != nil {
			err = errors.Wrap(err, self.Name)
			t.Fatal(err)
		}
		self.Robot = robot
	}

	policySet, err := wtype.GetLHPolicyForTest()
	if err != nil {
		err = errors.Wrap(err, self.Name)
		t.Fatal(err)
	}

	for _, rule := range self.Rules {
		rule.AddToPolicy(policySet)
	}

	set := NewRobotInstructionSet(self.Instruction)
	ris, err := set.Generate(ctx, policySet, self.Robot)
	if err != nil {
		if self.Error == "" {
			err = errors.Wrapf(err, "%s: unexpected error", self.Name)
			t.Error(err)
		} else if self.Error != err.Error() {
			t.Errorf("%s: errors don't match:\ne: \"%s\",\ng: \"%s\"", self.Name, self.Error, err.Error())
		}
		return
	}

	if self.Error != "" {
		t.Errorf("%s: error not generated: expected \"%s\"", self.Name, self.Error)
		return
	}

	if g := stringInstructions(ris); self.ExpectedInstructions != g {
		t.Errorf("%s: instruction types don't match\n  g: %s\n  e: %s", self.Name, g, self.ExpectedInstructions)
		return
	}

	for _, a := range self.Assertions {
		a.Assert(t, ris, self.Name)
	}
}

const (
	HVMinRate           = 0.225
	HVMaxRate           = 37.5
	LVMinRate           = 0.0225
	LVMaxRate           = 3.75
	defaultZSpeed       = 120.0
	defaultZOffset      = 0.5
	defaultPipetteSpeed = 3.0
)

func getHVConfig() *wtype.LHChannelParameter {
	minvol := wunit.NewVolume(10, "ul")
	maxvol := wunit.NewVolume(250, "ul")
	minspd := wunit.NewFlowRate(HVMinRate, "ml/min")
	maxspd := wunit.NewFlowRate(HVMaxRate, "ml/min")

	return wtype.NewLHChannelParameter("HVconfig", "GilsonPipetmax", minvol, maxvol, minspd, maxspd, 8, false, wtype.LHVChannel, 0)
}

func getLVConfig() *wtype.LHChannelParameter {
	newminvol := wunit.NewVolume(0.5, "ul")
	newmaxvol := wunit.NewVolume(20, "ul")
	newminspd := wunit.NewFlowRate(LVMinRate, "ml/min")
	newmaxspd := wunit.NewFlowRate(LVMaxRate, "ml/min")

	return wtype.NewLHChannelParameter("LVconfig", "GilsonPipetmax", newminvol, newmaxvol, newminspd, newmaxspd, 8, false, wtype.LHVChannel, 1)
}

func makeGilson() *LHProperties {
	// gilson pipetmax

	layout := make(map[string]wtype.Coordinates)
	i := 0
	x0 := 3.886
	y0 := 3.513
	z0 := -82.035
	xi := 149.86
	yi := 95.25
	xp := x0 // nolint
	yp := y0
	zp := z0
	for y := 0; y < 3; y++ {
		xp = x0
		for x := 0; x < 3; x++ {
			posname := fmt.Sprintf("position_%d", i+1)
			crds := wtype.Coordinates{X: xp, Y: yp, Z: zp}
			layout[posname] = crds
			i += 1
			xp += xi
		}
		yp += yi
	}
	lhp := NewLHProperties(9, "Pipetmax", "Gilson", LLLiquidHandler, DisposableTips, layout)
	// get tips permissible from the factory
	SetUpTipsFor(lhp)

	lhp.Tip_preferences = []string{"position_2", "position_3", "position_6", "position_9", "position_8", "position_5", "position_4", "position_7"}
	//lhp.Tip_preferences = []string{"position_2", "position_3", "position_6", "position_9", "position_8", "position_5", "position_7"}

	//lhp.Tip_preferences = []string{"position_9", "position_6", "position_3", "position_5", "position_2"} //jmanart i cut it down to 5, as it was hardcoded in the liquidhandler getInputs call before

	// original preferences
	lhp.Input_preferences = []string{"position_4", "position_5", "position_6", "position_9", "position_8", "position_3"}
	lhp.Output_preferences = []string{"position_8", "position_9", "position_6", "position_5", "position_3", "position_1"}

	// use these new preferences for gel loading: this is needed because outplate overlaps inplate otherwise so move inplate to position 5 rather than 4 (pos 4 deleted)
	//lhp.Input_preferences = []string{"position_5", "position_6", "position_9", "position_8", "position_3"}
	//lhp.Output_preferences = []string{"position_9", "position_8", "position_7", "position_6", "position_5", "position_3"}

	lhp.Wash_preferences = []string{"position_8"}
	lhp.Tipwaste_preferences = []string{"position_1", "position_7"}
	lhp.Waste_preferences = []string{"position_9"}
	//	lhp.Tip_preferences = []int{2, 3, 6, 9, 5, 8, 4, 7}
	//	lhp.Input_preferences = []int{24, 25, 26, 29, 28, 23}
	//	lhp.Output_preferences = []int{10, 11, 12, 13, 14, 15}
	hvconfig := getHVConfig()
	hvadaptor := wtype.NewLHAdaptor("DummyAdaptor", "Gilson", hvconfig)
	hvhead := wtype.NewLHHead("HVHead", "Gilson", hvconfig)
	hvhead.Adaptor = hvadaptor

	lvconfig := getLVConfig()
	lvadaptor := wtype.NewLHAdaptor("DummyAdaptor", "Gilson", lvconfig)
	lvhead := wtype.NewLHHead("LVHead", "Gilson", lvconfig)
	lvhead.Adaptor = lvadaptor

	lhp.Heads = append(lhp.Heads, hvhead)
	lhp.Heads = append(lhp.Heads, lvhead)
	lhp.HeadsLoaded = append(lhp.HeadsLoaded, hvhead)
	lhp.HeadsLoaded = append(lhp.HeadsLoaded, lvhead)

	return lhp
}

func makeTestGilsonWithPlates(ctx context.Context) (*LHProperties, error) {
	params, err := makeTestGilson(ctx)

	if err != nil {
		return nil, err
	}

	inputPlate, err := makeTestInputPlate(ctx)

	if err != nil {
		return nil, err
	}

	err = params.AddInputPlate(inputPlate)

	if err != nil {
		return nil, err
	}

	outputPlate, err := makeTestOutputPlate(ctx)

	if err != nil {
		return nil, err
	}

	err = params.AddOutputPlate(outputPlate)

	if err != nil {
		return nil, err
	}
	return params, nil
}
func makeTestGilson(ctx context.Context) (*LHProperties, error) {
	params := makeGilson()

	tw, err := inventory.NewTipwaste(ctx, "Gilsontipwaste")
	if err != nil {
		return nil, err
	}
	params.AddTipWaste(tw)

	tb, err := inventory.NewTipbox(ctx, "DL10 Tip Rack (PIPETMAX 8x20)")
	if err != nil {
		return nil, err
	}
	params.AddTipBox(tb)

	tb, err = inventory.NewTipbox(ctx, "DF200 Tip Rack (PIPETMAX 8x200)")
	if err != nil {
		return nil, err
	}
	params.AddTipBox(tb)

	return params, nil
}

func makeTestInputPlate(ctx context.Context) (*wtype.LHPlate, error) {
	p, err := inventory.NewPlate(ctx, "DWST12")

	if err != nil {
		return nil, err
	}

	c, err := inventory.NewComponent(ctx, "water")

	if err != nil {
		return nil, err
	}

	c.Vol = 5000.0 // ul

	p.AddComponent(c, true)

	return p, nil
}

func makeTestOutputPlate(ctx context.Context) (*wtype.LHPlate, error) {
	p, err := inventory.NewPlate(ctx, "DSW96")

	if err != nil {
		return nil, err
	}

	return p, nil
}
