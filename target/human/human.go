package human

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/laboratory/effects/id"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/handler"
	"github.com/antha-lang/antha/workflow"
)

var (
	_ effects.Device = &Human{}
)

// A Human is a device that can do anything
type Human struct {
	id workflow.DeviceInstanceID

	impl        *handler.GenericHandler
	canMix      bool
	canIncubate bool
}

// New returns a new human device
func New(idGen *id.IDGenerator) *Human {
	h := &Human{
		id: workflow.DeviceInstanceID(idGen.NextID()) + "_human",
	}
	h.impl = &handler.GenericHandler{
		GenFunc: h.generate,
	}

	return h
}

func (hum *Human) Id() workflow.DeviceInstanceID {
	return hum.id
}

func (hum *Human) CanCompile(req effects.Request) bool {
	can := effects.Request{
		Selector: []effects.NameValue{
			target.DriverSelectorV1Human,
		},
	}

	if hum.canMix {
		can.Selector = append(can.Selector, target.DriverSelectorV1Mixer)
	}
	if hum.canIncubate {
		can.Selector = append(can.Selector, target.DriverSelectorV1ShakerIncubator)
	}

	return can.Contains(req)
}

// Compile implements target.device Compile
func (hum *Human) Compile(labEffects *effects.LaboratoryEffects, dir string, nodes []effects.Node) ([]effects.Inst, error) {
	return hum.impl.Compile(labEffects, dir, nodes)
}

func (hum *Human) DetermineRole(tgt *target.Target) {
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

	hum.canMix = true
	hum.canIncubate = true
	for _, dev := range tgt.Devices {
		if hum.canMix && dev.CanCompile(mixReq) {
			hum.canMix = false
		}
		if hum.canIncubate && dev.CanCompile(incubateReq) {
			hum.canIncubate = false
		}
		if !hum.canMix && !hum.canIncubate {
			break
		}
	}

	if hum.canMix || hum.canIncubate {
		tgt.AddDevice(hum)
	}
}

func (hum *Human) Connect(*workflow.Workflow) error {
	return nil
}

func (hum *Human) Close() {}

func (hum *Human) generate(cmd interface{}) ([]effects.Inst, error) {
	instrs := make([]effects.Inst, 1)

	switch cmd := cmd.(type) {

	case *wtype.LHInstruction:
		instrs[0] = &target.Manual{
			Dev:     hum,
			Label:   "mix",
			Details: prettyMixDetails(cmd),
		}

	case *effects.IncubateInst:
		instrs[0] = &target.Manual{
			Dev:     hum,
			Label:   "incubate",
			Details: fmt.Sprintf("incubate at %s for %s", cmd.Temp.ToString(), cmd.Time.ToString()),
		}

	case *effects.PromptInst:
		instrs[0] = &target.Prompt{
			Message: cmd.Message,
		}

	case *wtype.PRInstruction:
		instrs[0] = &target.Manual{
			Dev:     hum,
			Label:   "plate-read",
			Details: fmt.Sprintf("plate-read instruction. Options:'%s'", cmd.Options),
		}

	case *effects.QPCRInstruction:
		instrs[0] = &target.Manual{
			Dev:     hum,
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
