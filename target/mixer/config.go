package mixer

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/driver/liquidhandling/client"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/workflow"
)

type GlobalMixerConfig struct {
	*workflow.GlobalMixerConfig
}

func (cfg *GlobalMixerConfig) validate(inv *inventory.Inventory) error {
	for _, plates := range [][]wtype.Plate{cfg.InputPlates, cfg.OutputPlates} {
		for _, plate := range plates {
			if _, err := inv.PlateTypes.NewPlateType(plate.Type); err != nil {
				return err
			}
		}
	}
	return nil
}

type GilsonPipetMaxInstances map[workflow.DeviceInstanceID]*GilsonPipetMaxInstanceConfig

type GilsonPipetMaxInstanceConfig struct {
	*workflow.GilsonPipetMaxInstanceConfig
	base   *BaseMixer
	Driver *client.LowLevelClient
}

func GilsonPipetMaxInstancesFromWorkflow(wf *workflow.Workflow) (GilsonPipetMaxInstances, error) {
	devices := wf.Config.GilsonPipetMax.Devices

	res := make(GilsonPipetMaxInstances, len(devices))
	for id, cfgWf := range devices {
		cfg := &GilsonPipetMaxInstanceConfig{
			GilsonPipetMaxInstanceConfig: cfgWf,
			base:                         NewBaseMixer(cfgWf.Connection, "GilsonPipetmax"),
		}
		if err := cfg.Connect(); err != nil {
			return nil, fmt.Errorf("Error when connecting to GilsonPipetmax at %s: %v", cfgWf.Connection, err)
		}
		res[id] = cfg
	}
	return res, nil
}

func (cfg *GilsonPipetMaxInstanceConfig) validate(id workflow.DeviceInstanceID, inv *inventory.Inventory) error {
	/*
		for _, ptns := range [][]wtype.PlateTypeName{cfg.InputPlateTypes, cfg.OutputPlateTypes} {
			for _, ptn := range ptns {
				if _, err := inv.PlateTypes.NewPlateType(ptn); err != nil {
					return err
				}
			}
		}

		// TODO: this check waste type creating first new tip boxes that we
		// throw away. We should have a tip box type.
		for _, tt := range cfg.TipTypes {
			if _, err := inv.TipBoxes.NewTipbox(tt); err != nil {
				return err
			}
		}
	*/
	return nil
}

func (cfg *GilsonPipetMaxInstanceConfig) Connect() error {
	if cfg.Driver == nil {
		if conn, err := cfg.base.ConnectInit(); err != nil {
			return err
		} else if conn != nil {
			cfg.Driver = client.NewLowLevelClientFromConn(conn)
		}
	}
	return nil
}
