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
	opt  Opt
	impl *handler.GenericHandler
}

// An Opt is a set of options to configure a human device
type Opt struct {
	CanMix      bool
	CanIncubate bool

	// CanHandle is deprecated
	CanHandle bool
}

// New returns a new human device
func New(opt Opt) *Human {
	h := &Human{opt: opt}
	h.impl = &handler.GenericHandler{
		GenFunc: h.generate,
	}

	return h
}

// CanCompile implements device CanCompile
func (a *Human) CanCompile(req effects.Request) bool {
	can := effects.Request{
		Selector: []effects.NameValue{
			target.DriverSelectorV1Human,
		},
	}

	if a.opt.CanIncubate {
		can.Selector = append(can.Selector, target.DriverSelectorV1ShakerIncubator)
	}

	if a.opt.CanMix {
		can.Selector = append(can.Selector, target.DriverSelectorV1Mixer)
	}

	if a.opt.CanHandle {
		can.Selector = append(can.Selector, req.Selector...)
	}

	return can.Contains(req)
}

// Compile implements target.device Compile
func (a *Human) Compile(labEffects *effects.LaboratoryEffects, nodes []effects.Node) ([]effects.Inst, error) {
	return a.impl.Compile(labEffects, nodes)
}

func (a *Human) generate(cmd interface{}) ([]effects.Inst, error) {

	var insts []effects.Inst

	switch cmd := cmd.(type) {

	case *wtype.LHInstruction:
		insts = append(insts, &target.Manual{
			Dev:     a,
			Label:   "mix",
			Details: prettyMixDetails(cmd),
		})

	case *effects.IncubateInst:
		insts = append(insts, &target.Manual{
			Dev:     a,
			Label:   "incubate",
			Details: fmt.Sprintf("incubate at %s for %s", cmd.Temp.ToString(), cmd.Time.ToString()),
		})

	case *effects.HandleInst:
		insts = append(insts, &target.Manual{
			Dev:   a,
			Label: cmd.Group,
		})

	case *effects.PromptInst:
		insts = append(insts, &target.Prompt{
			Message: cmd.Message,
		})

	case *wtype.PRInstruction:
		insts = append(insts, &target.Manual{
			Dev:     a,
			Label:   "plate-read",
			Details: fmt.Sprintf("plate-read instruction. Options:'%s'", cmd.Options),
		})

	case *effects.QPCRInstruction:
		insts = append(insts, &target.Manual{
			Dev:     a,
			Label:   "QPCR",
			Details: fmt.Sprintf("QPCR request, definition %s, barcode %s", cmd.Definition, cmd.Barcode),
		})

	default:
		return nil, fmt.Errorf("unknown inst %T", cmd)
	}

	return insts, nil
}

func prettyMixDetails(inst *wtype.LHInstruction) string {
	if len(inst.PlateName) != 0 || len(inst.Welladdress) != 0 {
		return fmt.Sprintf("mix %q[%q]", inst.PlateName, inst.Welladdress)
	}
	return "mix"
}
