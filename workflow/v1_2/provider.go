package v1_2

import (
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/workflow"
)

type V1_2WorkflowProvider struct {
	owf     *workflowv1_2 // the old, v1.2 workflow to migrate
	fm      *effects.FileManager
	repoMap workflow.ElementTypesByRepository
}

func NewV1_2WorkflowProvider(
	wf *workflowv1_2,
	fm *effects.FileManager,
	repoMap workflow.ElementTypesByRepository,
) *V1_2WorkflowProvider {
	return &V1_2WorkflowProvider{
		owf:     wf,
		fm:      fm,
		repoMap: repoMap,
	}
}

func (p *V1_2WorkflowProvider) GetMeta() (*workflow.Meta, error) {
	meta := &workflow.Meta{}
	if p.owf.Properties.Name != "" {
		meta.Name = p.owf.Properties.Name
	}
	return meta, nil
}

func (p *V1_2WorkflowProvider) GetRepositories() (*workflow.Repositories, error) {
	return nil, nil
}

func (p *V1_2WorkflowProvider) getElementInstances() (workflow.ElementInstances, error) {
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

func (p *V1_2WorkflowProvider) getElementTypes() (workflow.ElementTypes, error) {
	seen := make(map[string]struct{}, len(p.owf.Processes))
	types := make(workflow.ElementTypes, 0, len(p.owf.Processes))
	for _, v := range p.owf.Processes {
		if _, found := seen[v.Component]; found {
			continue
		}

		seen[v.Component] = struct{}{}
		et, err := uniqueElementType(p.repoMap, workflow.ElementTypeName(v.Component))
		if err != nil {
			return nil, err
		}
		types = append(types, et)
	}

	return types, nil
}

func (p *V1_2WorkflowProvider) getElementConnections() (workflow.ElementInstancesConnections, error) {
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

func (p *V1_2WorkflowProvider) GetElements() (*workflow.Elements, error) {
	instances, err := p.getElementInstances()
	if err != nil {
		return nil, err
	}

	types, err := p.getElementTypes()
	if err != nil {
		return nil, err
	}

	connections, err := p.getElementConnections()
	if err != nil {
		return nil, err
	}

	return &workflow.Elements{
		Instances:            instances,
		Types:                types,
		InstancesConnections: connections,
	}, nil
}

func (p *V1_2WorkflowProvider) GetInventory() (*workflow.Inventory, error) {
	return nil, nil
}

func (p *V1_2WorkflowProvider) GetConfig() (*workflow.Config, error) {
	return nil, nil
}

func (p *V1_2WorkflowProvider) GetTesting() (*workflow.Testing, error) {
	return nil, nil
}
