package provider

import (
	"github.com/antha-lang/antha/workflow"
)

// WorkflowProvider is an entity capable of providing the necessary information
// to generate a workflow, e.g. a v1.2 workflow JSON file
type WorkflowProvider interface {
	GetWorkflowID() (workflow.BasicId, error)
	GetMeta() (workflow.Meta, error)
	GetRepositories() (workflow.Repositories, error)
	GetElements() (workflow.Elements, error)
	GetInventory() (workflow.Inventory, error)
	GetConfig() (workflow.Config, error)
	GetTesting() (workflow.Testing, error)
}
