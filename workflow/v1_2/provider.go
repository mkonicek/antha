package v1_2

import (
	"github.com/antha-lang/antha/workflow"
)

type V1_2WorkflowProvider struct {
	owf *workflowv1_2
}

func NewV1_2WorkflowProvider(wf *workflowv1_2) *V1_2WorkflowProvider {
	return &V1_2WorkflowProvider{
		owf: wf,
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
	return nil, nil
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
