package pretty

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/target"
)

func prettyInst(inst *target.Manual) string {
	return fmt.Sprintf("[%s] %s", inst.Label, strings.Replace(inst.Details, "\n", "; ", -1))
}

func summarize(inst target.Inst) (string, error) {
	switch inst := inst.(type) {
	case target.RunInst:
		return fmt.Sprintf("Run file (size: %d)", len(inst.Data().Tarball)), nil
	case *target.Manual:
		return prettyInst(inst), nil
	case *target.Wait:
		return "", nil
	case *target.CmpError:
		return "", fmt.Errorf("Planning error: %s", inst.Error)
	default:
		return "", fmt.Errorf("unknown inst %T", inst)
	}
}

func Timeline(out io.Writer, result *execute.Result) error {
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
			if s, err := summarize(inst); err != nil {
				return err
			} else {
				lines = append(lines, fmt.Sprintf("    * %s\n", s))
				next = append(next, dag.Visit(n)...)
			}
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
