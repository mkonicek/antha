package auto

import (
	"fmt"
	"strings"

	"github.com/antha-lang/antha/target"
)

// Pretty returns a human description of instruction
func (a *Auto) Pretty(inst target.Inst) string {
	switch inst := inst.(type) {
	case *target.Mix:
		return prettyMix(inst)
	case *target.Run:
		return prettyRun(inst)
	case *target.Manual:
		return prettyManual(inst)
	case *target.Wait:
		return "Wait"
	default:
		return fmt.Sprintf("unknown instruction %T", inst)
	}
}

func prettyManual(inst *target.Manual) string {
	return fmt.Sprintf("[%s] %s", inst.Label, strings.Replace(inst.Details, "\n", "; ", -1))
}

func prettyMix(inst *target.Mix) string {
	return fmt.Sprintf("[mix] (size: %d)", len(inst.Files.Tarball))
}

func prettyRun(inst *target.Run) string {
	return fmt.Sprintf("[run] %s", inst.Label)
}
