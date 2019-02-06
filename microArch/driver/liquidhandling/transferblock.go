package liquidhandling

import (
	"context"
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
)

type TransferBlockInstruction struct {
	BaseRobotInstruction
	*InstructionType
	Inss []*wtype.LHInstruction
}

func NewTransferBlockInstruction(inss []*wtype.LHInstruction) *TransferBlockInstruction {
	tb := &TransferBlockInstruction{
		InstructionType: TFB,
		Inss:            inss,
	}
	tb.BaseRobotInstruction = NewBaseRobotInstruction(tb)
	return tb
}

func (ins *TransferBlockInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.TransferBlock(ins)
}

// this attempts to find arrays of destinations which can potentially be done simultaneously
// via multichannel operation. At present this means they must be aligned in rows or columns
// depending on the robot type and configuration

func (ti TransferBlockInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, robot *LHProperties) ([]RobotInstruction, error) {
	var tfr []*TransferInstruction
	var err error

	// assessing evaporation with this potentially
	//timer := robot.GetTimer()
	inss := make([]RobotInstruction, 0, 1)
	insm := make(map[string]*wtype.LHInstruction, len(ti.Inss))
	seen := make(map[string]bool)

	for _, ins := range ti.Inss {
		insm[ins.ID] = ins
	}

	// list of ids
	parallel_sets, prm, err := get_parallel_sets_robot(ctx, ti.Inss, robot, policy)

	// what if prm is nil?

	if err != nil {
		return inss, err
	}

	for _, set := range parallel_sets {
		// compile the instructions and pass them through
		insset := make([]*wtype.LHInstruction, len(set))

		for i, id := range set {
			// parallel sets are arranged in accordance with destination layout
			// hence can include gaps
			if id == "" {
				continue
			}
			seen[id] = true
			insset[i] = insm[id]
		}

		// aggregates across components
		//TODO --> allow setting legacy volume if necessary

		tfr, err = ConvertInstructions(ctx, insset, robot, prm, prm.Multi, false, policy)
		if err != nil {
			return inss, err
		}

		// we merge instructions which are compatible
		//tfr = mergeTransfers(tfr)

		for _, tf := range tfr {
			inss = append(inss, RobotInstruction(tf))
		}
	}

	toTransfers := func(inss []RobotInstruction) []*TransferInstruction {
		r := make([]*TransferInstruction, len(inss))

		for ix, ins := range inss {
			r[ix] = ins.(*TransferInstruction)
		}

		return r
	}

	fromTransfers := func(inss []*TransferInstruction) []RobotInstruction {
		r := make([]RobotInstruction, len(inss))

		for i, ins := range inss {
			r[i] = RobotInstruction(ins)
		}

		return r
	}

	// stuff that can't be done in parallel
	for _, ins := range ti.Inss {
		if seen[ins.ID] {
			continue
		}

		//insset = append(insset, ins)

		// now make transfer and append
		// prm here will be nil unless len(insset)==0
		// we must either tolerate this or do something else

		insset := []*wtype.LHInstruction{ins}

		tfr, err = ConvertInstructions(ctx, insset, robot, prm, 1, false, policy)

		if err != nil {
			return inss, err
		}

		for _, tf := range tfr {
			inss = append(inss, RobotInstruction(tf))
		}
	}

	//inss = append(inss, tfr...)
	inss = fromTransfers(mergeTransfers(toTransfers(inss), policy))

	return inss, nil
}

type IDSet []string
type SetOfIDSets []IDSet

func get_parallel_sets_robot(ctx context.Context, ins []*wtype.LHInstruction, robot *LHProperties, policy *wtype.LHPolicyRuleSet) (SetOfIDSets, *wtype.LHChannelParameter, error) {
	//  depending on the configuration and options we may have to try and
	//  use one or both of H / V or... whatever
	//  -- issue is this choice and choosechannel conflict with one another
	//  since we may only be able to do certain volumes with certain heads
	//  ... should account for that here, at least avoid passing things
	// that cannot work

	// part of the model here is just to make things possible, so that later
	// on we can at least make this choice

	headsLoaded := robot.GetLoadedHeads()
	possible_sets := make([]SetOfIDSets, 0, len(headsLoaded))
	corresponding_params := make([]*wtype.LHChannelParameter, 0, 1)

	for _, head := range headsLoaded {
		// ignore heads which do not have multi

		if head.GetParams().Multi == 1 {
			continue
		}

		// also TODO here -- allow adaptor changes
		sids, err := get_parallel_sets_head(ctx, head, ins)

		if err != nil {
			return SetOfIDSets{}, &wtype.LHChannelParameter{}, err
		}

		possible_sets = append(possible_sets, sids)
		corresponding_params = append(corresponding_params, head.GetParams())
	}

	// don't even bother if there are no possible sets

	if len(possible_sets) == 0 {
		return SetOfIDSets{}, &wtype.LHChannelParameter{}, nil
	}

	// now we make our choice
	return choose_parallel_sets(possible_sets, corresponding_params, ins)
}

