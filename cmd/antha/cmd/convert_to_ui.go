package cmd

import (
	"encoding/json"
	"os"
	"regexp"

	"github.com/antha-lang/antha/execute/executeutil"
	"github.com/antha-lang/antha/target/mixer"
	"github.com/antha-lang/antha/workflow"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var convertToUICmd = &cobra.Command{
	Use:   "to-ui",
	Short: "Convert antha parameters to UI bundle",
	RunE:  convertToUI,
}

type uiParameters map[string]map[string]json.RawMessage

// NB: inconsistent naming is required
type uiBundle struct {
	Processes   map[string]uiProcess  `json:"Processes"`
	Connections []workflow.Connection `json:"connections"`
	Parameters  uiParameters          `json:"Parameters"`
	Config      *uiMixerConfig        `json:"Config"`
	Version     string                `json:"version"`
	Properties  *uiProperties         `json:"Properties"`
}

type uiProcess struct {
	Component string      `json:"component"`
	Metadata  *uiMetadata `json:"metadata"`
}

type uiMetadata struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type uiProperties struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Serialization format for UI bundles. Would like to move to protobuf
// definition in api, but the UI format is not canonical protobuf
type uiMixerConfig struct {
	MaxPlates                         float64  `json:"maxPlates,omitempty"`
	MaxWells                          float64  `json:"maxWells,omitempty"`
	ResidualVolumeWeight              float64  `json:"residualVolumeWeight,omitempty"`
	InputPlateTypes                   []string `json:"inputPlateTypes"`
	OutputPlateTypes                  []string `json:"outputPlateTypes"`
	TipTypes                          []string `json:"tipTypes"`
	PlanningVersion                   string   `json:"planningVersion"`
	DriverSpecificInputPreferences    []string `json:"driverSpecificInputPreferences"`
	DriverSpecificOutputPreferences   []string `json:"driverSpecificOutputPreferences"`
	DriverSpecificTipPreferences      []string `json:"driverSpecificTipPreferences"`
	DriverSpecificTipWastePreferences []string `json:"driverSpecificTipWastePreferences"`
	DriverSpecificWashPreferences     []string `json:"driverSpecificWashPreferences"`
	ModelEvaporation                  bool     `json:"modelEvaporation,omitempty"`
	OutputSort                        bool     `json:"outputSort,omitempty"`
	PrintInstructions                 bool     `json:"printInstructions,omitempty"`
	UseDriverTipTracking              bool     `json:"useDriverTipTracking,omitempty"`
	LegacyVolume                      bool     `json:"legacyVolume,omitempty"`
}

func convertConfigToUI(in *mixer.Opt) *uiMixerConfig {
	if in == nil {
		return nil
	}

	getFloat := func(v *float64) float64 {
		if v == nil {
			return 0.0
		}
		return *v
	}

	getSlice := func(xs []string) []string {
		if xs == nil {
			return []string{}
		}
		return xs
	}

	ret := &uiMixerConfig{
		MaxPlates:                         getFloat(in.MaxPlates),
		MaxWells:                          getFloat(in.MaxWells),
		ResidualVolumeWeight:              getFloat(in.ResidualVolumeWeight),
		InputPlateTypes:                   getSlice(in.InputPlateTypes),
		OutputPlateTypes:                  getSlice(in.OutputPlateTypes),
		TipTypes:                          getSlice(in.TipTypes),
		PlanningVersion:                   in.PlanningVersion,
		DriverSpecificInputPreferences:    getSlice(in.DriverSpecificInputPreferences),
		DriverSpecificOutputPreferences:   getSlice(in.DriverSpecificOutputPreferences),
		DriverSpecificTipPreferences:      getSlice(in.DriverSpecificTipPreferences),
		DriverSpecificTipWastePreferences: getSlice(in.DriverSpecificTipWastePreferences),
		DriverSpecificWashPreferences:     getSlice(in.DriverSpecificWashPreferences),
		ModelEvaporation:                  in.ModelEvaporation,
		OutputSort:                        in.OutputSort,
		PrintInstructions:                 in.PrintInstructions,
		UseDriverTipTracking:              in.UseDriverTipTracking,
		LegacyVolume:                      in.LegacyVolume,
	}

	return ret
}

func convertProcessesToUI(in map[string]workflow.Process) (map[string]uiProcess, map[string]string) {
	ret := make(map[string]uiProcess)
	rename := make(map[string]string)

	pat := regexp.MustCompile(`.*_run(\d+)$`)

	for k, v := range in {
		name := k
		if !pat.MatchString(name) {
			name += "_run1"
		}
		ret[name] = uiProcess{
			Component: v.Component,
			// Must be non-zero
			Metadata: &uiMetadata{
				X: 1,
				Y: 1,
			},
		}
		rename[k] = name
	}

	return ret, rename
}

func convertParametersToUI(params uiParameters, rename map[string]string) uiParameters {
	ret := make(uiParameters)
	for k, v := range params {
		ret[rename[k]] = v
	}
	return ret
}

func convertConnectionsToUI(in []workflow.Connection, rename map[string]string) (ret []workflow.Connection) {
	for _, v := range in {
		ret = append(ret, workflow.Connection{
			Src: workflow.Port{
				Process: rename[v.Src.Process],
				Port:    v.Src.Port,
			},
			Tgt: workflow.Port{
				Process: rename[v.Tgt.Process],
				Port:    v.Tgt.Port,
			},
		})
	}
	return
}

func convertToUI(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	bundle, err := executeutil.UnmarshalSingle(viper.GetString("bundle"), viper.GetString("workflow"), viper.GetString("parameters"))
	if err != nil {
		return err
	}

	processes, rename := convertProcessesToUI(bundle.Processes)
	ui := &uiBundle{
		Processes:   processes,
		Connections: convertConnectionsToUI(bundle.Connections, rename),
		Parameters:  convertParametersToUI(bundle.Parameters, rename),
		Config:      convertConfigToUI(bundle.Config),
		Version:     "1.2.0",
		// Empty properties is required
		Properties: &uiProperties{},
	}

	f, err := os.Create(viper.GetString("output"))
	if err != nil {
		return err
	}
	defer f.Close() // nolint

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(ui)
}

func init() {
	c := convertToUICmd
	convertCmd.AddCommand(c)

	flags := c.Flags()
	flags.String("parameters", "parameters.json", "Parameters to workflow")
	flags.String("workflow", "workflow.json", "Workflow definition file")
	flags.String("bundle", "", "Input bundle with parameters and workflow together (overrides parameter and workflow arguments)")
	flags.String("output", "output.json", "Output bundle in UI format")
}
