package pretty

import (
	"fmt"
	"io"
	"strings"

	"github.com/antha-lang/antha/cmd/antha/comp"
)

// Components returns a pretty printed version of comp.Component
func Components(out io.Writer, cs []comp.Component) error {
	var lines []string

	for _, c := range cs {
		lines = append(lines, fmt.Sprintf("%s:\n", c.Name))
		lines = append(lines, fmt.Sprintf("\tInputs:\n"))
		for _, p := range c.InPorts {
			lines = append(lines, fmt.Sprintf("\t\t%s %s\n", p.Name, p.Type))
		}
		lines = append(lines, fmt.Sprintf("\tOutputs:\n"))
		for _, p := range c.OutPorts {
			lines = append(lines, fmt.Sprintf("\t\t%s %s\n", p.Name, p.Type))
		}
		lines = append(lines, fmt.Sprintf("\tPath:\n"))
		lines = append(lines, fmt.Sprintf("\t\t%s\n", c.Path))
	}

	_, err := fmt.Fprint(out, strings.Join(lines, ""))
	return err
}
