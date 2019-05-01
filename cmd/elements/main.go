package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/antha-lang/antha/logger"
)

func main() {
	l := logger.NewLogger()

	subCmds := subCommands{
		"list":         list,
		"makeworkflow": makeWorkflowCmd,
		"describe":     describe,
		"defaults":     defaults,
	}

	args := os.Args
	if len(args) < 2 {
		logger.Fatal(l, fmt.Errorf("Subcommand needed. One of: %v", subCmds.List()))
	}

	if cmd, found := subCmds[strings.ToLower(args[1])]; found {
		if err := cmd(l, args[2:]); err != nil {
			logger.Fatal(l, err)
		}
	} else {
		logger.Fatal(l, fmt.Errorf("Unknown subcommand: %s. Available: %v", args[1], subCmds.List()))
	}
}

type subCommands map[string]func(*logger.Logger, []string) error

func (sc subCommands) List() []string {
	res := make([]string, 0, len(sc))
	for k := range sc {
		res = append(res, k)
	}
	return res
}
