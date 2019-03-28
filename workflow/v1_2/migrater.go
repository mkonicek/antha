package v1_2

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory/effects"

	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/utils"
	"github.com/antha-lang/antha/workflow"
)

// Migrater handles migration to updated workflow format.
type Migrater struct {
	Logger       *logger.Logger
	Cur          *workflow.Workflow
	Old          *workflowv1_2
	GilsonDevice string // The name of a gilson device to create
	OutDir       string
	FileManager  *effects.FileManager
}

// NewMigrater creates a new migration object.
func NewMigrater(logger *logger.Logger, mergePaths []string, migratePath, outDir string, gilsonDevice string) (*Migrater, error) {
	owf, cwf, err := readWorkflows(migratePath, mergePaths)
	if err != nil {
		return nil, err
	}

	if outDir == "" {
		if outDir, err = ioutil.TempDir("", "antha-migrater"); err != nil {
			return nil, err
		}
	}
	for _, leaf := range []string{"workflow", "data"} {
		if err := os.MkdirAll(filepath.Join(outDir, leaf), 0700); err != nil {
			return nil, err
		}
	}

	dataDir := filepath.Join(outDir, "data")
	fm, err := effects.NewFileManager(dataDir, dataDir)
	if err != nil {
		return nil, err
	}

	logger.Log("outdir", outDir)

	return &Migrater{
		Logger:       logger,
		Cur:          cwf,
		Old:          owf,
		GilsonDevice: gilsonDevice,
		OutDir:       outDir,
		FileManager:  fm,
	}, nil
}

// MigrateAll perform all migration steps.
func (m *Migrater) MigrateAll() error {
	return utils.ErrorSlice{
		m.migrateJobIdAndMeta(),
		m.migrateElements(),
		m.migrateConnections(),
		m.migrateConfig(),
		m.migrateTesting(),
	}.Pack()
}

func (m *Migrater) SaveCur() error {
	p := filepath.Join(m.OutDir, "workflow", "workflow.json")
	return m.Cur.WriteToFile(p, true)
}

func (m *Migrater) migrateTesting() error {
	if len(m.Old.testOpt.Results.MixTaskResults) == 0 {
		return nil
	}

	mixChecks := make([]workflow.MixTaskCheck, 0, len(m.Old.testOpt.Results.MixTaskResults))

	for _, check := range m.Old.testOpt.Results.MixTaskResults {

		instructions, err := json.Marshal(check.Instructions)
		if err != nil {
			return err
		}

		mixChecks = append(mixChecks, workflow.MixTaskCheck{
			Instructions: json.RawMessage(instructions),
			Outputs:      check.Outputs,
			TimeEstimate: check.TimeEstimate,
		})
	}

	m.Cur.Testing = &workflow.Testing{
		MixTaskChecks: mixChecks,
	}

	return nil
}

func (m *Migrater) migrateConfig() error {
	m.migrateGlobalMixerConfig()
	m.migrateGilsonConfigs()
	return nil
}

func (m *Migrater) migrateGilsonConfigs() {
	// If no gilson device specified, do nothing.
	if m.GilsonDevice == "" {
		return
	}

	devId := workflow.DeviceInstanceID(m.GilsonDevice)

	if _, found := m.Cur.Config.GilsonPipetMax.Devices[devId]; found {
		m.Logger.Log("warning", fmt.Sprintf("Gilson device %s already exists, and will have configuration replaced with migrated configuration.", m.GilsonDevice))
	}

	m.Cur.Config.GilsonPipetMax.Devices[devId] = m.migrateGilsonConfig()
}

func (m *Migrater) migrateGlobalMixerConfig() {
	m.Cur.Config.GlobalMixer = workflow.GlobalMixerConfig{
		CustomPolicyRuleSet:      m.Old.Config.CustomPolicyRuleSet,
		IgnorePhysicalSimulation: m.Old.Config.IgnorePhysicalSimulation,
		InputPlates:              m.Old.Config.InputPlates,
		PrintInstructions:        m.Old.Config.PrintInstructions,
		UseDriverTipTracking:     m.Old.Config.UseDriverTipTracking,
	}
}

func (m *Migrater) migrateLayoutPreferences() *workflow.LayoutOpt {
	return &workflow.LayoutOpt{
		Inputs:    m.Old.Config.DriverSpecificInputPreferences,
		Outputs:   m.Old.Config.DriverSpecificOutputPreferences,
		Tipboxes:  m.Old.Config.DriverSpecificTipPreferences,
		Tipwastes: m.Old.Config.DriverSpecificTipWastePreferences,
		Washes:    m.Old.Config.DriverSpecificWashPreferences,
	}
}

func updatePlateTypes(names []string) []wtype.PlateTypeName {
	ptnames := make([]wtype.PlateTypeName, len(names))
	for i, v := range names {
		ptnames[i] = wtype.PlateTypeName(v)
	}
	return ptnames
}

