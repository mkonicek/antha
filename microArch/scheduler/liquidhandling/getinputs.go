package liquidhandling

import (
	"fmt"
	"sort"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

//InputSolutions properties to do with the input Liquids required for the mix
type InputSolutions struct {
	Solutions       map[string][]*wtype.Liquid //the solutions explicitly supplied to the protocol
	Order           []string                   //the order in which liquids are required in the ordered LH instructions
	VolumesSupplied map[string]wunit.Volume    //the volumes of the solutions explicitly supplied to the protocol
	VolumesRequired map[string]wunit.Volume    //the estimated volumes of solutions required to carry out the protocol
	VolumesWanting  map[string]wunit.Volume    //the estimated shortfall between the supplied and required volumes which must be auto-allocated
}

// String return a string representation, useful for debugging
func (self *InputSolutions) String() string {
	allNamesMap := make(map[string]bool)
	for name := range self.VolumesSupplied {
		allNamesMap[name] = true
	}
	for name := range self.VolumesRequired {
		allNamesMap[name] = true
	}
	for name := range self.VolumesWanting {
		allNamesMap[name] = true
	}
	allNames := make([]string, 0, len(allNamesMap))
	for name := range allNamesMap {
		allNames = append(allNames, name)
	}
	sort.Strings(allNames)

	ret := []string{"name, supplied, required, wanting"}
	for _, name := range allNames {
		ret = append(ret, fmt.Sprintf("%s, %v, %v, %v", name, self.VolumesSupplied[name], self.VolumesRequired[name], self.VolumesWanting[name]))
	}

	return strings.Join(ret, "\n")
}

// getInputs calculate the volumes required for each input solution by looking
// through the high level liquidhandling (LH) instructions.
// inputSolutions gives the solutions explicitly provided to the protocol, and
// is used to calculate the shortfall, i.e how much of each must be auto-allocated
func (rq *LHRequest) getInputs(carryVolume wunit.Volume) (*InputSolutions, error) {

	orderedInstructions, err := rq.GetOrderedLHInstructions()
	if err != nil {
		return nil, err
	}

	//find what liquids are explicitly provided by the user
	inputSolutions, err := rq.GetSolutionsFromInputPlates()
	if err != nil {
		return nil, err
	}

	inputs := make(map[string][]*wtype.Liquid)
	volsRequired := make(map[string]wunit.Volume)
	var allinputs []string
	ordH := make(map[string]int, len(orderedInstructions))

	for _, instruction := range orderedInstructions {
		// ignore non-mixes
		if instruction.InsType() != "MIX" {
			continue
		}

		for ix, component := range instruction.Inputs {
			// Ignore components which already exist
			if component.IsInstance() {
				continue
			}

			// what if this is a mix in place?
			if ix == 0 && !component.IsSample() {
				// these components come in as instances -- hence 1 per well
				// but if not allocated we need to do so
				inputs[component.CNID()] = []*wtype.Liquid{component}
				allinputs = append(allinputs, component.CNID())
				volsRequired[component.CNID()] = component.Volume()
				component.DeclareInstance()

				if _, ok := ordH[component.CNID()]; !ok {
					ordH[component.CNID()] = len(ordH)
				}
			} else {
				cmps, ok := inputs[component.Kind()]
				if !ok {
					cmps = make([]*wtype.Liquid, 0, 3)
					allinputs = append(allinputs, component.Kind())
				}

				cmps = append(cmps, component)
				inputs[component.Kind()] = cmps

				if _, ok = ordH[component.Kind()]; !ok {
					ordH[component.Kind()] = len(ordH)
				}

				// similarly add the volumes up

				vol := volsRequired[component.Kind()]

				if vol.IsNil() {
					vol = wunit.NewVolume(0.0, "ul")
				}

				v2a := wunit.NewVolume(component.Vol, component.Vunit)

				// we have to add the carry volume here
				// this is roughly per transfer so should be OK
				v2a.Add(carryVolume)
				vol.Add(v2a)

				volsRequired[component.Kind()] = vol
			}
		}
	}

	// work out how much we have and how much we need
	// need to consider what to do with IDs

	// invert the Hash

	inputOrder, err := OrdinalFromHash(ordH)
	if err != nil {
		return nil, err
	}

	volsSupplied := make(map[string]wunit.Volume, len(volsRequired))
	volsWanting := make(map[string]wunit.Volume, len(volsRequired))

	for _, k := range allinputs {
		// volSupplied: how much comes in
		volSupplied := wunit.NewVolume(0.0, "ul")
		for _, cmp := range inputSolutions[k] {
			volSupplied.Add(cmp.Volume())
		}
		volsSupplied[k] = volSupplied

		// volWanted: how much extra we wanted
		if volWanted := wunit.SubtractVolumes(volsRequired[k], volSupplied); volWanted.IsPositive() {
			volsWanting[k] = volWanted
		}
	}

	// add any new inputs
	for k, v := range inputs {
		if inputSolutions[k] == nil {
			inputSolutions[k] = v
		}
	}

	return &InputSolutions{
		Solutions:       inputSolutions,
		Order:           inputOrder,
		VolumesSupplied: volsSupplied,
		VolumesRequired: volsRequired,
		VolumesWanting:  volsWanting,
	}, nil
}

func OrdinalFromHash(m map[string]int) ([]string, error) {
	s := make([]string, len(m))

	// no collisions allowed!

	for k, v := range m {
		if s[v] != "" {
			return nil, fmt.Errorf("Error: ordinal %d appears twice!", v)
		}

		s[v] = k
	}

	return s, nil
}
