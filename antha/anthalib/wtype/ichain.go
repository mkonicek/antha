package wtype

import (
	"fmt"
	"sort"
	"strings"
)

type IChain struct {
	Parent *IChain
	Child  *IChain
	Values []*LHInstruction
	Depth  int
}

func NewIChain(parent *IChain) *IChain {
	var it IChain
	it.Parent = parent
	it.Values = make([]*LHInstruction, 0, 1)
	if parent != nil {
		it.Depth = parent.Depth + 1
		parent.Child = &it
	}
	return &it
}

// depth from here
func (it *IChain) Height() int {
	if it == nil {
		return 0
	}

	return it.Child.Height() + 1
}

func (it *IChain) PruneOut(Remove map[string]bool) *IChain {
	if it == nil || len(Remove) == 0 || len(it.Values) == 0 {
		return it
	}

	it.Child = it.Child.PruneOut(Remove)

	newValues := make([]*LHInstruction, 0, len(it.Values))

	for _, v := range it.Values {
		if Remove[v.ID] {
			continue
		}
		newValues = append(newValues, v)
		delete(Remove, v.ID)
	}

	// if we've removed a whole layer, get rid of it

	if len(newValues) == 0 {

		if it.Child != nil {
			it.Child.Parent = it.Parent
		}

		if it.Parent != nil {
			it.Parent.Child = it.Child
		}

		return it.Child

	} else {
		it.Values = newValues
		return it
	}

}

func (it *IChain) AsList(ica []*IChain) []*IChain {
	if it == nil {
		return ica
	}

	ica = append(ica, it)

	return it.Child.AsList(ica)
}

func (it *IChain) Reverse() {
	if it.Child != nil {
		it.Child.Reverse()
	}
	// swap parent and child
	p := it.Parent
	it.Parent = it.Child
	it.Child = p
}

func (it *IChain) ValueIDs() []string {
	r := make([]string, 0, 1)

	for _, v := range it.Values {
		r = append(r, v.ID)
	}
	return r
}

func (it *IChain) Print() {
	fmt.Println("****")
	fmt.Println("\tPARENT NIL: ", it.Parent == nil)
	if len(it.Values) > 0 {
		for j := 0; j < len(it.Values); j++ {
			if it.Values[j].Type == LHIMIX {
				fmt.Printf("MIX    %2d: %s \n", j, it.Values[j].ID)
				for i := 0; i < len(it.Values[j].Inputs); i++ {
					fmt.Print(" ", it.Values[j].Inputs[i].ID, ":", it.Values[j].Inputs[i].FullyQualifiedName(), "@", it.Values[j].Inputs[i].Volume().ToString(), " \n")
				}
				fmt.Println(":", it.Values[j].Outputs[0].ID, ":", it.Values[j].Platetype, " ", it.Values[j].PlateName, " ", it.Values[j].Welladdress)
				fmt.Printf("-- ")
			} else if it.Values[j].Type == LHIPRM {
				fmt.Println("PROMPT ", it.Values[j].Message, "-- ")
				for i := range it.Values[j].Inputs {
					fmt.Println(it.Values[j].Inputs[i].ID, ":::", it.Values[j].Outputs[i].ID, " --")
				}
			} else if it.Values[j].Type == LHISPL {
				fmt.Printf("SPLIT %2d: %s ", j, it.Values[j].ID)
				fmt.Println(" ", it.Values[j].Inputs[0].ID, ":", it.Values[j].Inputs[0].FullyQualifiedName(), " : ", it.Values[j].PlateName, " ", it.Values[j].Welladdress)
				fmt.Println(" MOVE:", it.Values[j].Outputs[0].ID, ":", it.Values[j].Outputs[0].FullyQualifiedName(), "@", it.Values[j].Outputs[0].Volume().ToString())
				fmt.Println(" STAY:", it.Values[j].Outputs[1].ID, ":", it.Values[j].Outputs[1].FullyQualifiedName(), "@", it.Values[j].Outputs[1].Volume().ToString())
				fmt.Printf("-- \n")
			} else {
				fmt.Println("WTF?   ", InsType(it.Values[j].Type), "-- ")
			}
		}
		fmt.Println()
	}
	if it.Child != nil {
		it.Child.Print()
	}
}

//FlattenInstructionIDs returns a slice containing the IDs of each instruction
//in the chain in order
func (it *IChain) FlattenInstructionIDs() []string {
	return it.flattenInstructionIDs(nil)
}

func (it *IChain) flattenInstructionIDs(acc []string) []string {
	if it == nil {
		return acc
	} else {
		for _, v := range it.Values {
			acc = append(acc, v.ID)
		}
		return it.Child.flattenInstructionIDs(acc)
	}
}

//GetOrderedLHInstructions get the instructions in order
func (it *IChain) GetOrderedLHInstructions() []*LHInstruction {
	return it.getOrderedLHInstructions(nil)
}