func (m *Migrater) migrateGilsonConfig() *workflow.GilsonPipetMaxInstanceConfig {
	config := workflow.GilsonPipetMaxInstanceConfig{}
	config.InputPlateTypes = updatePlateTypes(m.Old.Config.InputPlateTypes)
	config.MaxPlates = m.Old.Config.MaxPlates
	config.MaxWells = m.Old.Config.MaxWells
	config.OutputPlateTypes = updatePlateTypes(m.Old.Config.OutputPlateTypes)
	config.ResidualVolumeWeight = m.Old.Config.ResidualVolumeWeight
	config.TipTypes = m.Old.Config.TipTypes
	config.LayoutPreferences = m.migrateLayoutPreferences()
	return &config
}

func (m *Migrater) migrateElements() error {
	return utils.ErrorSlice{
		m.migrateElementInstances(),
		m.migrateElementTypes(),
	}.Pack()
}

func (m *Migrater) migrateElementInstances() error {
	for k := range m.Old.Processes {
		name := workflow.ElementInstanceName(k)
		ei, err := m.Old.MigrateElement(m.FileManager, k)
		if err != nil {
			return err
		}
		m.Cur.Elements.Instances[name] = ei
	}
	return nil
}

func uniqueElementType(types workflow.ElementTypesByRepository, name workflow.ElementTypeName) (*workflow.ElementType, error) {
	var et *workflow.ElementType
	for _, rmap := range types {
		if v, found := rmap[name]; found {
			if et != nil {
				return nil, fmt.Errorf("element type %v is found in multiple repositories", name)
			}
			et = &v
		}
	}

	if et == nil {
		return nil, fmt.Errorf("element type %v could not be found in the supplied repositories", name)
	}
	return et, nil
}

func (m *Migrater) migrateElementTypes() error {
	repoMaps, err := m.Cur.Repositories.FindAllElementTypes()
	if err != nil {
		return err
	}

	seen := make(map[string]struct{}, len(m.Old.Processes))
	ets := make(workflow.ElementTypes, 0, len(m.Old.Processes))
	for _, v := range m.Old.Processes {
		if _, found := seen[v.Component]; found {
			continue
		}

		seen[v.Component] = struct{}{}
		et, err := uniqueElementType(repoMaps, workflow.ElementTypeName(v.Component))
		if err != nil {
			return err
		}
		ets = append(ets, et)
	}

	m.Cur.Elements.Types = ets
	return nil
}

func (m *Migrater) migrateConnections() error {
	connections := make(workflow.ElementInstancesConnections, 0, len(m.Old.Connections))
	for _, c := range m.Old.Connections {
		connections = append(connections, workflow.ElementConnection{
			Source: workflow.ElementSocket{
				ElementInstance: workflow.ElementInstanceName(c.Src.Process),
				ParameterName:   workflow.ElementParameterName(c.Src.Port),
			},
			Target: workflow.ElementSocket{
				ElementInstance: workflow.ElementInstanceName(c.Tgt.Process),
				ParameterName:   workflow.ElementParameterName(c.Tgt.Port),
			},
		})
	}

	m.Cur.Elements.InstancesConnections = connections
	return nil
}

func (m *Migrater) migrateJobIdAndMeta() error {
	if m.Old.Properties.Name != "" {
		m.Cur.JobId = workflow.JobId(m.Old.Properties.Name)
	}
	if desc := m.Old.Properties.Description; desc != "" {
		m.Cur.Meta.Rest["Description"] = desc
	}
	return nil
}

func readWorkflows(migrate string, merges []string) (*workflowv1_2, *workflow.Workflow, error) {
	rs, err := workflow.ReadersFromPaths(append(merges, migrate))

	if err != nil {
		return nil, nil, err
	}

	cwf, err := workflow.WorkflowFromReaders(rs[:len(rs)-1]...)
	if err != nil {
		return nil, nil, err
	}

	owf, err := readWorkflowV1_2(rs[len(rs)-1])
	if err != nil {
		return nil, nil, err
	}
	return owf, cwf, nil
}

func readWorkflowV1_2(r io.ReadCloser) (*workflowv1_2, error) {
	defer r.Close()
	wf := &workflowv1_2{}
	if err := json.NewDecoder(r).Decode(wf); err != nil {
		return nil, err
	}
	return wf, nil
}

// Validate performs standard validation on the migrated version of the file (Migration should be performed first)
func (m *Migrater) ValidateCur() error {
	return m.Cur.Validate()
}

// ValidateOld checks that the workflow supplied as input is of the correct version
func (m *Migrater) ValidateOld() error {
	if !(m.Old.Version == "1.2.0" || m.Old.Version == "") {
		return fmt.Errorf("Unexpected version '%s'", m.Old.Version)
	}
	return nil
}
