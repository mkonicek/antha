package pretty

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/auto"
)

func Timeline(out io.Writer, a *auto.Auto, result *execute.Result) error {
	g := &target.Graph{
		Insts: result.Insts,
	}

	dag := graph.Schedule(graph.Reverse(g))
	var lines []string
	for round := 1; len(dag.Roots) != 0; round += 1 {
		lines = append(lines, fmt.Sprintf("== Round %2d:\n", round))
		var next []graph.Node
		for _, n := range dag.Roots {
			inst := n.(target.Inst)
			lines = append(lines, fmt.Sprintf("    * %s\n", a.Pretty(inst)))
			next = append(next, dag.Visit(n)...)
		}

		dag.Roots = next
	}

	lines = append(lines, fmt.Sprint("== Workflow Outputs:\n"))

	for k, v := range result.Workflow.Outputs {
		var s string
		bs, err := json.Marshal(v)
		if err == nil {
			s = string(bs)
		} else {
			s = fmt.Sprintf("<cannot unmarshal: %s>", err)
		}
		lines = append(lines, fmt.Sprintf("    - %s: %s\n", k, s))
	}

	_, err := fmt.Fprint(out, strings.Join(lines, ""))
	return err
}
