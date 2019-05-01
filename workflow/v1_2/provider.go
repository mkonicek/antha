package v1_2

import (
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/workflow"
)

type V1_2WorkflowProvider struct {
	owf *workflowv1_2 // the old, v1.2 workflow to migrate
	fm  *effects.FileManager
}

func NewV1_2WorkflowProvider(wf *workflowv1_2, fm *effects.FileManager) *V1_2WorkflowProvider {
	return &V1_2WorkflowProvider{
		owf: wf,
		fm:  fm,
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

func (p *V1_2WorkflowProvider) GetElements() (*workflow.Elements, error) {
	els := &workflow.Elements{}
	instances := workflow.ElementInstances{}
	for k := range p.owf.Processes {
		name := workflow.ElementInstanceName(k)
		ei, err := p.owf.MigrateElement(p.fm, k)
		if err != nil {
			return nil, err
		}
		instances[name] = ei
	}
	els.Instances = instances
	return els, nil
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
