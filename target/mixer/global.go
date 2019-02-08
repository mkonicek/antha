package mixer

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
	"github.com/antha-lang/antha/workflow"
)

type GlobalMixerConfig struct {
	*workflow.GlobalMixerConfig
}

func NewGlobalMixerConfig(cfg *workflow.GlobalMixerConfig) *GlobalMixerConfig {
	return &GlobalMixerConfig{
		GlobalMixerConfig: cfg,
	}
}

func (cfg *GlobalMixerConfig) Validate(inv *inventory.Inventory) error {
	for _, plates := range [][]*wtype.Plate{cfg.InputPlates, cfg.OutputPlates} {
		for _, plate := range plates {
			if _, err := inv.PlateTypes.NewPlateType(plate.Type); err != nil {
				return err
			}
		}
	}
	return nil
}

func (cfg *GlobalMixerConfig) ApplyToLHRequest(req *liquidhandling.LHRequest) error {
	if cfg.CustomPolicyRuleSet != nil {
		req.AddUserPolicies(cfg.CustomPolicyRuleSet)
	}
	if err := req.PolicyManager.SetOption("USE_DRIVER_TIP_TRACKING", cfg.UseDriverTipTracking); err != nil {
		return err
	}
	req.Options.PrintInstructions = cfg.PrintInstructions
	req.Options.IgnorePhysicalSimulation = cfg.IgnorePhysicalSimulation
	return nil
}
