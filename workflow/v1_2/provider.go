package v1_2

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wtype/liquidtype"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/workflow"
	"github.com/antha-lang/antha/workflow/migrate"
)

type Provider struct {
	owf              *workflowv1_2 // the old, v1.2 workflow to migrate
	fm               *effects.FileManager
	repoMap          workflow.ElementTypesByRepository
	gilsonDeviceName string
	logger           *logger.Logger
}

func NewProvider(
	oldWorkflowReader io.Reader,
	fm *effects.FileManager,
	repoMap workflow.ElementTypesByRepository,
	gilsonDeviceName string,
	logger *logger.Logger,
) (*Provider, error) {
	d := json.NewDecoder(oldWorkflowReader)
	wf := &workflowv1_2{}
	err := d.Decode(wf)
	if err != nil {
		return nil, err
	}

	expectedVersion := "1.2.0"
	if wf.Version != expectedVersion {
		return nil, fmt.Errorf("Invalid version in old workflow: expected %v, got %v", expectedVersion, wf.Version)
	}

	return &Provider{
		owf:              wf,
		fm:               fm,
		repoMap:          repoMap,
		gilsonDeviceName: gilsonDeviceName,
		logger:           logger,
	}, nil
}

func (p *Provider) GetWorkflowID() (workflow.BasicId, error) {
	if id, err := workflow.RandomBasicId(""); err != nil {
		return "", err
	} else {
		return id, nil
	}
}

func (p *Provider) GetMeta() (workflow.Meta, error) {
	meta := workflow.EmptyMeta()
	if p.owf.Properties.Name != "" {
		meta.Name = p.owf.Properties.Name
	}
	if p.owf.Properties.Description != "" {
		meta.Set("Description", p.owf.Properties.Description)
	}
	return meta, nil
}

func (p *Provider) GetRepositories() (workflow.Repositories, error) {
	return workflow.Repositories{}, nil
}

func (p *Provider) getElementInstances() (workflow.ElementInstances, error) {
	instances := workflow.ElementInstances{}
	for k := range p.owf.Processes {
		name := workflow.ElementInstanceName(k)
		ei, err := p.owf.MigrateElement(p.fm, k)
		if err != nil {
			return nil, err
		}
		instances[name] = ei
	}
	return instances, nil
}

func (p *Provider) getElementTypes() (workflow.ElementTypes, error) {
	seen := make(map[string]struct{}, len(p.owf.Processes))
	types := make(workflow.ElementTypes, 0, len(p.owf.Processes))
	for _, v := range p.owf.Processes {
		if _, found := seen[v.Component]; found {
			continue
		}

		seen[v.Component] = struct{}{}
		et, err := migrate.UniqueElementType(p.repoMap, workflow.ElementTypeName(v.Component))
		if err != nil {
			return nil, err
		}
		types = append(types, et)
	}

	return types, nil
}

