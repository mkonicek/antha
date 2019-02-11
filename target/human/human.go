package human

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/handler"
)

var (
	_ effects.Device = &Human{}
)

// A Human is a device that can do anything
type Human struct {
	impl        *handler.GenericHandler
	canMix      bool
	canIncubate bool
}

// New returns a new human device
func New() *Human {
	h := &Human{}
	h.impl = &handler.GenericHandler{
		GenFunc: h.generate,
	}

	return h
}

func (a *Human) CanCompile(req effects.Request) bool {
	can := effects.Request{
		Selector: []effects.NameValue{
			target.DriverSelectorV1Human,
		},
	}

	if a.canMix {
		can.Selector = append(can.Selector, target.DriverSelectorV1Mixer)
	}
	if a.canIncubate {
		can.Selector = append(can.Selector, target.DriverSelectorV1ShakerIncubator)
	}

	return can.Contains(req)
}

// Compile implements target.device Compile
func (a *Human) Compile(labEffects *effects.LaboratoryEffects, nodes []effects.Node) ([]effects.Inst, error) {
	return a.impl.Compile(labEffects, nodes)
}

func (a *Human) DetermineRole(tgt *target.Target) {
	mixReq := effects.Request{
		Selector: []effects.NameValue{
			target.DriverSelectorV1Mixer,
		},
	}

	incubateReq := effects.Request{
		Selector: []effects.NameValue{
			target.DriverSelectorV1ShakerIncubator,
		},
	}

	mix := true
	incubate := true
	for _, dev := range tgt.Devices {
		if mix && dev.CanCompile(mixReq) {
			mix = false
		}
		if incubate && dev.CanCompile(incubateReq) {
			incubate = false
		}
		if !(mix || incubate) {
			break
		}
	}
	a.canMix = mix
	a.canIncubate = incubate

	if a.canMix || a.canIncubate {
		tgt.AddDevice(a)
	}
}

func (a *Human) generate(cmd interface{}) ([]effects.Inst, error) {
	instrs := make([]effects.Inst, 1)

	switch cmd := cmd.(type) {

	case *wtype.LHInstruction:
		instrs[0] = &target.Manual{
			Dev:     a,
			Label:   "mix",
			Details: prettyMixDetails(cmd),
		}

	case *effects.IncubateInst:
		instrs[0] = &target.Manual{
			Dev:     a,
			Label:   "incubate",
			Details: fmt.Sprintf("incubate at %s for %s", cmd.Temp.ToString(), cmd.Time.ToString()),
		}

	case *effects.HandleInst:
		instrs[0] = &target.Manual{
			Dev:   a,
			Label: cmd.Group,
		}

	case *effects.PromptInst:
		instrs[0] = &target.Prompt{
			Message: cmd.Message,
		}

	case *wtype.PRInstruction:
		instrs[0] = &target.Manual{
			Dev:     a,
			Label:   "plate-read",
			Details: fmt.Sprintf("plate-read instruction. Options:'%s'", cmd.Options),
		}

	case *effects.QPCRInstruction:
		instrs[0] = &target.Manual{
			Dev:     a,
			Label:   "QPCR",
			Details: fmt.Sprintf("QPCR request, definition %s, barcode %s", cmd.Definition, cmd.Barcode),
		}

	default:
		return nil, fmt.Errorf("unknown inst %T", cmd)
	}

	return instrs, nil
}

func prettyMixDetails(inst *wtype.LHInstruction) string {
	if len(inst.PlateName) != 0 || len(inst.Welladdress) != 0 {
		return fmt.Sprintf("mix %q[%q]", inst.PlateName, inst.Welladdress)
	}
	return "mix"
}
