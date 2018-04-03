package cmd

import (
	"encoding/json"
	"os"

	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/target/mixer"
	"github.com/antha-lang/antha/workflow"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var convertFromUICmd = &cobra.Command{
	Use:   "from-ui",
	Short: "Convert UI bundle to antha bundle",
	RunE:  convertFromUI,
}

type anthaBundle struct {
	workflow.Desc
	execute.RawParams
}

func convertConfigToAntha(in *uiMixerConfig) *mixer.Opt {
	if in == nil {
		return nil
	}

	getFloat := func(v float64) *float64 {
		if v == 0.0 {
			return nil
		}
		return &v
	}

	ret := &mixer.Opt{
		MaxPlates:                         getFloat(in.MaxPlates),
		MaxWells:                          getFloat(in.MaxWells),
		ResidualVolumeWeight:              getFloat(in.ResidualVolumeWeight),
		InputPlateTypes:                   in.InputPlateTypes,
		OutputPlateTypes:                  in.OutputPlateTypes,
		TipTypes:                          in.TipTypes,
		PlanningVersion:                   in.PlanningVersion,
		DriverSpecificInputPreferences:    in.DriverSpecificInputPreferences,
		DriverSpecificOutputPreferences:   in.DriverSpecificOutputPreferences,
		DriverSpecificTipPreferences:      in.DriverSpecificTipPreferences,
		DriverSpecificTipWastePreferences: in.DriverSpecificTipWastePreferences,
		DriverSpecificWashPreferences:     in.DriverSpecificWashPreferences,
		ModelEvaporation:                  in.ModelEvaporation,
		OutputSort:                        in.OutputSort,
		PrintInstructions:                 in.PrintInstructions,
		UseDriverTipTracking:              in.UseDriverTipTracking,
		LegacyVolume:                      in.LegacyVolume,
	}

	return ret
}

func convertProcessesToAntha(in map[string]uiProcess) map[string]workflow.Process {
	ret := make(map[string]workflow.Process)

	for k, v := range in {
		ret[k] = workflow.Process{
			Component: v.Component,
		}
	}

	return ret
}

func convertFromUI(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	inFile, err := os.Open(viper.GetString("bundle"))
	if err != nil {
		return err
	}
	defer inFile.Close() // nolint

	var ui uiBundle
	dec := json.NewDecoder(inFile)
	if err := dec.Decode(&ui); err != nil {
		return err
	}

	var bundle anthaBundle
	bundle.Connections = ui.Connections
	bundle.Parameters = ui.Parameters
	bundle.Processes = convertProcessesToAntha(ui.Processes)
	bundle.Config = convertConfigToAntha(ui.Config)

	outFile, err := os.Create(viper.GetString("output"))
	if err != nil {
		return err
	}
	defer outFile.Close() // nolint

	enc := json.NewEncoder(outFile)
	enc.SetIndent("", "  ")
	return enc.Encode(bundle)
}

func init() {
	c := convertFromUICmd
	convertCmd.AddCommand(c)

	flags := c.Flags()
	flags.String("bundle", "", "Input bundle")
	flags.String("output", "output.json", "Output bundle in antha format")
}