func (p *Provider) getElementConnections() (workflow.ElementInstancesConnections, error) {
	connections := make(workflow.ElementInstancesConnections, 0, len(p.owf.Connections))
	for _, c := range p.owf.Connections {
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
	return connections, nil
}

func (p *Provider) GetElements() (workflow.Elements, error) {
	instances, err := p.getElementInstances()
	if err != nil {
		return workflow.Elements{}, err
	}

	types, err := p.getElementTypes()
	if err != nil {
		return workflow.Elements{}, err
	}

	connections, err := p.getElementConnections()
	if err != nil {
		return workflow.Elements{}, err
	}

	return workflow.Elements{
		Instances:            instances,
		Types:                types,
		InstancesConnections: connections,
	}, nil
}

func (p *Provider) GetInventory() (workflow.Inventory, error) {
	return workflow.Inventory{}, nil
}

func (p *Provider) getGlobalMixerConfig() (workflow.GlobalMixerConfig, error) {
	customPolicyRuleSet := p.owf.Config.CustomPolicyRuleSet
	if p.owf.Config.LiquidHandlingPolicyXlsxJmpFile != nil {
		policyMap, err := liquidtype.PolicyMakerFromBytes(p.owf.Config.LiquidHandlingPolicyXlsxJmpFile, wtype.PolicyName(liquidtype.BASEPolicy))
		if err != nil {
			return workflow.GlobalMixerConfig{}, err
		}
		lhpr := wtype.NewLHPolicyRuleSet()
		lhpr, err = wtype.AddUniversalRules(lhpr, policyMap)
		if err != nil {
			return workflow.GlobalMixerConfig{}, err
		}
		if customPolicyRuleSet == nil {
			customPolicyRuleSet = lhpr
		} else {
			customPolicyRuleSet.MergeWith(lhpr)
		}
	}

	return workflow.GlobalMixerConfig{
		CustomPolicyRuleSet:      customPolicyRuleSet,
		IgnorePhysicalSimulation: p.owf.Config.IgnorePhysicalSimulation,
		InputPlates:              p.owf.Config.InputPlates,
		PrintInstructions:        p.owf.Config.PrintInstructions,
		UseDriverTipTracking:     p.owf.Config.UseDriverTipTracking,
	}, nil
}

func (p *Provider) getLayoutPreferences() *workflow.LayoutOpt {
	if p.owf.Config == nil {
		return nil
	}
	return &workflow.LayoutOpt{
		Inputs:    p.owf.Config.DriverSpecificInputPreferences,
		Outputs:   p.owf.Config.DriverSpecificOutputPreferences,
		Tipboxes:  p.owf.Config.DriverSpecificTipPreferences,
		Tipwastes: p.owf.Config.DriverSpecificTipWastePreferences,
		Washes:    p.owf.Config.DriverSpecificWashPreferences,
	}
}

func (p *Provider) getGilsonPipetMaxInstanceConfig() (*workflow.GilsonPipetMaxInstanceConfig, error) {
	config := workflow.GilsonPipetMaxInstanceConfig{}
	if p.owf.Config != nil {
		config.InputPlateTypes = migrate.UpdatePlateTypes(p.owf.Config.InputPlateTypes)
		config.MaxPlates = p.owf.Config.MaxPlates
		config.MaxWells = p.owf.Config.MaxWells
		config.OutputPlateTypes = migrate.UpdatePlateTypes(p.owf.Config.OutputPlateTypes)
		config.ResidualVolumeWeight = p.owf.Config.ResidualVolumeWeight
		config.TipTypes = p.owf.Config.TipTypes
		config.LayoutPreferences = p.getLayoutPreferences()
	}
	return &config, nil
}

func (p *Provider) getGilsonPipetMaxConfig() (workflow.GilsonPipetMaxConfig, error) {
	if p.gilsonDeviceName == "" {
		return workflow.GilsonPipetMaxConfig{}, nil
	}

	devices := map[workflow.DeviceInstanceID]*workflow.GilsonPipetMaxInstanceConfig{}
	devID := workflow.DeviceInstanceID(p.gilsonDeviceName)

	devConfig, err := p.getGilsonPipetMaxInstanceConfig()
	if err != nil {
		return workflow.GilsonPipetMaxConfig{}, err
	}

	devices[devID] = devConfig

	return workflow.GilsonPipetMaxConfig{
		Devices: devices,
	}, nil
}

func (p *Provider) GetConfig() (workflow.Config, error) {
	gmc, err := p.getGlobalMixerConfig()
	if err != nil {
		return workflow.Config{}, err
	}

	gpmc, err := p.getGilsonPipetMaxConfig()
	if err != nil {
		return workflow.Config{}, err
	}

	return workflow.Config{
		GlobalMixer:    gmc,
		GilsonPipetMax: gpmc,
	}, nil
}

func (p *Provider) GetTesting() (workflow.Testing, error) {
	if len(p.owf.testOpt.Results.MixTaskResults) == 0 {
		return workflow.Testing{}, nil
	}

	mixChecks := make([]workflow.MixTaskCheck, 0, len(p.owf.testOpt.Results.MixTaskResults))

	for _, check := range p.owf.testOpt.Results.MixTaskResults {

		instructions, err := json.Marshal(check.Instructions)
		if err != nil {
			return workflow.Testing{}, err
		}

		mixChecks = append(mixChecks, workflow.MixTaskCheck{
			Instructions: json.RawMessage(instructions),
			Outputs:      check.Outputs,
			TimeEstimate: check.TimeEstimate,
		})
	}

	return workflow.Testing{
		MixTaskChecks: mixChecks,
	}, nil
}