type InsByComponent []*wtype.LHInstruction

func (ibc InsByComponent) Len() int      { return len(ibc) }
func (ibc InsByComponent) Swap(i, j int) { ibc[i], ibc[j] = ibc[j], ibc[i] }
func (ibc InsByComponent) Less(i, j int) bool {
	return strings.Compare(ibc[i].Outputs[0].CName, ibc[j].Outputs[0].CName) < 0
}

type InsByRow []*wtype.LHInstruction

func (ibc InsByRow) Len() int      { return len(ibc) }
func (ibc InsByRow) Swap(i, j int) { ibc[i], ibc[j] = ibc[j], ibc[i] }
func (ibc InsByRow) Less(i, j int) bool {
	wca := wtype.MakeWellCoords(ibc[i].Welladdress)
	wcb := wtype.MakeWellCoords(ibc[j].Welladdress)
	return wtype.CompareWellCoordsRow(wca, wcb) < 0
}

type InsByCol []*wtype.LHInstruction

func (ibc InsByCol) Len() int      { return len(ibc) }
func (ibc InsByCol) Swap(i, j int) { ibc[i], ibc[j] = ibc[j], ibc[i] }
func (ibc InsByCol) Less(i, j int) bool {
	wca := wtype.MakeWellCoords(ibc[i].Welladdress)
	wcb := wtype.MakeWellCoords(ibc[j].Welladdress)
	return wtype.CompareWellCoordsCol(wca, wcb) < 0
}

// limited to SBS format plates for now
func get_parallel_sets_head(ctx context.Context, head *wtype.LHHead, ins []*wtype.LHInstruction) (SetOfIDSets, error) {
	// surely not

	if len(ins) == 0 {
		return nil, fmt.Errorf("No instructions")
	}

	// sort instructions to keep components together

	//sort.Sort(InsByComponent(ins))

	ret := make(SetOfIDSets, 0, 1)

	// h maps plate IDs to platedestmaps
	// platedestmaps are 2d arrays of instructions arranged
	// to mirror the layout of a plate (in fact limited to a 96x96 grid, but
	// that's pretty big by comparison to any existing plate)

	h := make(map[string]wtype.Platedestmap, 2)

	platedims := make(map[string]wtype.Rational)

	prm := head.GetParams()

	for _, i := range ins {
		// ignore empty instructions
		if len(i.Inputs) == 0 {
			continue
		}

		wc := wtype.MakeWellCoords(i.Welladdress)

		_, ok := h[i.PlateID]

		if !ok {
			h[i.PlateID] = wtype.NewPlatedestmap()

			// gerrabirrovinfo on the plate type
			// is this always set??
			pt, err := inventory.NewPlate(ctx, i.Platetype)

			if err != nil {
				return ret, wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("No plate type exists for ID %s - requested was %s", i.PlateID, i.Platetype))
			}

			platedims[i.PlateID] = wtype.Rational{N: pt.WellsX(), D: pt.WellsY()}
		}

		h[i.PlateID][wc.X][wc.Y] = append(h[i.PlateID][wc.X][wc.Y], i)
	}

	if len(h) == 0 {
		return nil, fmt.Errorf("No plate destinations")
	}

	for id, pdm := range h {
		dims := platedims[id]
		switch prm.Orientation {
		case wtype.LHHChannel:
			r := get_rows(pdm, prm.Multi, dims.N, !prm.Independent, false)
			if len(ret) == 0 {
				ret = r
			} else {
				ret = append(ret, r...)
			}
		case wtype.LHVChannel:
			r := get_cols(pdm, prm.Multi, dims.D, !prm.Independent, false)

			if len(ret) == 0 {
				ret = r
			} else {
				ret = append(ret, r...)
			}
		}
	}

	return ret, nil
}

func get_rows(pdm wtype.Platedestmap, multi, wells int, contiguous, full bool) SetOfIDSets {
	ret := make(SetOfIDSets, 0, 1)
	row := 0
	for {
		if row >= len(pdm) {
			break
		}

		rowset := get_row(pdm, row, multi, wells, contiguous, full)

		if len(rowset) != 0 {
			ret = append(ret, rowset)
		} else {
			row += 1
		}
	}

	return ret
}

func get_row(pdm wtype.Platedestmap, row, multi, wells int, contiguous, full bool) IDSet {
	var ret IDSet
	for s := 0; s < len(pdm[row])-2; s++ {
		ret = make(IDSet, 0, multi)
		newrow := make([][]*wtype.LHInstruction, len(pdm[row]))

		tipsperwell := 1

		if wells < multi {
			// if this isn't an even multiple it should be rejected

			if multi%wells != 0 {
				///urrr
				return ret
			}

			tipsperwell = multi / wells

		}

		for c := s; c < len(pdm[row]); c++ {
			if len(pdm[row][c]) >= tipsperwell {
				for x := 0; x < tipsperwell; x++ {
					id := pdm[row][c][x].ID
					ret = append(ret, id)
				}
				newrow[c] = pdm[row][c][tipsperwell:]
			} else if contiguous {
				break
			}
		}

		if len(ret) != multi && full {
			return make(IDSet, 0, 1)
		} else if len(ret) == 0 {
			continue
		} else {
			pdm[row] = newrow
			return ret
		}
	}

	return ret
}

