package pretty

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/auto"
	"golang.org/x/net/context"
)

func shouldWait(inst ast.Inst) bool {
	switch inst.(type) {
	case *target.Run:
		return true
	}
	return false
}

// Run executes an execute.Result against the given auto target.
func Run(out io.Writer, in io.Reader, a *auto.Auto, result *execute.Result) error {
	if _, err := fmt.Fprintf(out, "== Running Workflow:\n"); err != nil {
		return err
	}

	bin := bufio.NewReader(in)
	ctx := context.Background()
	for _, inst := range result.Insts {
		if _, err := fmt.Fprintf(out, "    * %s", a.Pretty(inst)); err != nil {
			return err
		}

		run := true
		if shouldWait(inst) {
			if _, err := fmt.Fprintf(out, " (Run? [yes,skip]) "); err != nil {
				return err
			} else if s, err := bin.ReadString('\n'); err != nil {
				return err
			} else {
				run = strings.HasPrefix(strings.ToLower(s), "yes")
			}
		}

		if run {
			if err := a.Execute(ctx, inst); err != nil {
				fmt.Fprintf(out, " [FAIL]\n") // nolint
				return err
			}
		}

		if _, err := fmt.Fprintf(out, " [OK]\n"); err != nil {
			return err
		}
	}
	return nil
}