func (it *IChain) getOrderedLHInstructions(acc []*LHInstruction) []*LHInstruction {
	if it == nil {
		return acc
	} else {
		acc = append(acc, it.Values...)
		return it.Child.getOrderedLHInstructions(acc)
	}
}

func (it *IChain) SplitMixedNodes() *IChain {
	if it.Child != nil {
		it.Child.SplitMixedNodes()
	}

	if len(it.getInstructionTypes()) > 1 {
		return it.splitMixedNode()
	}
	return it
}

func (it *IChain) FindEnd() *IChain {
	if it.Child == nil {
		return it
	}

	return it.Child.FindEnd()
}

func (it *IChain) splitMixedNode() *IChain {
	// put mixes first, then splits, then prompts
	mixValues := make([]*LHInstruction, 0, len(it.Values))
	splitValues := make([]*LHInstruction, 0, len(it.Values))
	promptValues := make([]*LHInstruction, 0, len(it.Values))

	for _, v := range it.Values {
		if v.Type == LHIMIX {
			mixValues = append(mixValues, v)
		} else if v.Type == LHISPL {
			splitValues = append(splitValues, v)
		} else if v.Type == LHIPRM {
			promptValues = append(promptValues, v)
		} else {
			panic("Wrong instruction type passed through to instruction chain split")
		}
	}

	// make new chain

	newch := makeNewIChain(mixValues, splitValues, promptValues)

	// swap it in

	r := it.SwapForChain(newch)

	return r
}

// return a chain containing one node for each argument, linked in sequence
// skip any empty sets
func makeNewIChain(vals ...[]*LHInstruction) *IChain {
	var top, cur *IChain

	for _, v := range vals {
		if len(v) == 0 {
			continue
		}

		cur = NewIChain(cur)

		if top == nil {
			top = cur
		}

		cur.Values = v
	}

	return top
}

//SwapForChain replace node ch with the node chain starting with newch
//             if ch is the head of the chain, return newch as the new head
//             otherwise return nil
func (ch *IChain) SwapForChain(newch *IChain) *IChain {
	if ch.Child != nil {
		end := newch.FindEnd()
		ch.Child.Parent = end
		end.Child = ch.Child
	}

	if ch.Parent != nil {
		ch.Parent.Child = newch
		newch.Parent = ch.Parent
	} else {
		return newch
	}

	return nil
}

func (self *IChain) getInstructionTypes() map[string]bool {
	types := make(map[string]bool, len(self.Values))
	for _, v := range self.Values {
		types[v.InsType()] = true
	}

	return types
}

//AssertInstructionsSeparate check that there's only one type of instruction
//in each link of the chain
func (self *IChain) AssertInstructionsSeparate() error {
	if self == nil {
		return nil
	}

	types := self.getInstructionTypes()

	if len(types) != 1 {
		return fmt.Errorf("Only one instruction type per stage is allowed, found %v at stage %d, %v", len(types), self.Depth, types)
	}

	return self.Child.AssertInstructionsSeparate()
}

type ByColumn []*LHInstruction

func (bg ByColumn) Len() int      { return len(bg) }
func (bg ByColumn) Swap(i, j int) { bg[i], bg[j] = bg[j], bg[i] }
func (bg ByColumn) Less(i, j int) bool {
	// compare any messages present (only really applies to prompts)
	c := strings.Compare(bg[i].Message, bg[j].Message)

	if c != 0 {
		return c < 0
	}
	// compare the plate names (which must exist now)
	//	 -- oops, I think this has ben violated by moving the sort
	// 	 TODO check and fix

	c = strings.Compare(bg[i].PlateName, bg[j].PlateName)

	if c != 0 {
		return c < 0
	}

	// Go Down Columns

	return CompareStringWellCoordsCol(bg[i].Welladdress, bg[j].Welladdress) < 0
}

// Optimally - order by component.
type ByResultComponent []*LHInstruction

func (bg ByResultComponent) Len() int      { return len(bg) }
func (bg ByResultComponent) Swap(i, j int) { bg[i], bg[j] = bg[j], bg[i] }
func (bg ByResultComponent) Less(i, j int) bool {
	// compare any messages present

	c := strings.Compare(bg[i].Message, bg[j].Message)

	if c != 0 {
		return c < 0
	}

	// compare the names of the resultant components
	c = strings.Compare(bg[i].Outputs[0].CName, bg[j].Outputs[0].CName)

	if c != 0 {
		return c < 0
	}

	// if two components names are equal, then compare the plates
	c = strings.Compare(bg[i].PlateName, bg[j].PlateName)

	if c != 0 {
		return c < 0
	}

	// finally go down columns (nb need to add option)

	return CompareStringWellCoordsCol(bg[i].Welladdress, bg[j].Welladdress) < 0
}

//SortInstructions sort the instructions within each link of the chain
func (ic *IChain) SortInstructions(byComponent bool) {
	if ic == nil {
		return
	}

	if byComponent {
		sort.Sort(ByResultComponent(ic.Values))
	} else {
		sort.Sort(ByColumn(ic.Values))
	}

	ic.Child.SortInstructions(byComponent)
}
