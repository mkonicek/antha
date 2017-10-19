package comp

import (
	"fmt"

	"github.com/antha-lang/antha/component"
)

// TODO: replace with api defintions

// A Port is an input to a component
type Port struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
}

// A Component is an element to execute
type Component struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	InPorts     []Port `json:"in_ports"`
	OutPorts    []Port `json:"out_ports"`
}

// New creates Component messages from components
func New(lib []component.Component) ([]Component, error) {
	var cs []Component
	for _, v := range lib {
		c := Component{
			ID:          v.Name,
			Name:        v.Name,
			Description: v.Description.Desc,
			Path:        v.Description.Path,
		}
		for _, p := range v.Description.Params {
			port := Port{
				Name:        p.Name,
				Type:        p.Type,
				Description: p.Desc,
				Kind:        p.Kind,
			}
			switch p.Kind {
			case "Outputs", "Data":
				c.OutPorts = append(c.OutPorts, port)
			case "Inputs", "Parameters":
				c.InPorts = append(c.InPorts, port)
			default:
				return nil, fmt.Errorf("unknown parameter kind %q", p.Kind)
			}
		}
		cs = append(cs, c)
	}
	return cs, nil
}
