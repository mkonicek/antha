package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/workflow"
)

func main() {
	l := logger.NewLogger()

	args := os.Args
	if len(args) < 2 {
		l.Fatal(errors.New("Subcommand needed"))
	}

	subCmds := map[string]func(*logger.Logger, []string){
		"find": find,
	}

	if cmd, found := subCmds[args[1]]; found {
		cmd(l, args[2:])
	} else {
		l.Fatal(fmt.Errorf("Unknown subcommand: %s", args[1]))
	}
}

func find(l *logger.Logger, paths []string) {
	if rs, err := workflow.ReadersFromPaths(paths); err != nil {
		l.Fatal(err)
	} else if wf, err := workflow.WorkflowFromReaders(rs...); err != nil {
		l.Fatal(err)
	} else {
		for _, r := range wf.Repositories {
			err := r.Walk(func(f *workflow.File) error {
				if f == nil || !f.IsRegular {
					return nil
				}
				if filepath.Ext(f.Name) != ".an" {
					return nil
				}
				l.Log("element", filepath.Dir(f.Name))
				return nil
			})
			if err != nil {
				l.Fatal(err)
			}
		}
	}
}
