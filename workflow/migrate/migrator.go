package migrate

import (
	"github.com/antha-lang/antha/workflow"
)

// Migrator migrates data from a previous format to the v2.0 format
type Migrator struct {
	provider WorkflowProvider
}

// NewMigrator creates and returns a migrator
func NewMigrator(provider WorkflowProvider) *Migrator {
	return &Migrator{
		provider: provider,
	}
}

// Workflow returns the workflow resulting from exercising this migrator
func (m *Migrator) Workflow() (*workflow.Workflow, error) {

	wf := workflow.EmptyWorkflow()

	id, err := m.provider.GetWorkflowID()
	if err != nil {
		return nil, err
	}
	wf.WorkflowId = id

	meta, err := m.provider.GetMeta()
	if err != nil {
		return nil, err
	}
	wf.Meta = meta

	repos, err := m.provider.GetRepositories()
	if err != nil {
		return nil, err
	}
	wf.Repositories = repos

	elements, err := m.provider.GetElements()
	if err != nil {
		return nil, err
	}
	wf.Elements = elements

	inventory, err := m.provider.GetInventory()
	if err != nil {
		return nil, err
	}
	wf.Inventory = inventory

	config, err := m.provider.GetConfig()
	if err != nil {
		return nil, err
	}
	wf.Config = config

	testing, err := m.provider.GetTesting()
	if err != nil {
		return nil, err
	}
	wf.Testing = testing

	return wf, nil
}
