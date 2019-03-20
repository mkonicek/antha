package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/workflow"
)

func main() {
	l := logger.NewLogger()

	subCmds := subCommands{
		"list":         list,
		"makeWorkflow": makeWorkflow,
	}

	args := os.Args
	if len(args) < 2 {
		l.Fatal(fmt.Errorf("Subcommand needed. One of: %v", subCmds.List()))
	}

	if cmd, found := subCmds[args[1]]; found {
		if err := cmd(l, args[2:]); err != nil {
			l.Fatal(err)
		}
	} else {
		l.Fatal(fmt.Errorf("Unknown subcommand: %s. Available: %v", args[1], subCmds.List()))
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

func list(l *logger.Logger, args []string) error {
	flagSet := flag.NewFlagSet(flag.CommandLine.Name()+" list", flag.ContinueOnError)
	flagSet.Usage = workflow.NewFlagUsage(flagSet, "List all found element types, tab separated")

	inDir := ""
	flagSet.StringVar(&inDir, "indir", "", "Directory from which to read files (optional)")

	if err := flagSet.Parse(args); err != nil {
		return err
	} else if wfPaths, err := workflow.GatherPaths(flagSet, inDir); err != nil {
		return err
	} else {
		return findElements(l, wfPaths, func(r *workflow.Repository, et *workflow.ElementType) error {
			fmt.Printf("%v\t%v\t%v\n", et.Name(), et.ElementPath, et.RepositoryPrefix)
			return nil
		})
	}
}

func makeWorkflow(l *logger.Logger, args []string) error {
	flagSet := flag.NewFlagSet(flag.CommandLine.Name()+" makeWorkflow", flag.ContinueOnError)
	flagSet.Usage = workflow.NewFlagUsage(flagSet, "Modify workflow adding selected elements")

	var regexStr, inDir, toFile string
	flagSet.StringVar(&regexStr, "regex", "", "Regular expression to match against element type path (optional)")
	flagSet.StringVar(&inDir, "indir", "", "Directory from which to read files (optional)")
	flagSet.StringVar(&toFile, "to", "", "File to write to (default: will write to stdout)")

	if err := flagSet.Parse(args); err != nil {
		return err
	} else if wfPaths, err := workflow.GatherPaths(flagSet, inDir); err != nil {
		return err
	} else if regex, err := regexp.Compile(regexStr); err != nil {
		return err
	} else if rs, err := workflow.ReadersFromPaths(wfPaths); err != nil {
		return err
	} else if acc, err := workflow.WorkflowFromReaders(rs...); err != nil {
		return err
	} else {
		err := findElements(l, wfPaths, func(r *workflow.Repository, et *workflow.ElementType) error {
			if !regex.MatchString(string(et.ElementPath)) {
				return nil
			}
			etCopy := *et
			wf := &workflow.Workflow{
				SchemaVersion: workflow.CurrentSchemaVersion,
				Elements: workflow.Elements{
					Types: workflow.ElementTypes{&etCopy},
					Instances: workflow.ElementInstances{
						workflow.ElementInstanceName(etCopy.Name()): &workflow.ElementInstance{
							ElementTypeName: etCopy.Name(),
						},
					},
				},
			}
			return acc.Merge(wf)
		})
		if err != nil {
			return err
		}
		return acc.WriteToFile(toFile)
	}
}

func findElements(l *logger.Logger, paths []string, consumer func(*workflow.Repository, *workflow.ElementType) error) error {
	if rs, err := workflow.ReadersFromPaths(paths); err != nil {
		return err
	} else if wf, err := workflow.WorkflowFromReaders(rs...); err != nil {
		return err
	} else if repoToEts, err := wf.Repositories.FindAllElementTypes(); err != nil {
		return err
	} else {
		for repo, ets := range repoToEts {
			for _, et := range ets {
				if err := consumer(repo, &et); err != nil {
					return err
				}
			}
		}
		return nil
	}
}