func get_cols(pdm wtype.Platedestmap, multi, wells int, contiguous, full bool) SetOfIDSets {
	ret := make(SetOfIDSets, 0, 1)
	col := 0

	countUsed := func(sa []string) int {
		c := 0
		for _, v := range sa {
			if v != "" {
				c += 1
			}
		}
		return c
	}

	for {
		if col >= len(pdm) {
			break
		}

		wellsPerTip := 1

		if wells > multi {
			wellsPerTip = wells / multi
		}
		for modulo := 0; modulo < wellsPerTip; modulo++ {
			colset := get_col(pdm, col, multi, wells, contiguous, full, modulo)
			if countUsed(colset) != 0 {
				ret = append(ret, colset)
				//} else {
				//	break
			}
		}

		col += 1

	}

	return ret
}

func colDup(in [][]*wtype.LHInstruction) [][]*wtype.LHInstruction {
	out := make([][]*wtype.LHInstruction, len(in))

	for i, v := range in {
		r := make([]*wtype.LHInstruction, len(v))

		for i := 0; i < len(v); i++ {
			r[i] = v[i]
		}
		out[i] = r
	}

	return out
}

func get_col(pdm wtype.Platedestmap, col, multi, wells int, contiguous, full bool, modulo int) IDSet {
	var ret IDSet
	tipsperwell := 1

	if wells < multi {
		// if this isn't an even multiple it should be rejected

		if multi%wells != 0 {
			return ret
		}

		tipsperwell = multi / wells
	}

	wellsPerTip := 1

	if wells > multi {
		// as above if not an even multiple it should be rejected

		if wells%multi != 0 {
			return ret
		}

		wellsPerTip = wells / multi
	}

	for s := 0; s < len(pdm[col])-2; s++ {
		ret = make(IDSet, multi)
		//newcol := make([][]*wtype.LHInstruction, len(pdm[col]))
		newcol := colDup(pdm[col])
		used := 0 // number of instructions returned
		offset := 0
		for c := s + modulo; c < len(pdm[col]); c += wellsPerTip {
			if len(pdm[col][c]) >= tipsperwell {
				for x := 0; x < tipsperwell; x++ {
					id := pdm[col][c][x].ID
					//ret = append(ret, id)
					ret[offset] = id
					offset += 1
					used += 1
				}
				newcol[c] = pdm[col][c][tipsperwell:]
			} else if contiguous {
				break
			} else {
				offset += 1
			}

			if used == multi {
				break
			}
		}

		if used != multi && full {
			return make(IDSet, 0, 1)
		} else if used == 0 {
			continue
		} else {
			pdm[col] = newcol

			return ret
		}
	}

	return ret
}

func choose_parallel_sets(sets []SetOfIDSets, params []*wtype.LHChannelParameter, ins []*wtype.LHInstruction) (SetOfIDSets, *wtype.LHChannelParameter, error) {
	var ret SetOfIDSets
	var retp *wtype.LHChannelParameter
	// just one or the other to start with

	mx := 0
	for i, s := range sets {
		if len(s) > mx {
			mx = len(s)
			ret = s
			retp = params[i]
		}
	}

	return ret, retp, nil
}

func (ti TransferBlockInstruction) GetParameter(name InstructionParameter) interface{} {
	return ti.BaseRobotInstruction.GetParameter(name)
}

func mergeTransfers(tfrs []*TransferInstruction, policy *wtype.LHPolicyRuleSet) []*TransferInstruction {
	ret := make([]*TransferInstruction, 0, len(tfrs))

	for _, tf := range tfrs {
		forMerge := findTransferForMerge(tf, ret, policy)

		// true if ret is empty or nothing mergeable within
		if forMerge == nil {
			ret = append(ret, tf)
		} else {
			// forMerge is already in ret
			forMerge.MergeWith(tf)
		}
	}

	return ret
}

func findTransferForMerge(ins *TransferInstruction, arr []*TransferInstruction, policy *wtype.LHPolicyRuleSet) *TransferInstruction {
	for _, ins2 := range arr {
		if canMerge(ins, ins2, policy) {
			return ins2
		}
	}

	return nil
}

func canMulti(policy wtype.LHPolicy) bool {
	return policy["CAN_MULTI"].(bool)
}

func canMerge(ins, ins2 *TransferInstruction, policy *wtype.LHPolicyRuleSet) bool {
	// merge only if the merge doesn't break either

	ins3 := ins.Dup()
	ins3.MergeWith(ins2)

	pol1, _ := GetPolicyFor(policy, ins)
	pol2, _ := GetPolicyFor(policy, ins2)
	pol3, _ := GetPolicyFor(policy, ins3)

	if canMulti(pol1) == canMulti(pol2) {
		return canMulti(pol1) == canMulti(pol3)
	}

	return false
}
