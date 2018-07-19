package liquidhandling

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/logger"
)

//InputSolutions properties to do with the input Liquids required for the mix
type InputSolutions struct {
	Solutions       map[string][]*wtype.Liquid
	Order           []string
	VolumesSupplied map[string]wunit.Volume
	VolumesRequired map[string]wunit.Volume
	VolumesWanting  map[string]wunit.Volume
}

// sort out inputs
func GetInputs(orderedInstructions []*wtype.LHInstruction, inputSolutions map[string][]*wtype.Liquid, carryVolume wunit.Volume) (*InputSolutions, error) {
	inputs := make(map[string][]*wtype.Liquid)
	volsRequired := make(map[string]wunit.Volume)
	var allinputs []string
	ordH := make(map[string]int, len(orderedInstructions))

	for _, instruction := range orderedInstructions {
		// ignore non-mixes
		if instruction.InsType() != "MIX" {
			continue
		}

		for ix, component := range instruction.Components {
			// Ignore components which already exist
			if component.IsInstance() {
				continue
			}

			// what if this is a mix in place?
			if ix == 0 && !component.IsSample() {
				// these components come in as instances -- hence 1 per well
				// but if not allocated we need to do so
				inputs[component.CNID()] = make([]*wtype.Liquid, 0, 1)
				inputs[component.CNID()] = append(inputs[component.CNID()], component)
				allinputs = append(allinputs, component.CNID())
				volsRequired[component.CNID()] = component.Volume()
				component.DeclareInstance()

				// if this already exists do nothing
				_, ok := ordH[component.CNID()]

				if !ok {
					ordH[component.CNID()] = len(ordH)
				}
			} else {
				cmps, ok := inputs[component.Kind()]
				if !ok {
					cmps = make([]*wtype.Liquid, 0, 3)
					allinputs = append(allinputs, component.Kind())
				}

				_, ok = ordH[component.Kind()]

				if !ok {
					ordH[component.Kind()] = len(ordH)
				}

				cmps = append(cmps, component)
				inputs[component.Kind()] = cmps

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

	if len(inputSolutions) == 0 {
		inputSolutions = make(map[string][]*wtype.Liquid, 5)
	}

	volsSupplied := make(map[string]wunit.Volume, len(volsRequired))
	volsWanting := make(map[string]wunit.Volume, len(volsRequired))

	for _, k := range allinputs {
		// volSupplied: how much comes in
		volSupplied := wunit.NewVolume(0.00, "ul")
		for _, cmp := range inputSolutions[k] {
			volSupplied.Add(cmp.Volume())
		}
		volsSupplied[k] = volSupplied

		// volWanted: how much extra we wanted
		volWanted := wunit.SubtractVolumes(volsRequired[k], volSupplied)

		if volWanted.GreaterThanFloat(0.0001) {
			volsWanting[k] = volWanted
		}
		// toggle HERE for DEBUG
		if false {
			logger.Debug(fmt.Sprintf("COMPONENT %s HAVE : %v WANT: %v DIFF: %v", k, volSupplied.ToString(), volsRequired[k].ToString(), volWanted.ToString()))
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
