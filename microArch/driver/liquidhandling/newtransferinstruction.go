package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type TransferBlockInstruction struct {
	GenericRobotInstruction
	Ins []*wtype.LHInstruction
}

func NewTransferBlockInstruction(inss []*wtype.LHInstruction) TransferBlockInstruction {
	tb := TransferBlockInstruction{}
	tb.Ins = inss
	tb.GenericRobotInstruction.Ins = RobotInstruction(&tb)
	return tb
}

func (ti TransferBlockInstruction) InstructionType() int {
	return TFB
}

func (ti TransferBlockInstruction) Generate(policy *wtype.LHPolicyRuleSet, robot *LHProperties) ([]RobotInstruction, error) {
	// timer for assessing evaporation
	// need to define how to make this optional
	inss := make([]RobotInstruction, 0, 1)
	timer := prms.GetTimer()

	seen := make(map[string]bool)

	// list of ids
	parallel_sets := get_parallel_sets_robot(ti.Ins, robot, policy)

	for _, set := range parallel_sets {
		for _, id := range set {
			seen[id] = true
		}

	}

	// stuff that can't be done in parallel

	for _, ins := range ti.Inss {
		if seen[ins.ID] {
			continue
		}
	}

	return inss, nil
}

type IDset []string
type SetOfIDSets []IDset

func get_parallel_sets_robot(ins []*wtype.LHInstruction, robot *LHProperties, policy *LHPolicyRuleSet) SetOfIDSets {
	//  depending on the configuration and options we may have to try and
	//  use one or both of H / V or... whatever
	//  -- issue is this choice and choosechannel conflict with one another
	//  since we may only be able to do certain volumes with certain heads
	//  ... should account for that here, at least avoid passing things
	// that cannot work

	// part of the model here is just to make things possible, so that later
	// on we can at least make this choice

	possible_sets := make([]SetOfIDSets, 0, len(robot.HeadsLoaded))

	for _, head := range robot.HeadsLoaded {
		// ignore heads which do not have multi

		if head.Multi == 1 {
			continue
		}

		// also TODO here -- allow adaptor changes
		sids := get_parallel_sets_head(head, ins)
		possible_sets = append(possible_sets, sids)
	}

	// now we make our choice

	return choose_parallel_sets(possible_sets, ins)
}

func get_parallel_sets_head(head wtype.LHHead, ins []*wtype.LHInstruction) SetOfIDSets {
	ret := make(SetOfIDSets, 0, 1)

	return ret
}

func choose_parallel_sets(sets []SetOfIDSets, ins []*wtype.LHInstruction) SetOfIDSets {
	ret := make(SetOfIDSets, 0, 1)

	return ret
}

func (ti TransferBlockInstruction) GetParameter(p string) interface{} {
	return nil
}
