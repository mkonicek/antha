package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/antha-lang/antha/composer"
	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/workflow"
)

func main() {
	l := logger.NewLogger()

	subCmds := subCommands{
		"list":         list,
		"makeWorkflow": makeWorkflow,
		"describe":     describe,
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

	var regexStr, inDir string
	flagSet.StringVar(&regexStr, "regex", "", "Regular expression to match against element type path (optional)")
	flagSet.StringVar(&inDir, "indir", "", "Directory from which to read files (optional)")

	if err := flagSet.Parse(args); err != nil {
		return err
	} else if wfPaths, err := workflow.GatherPaths(flagSet, inDir); err != nil {
		return err
	} else if regex, err := regexp.Compile(regexStr); err != nil {
		return err
	} else {
		return findElements(l, wfPaths, func(r *workflow.Repository, et *workflow.ElementType) error {
			if !regex.MatchString(string(et.ElementPath)) {
				return nil
			}
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
	} else if rs, err := workflow.ReadersFromPaths(wfPaths); err != nil {
		return err
	} else if acc, err := workflow.WorkflowFromReaders(rs...); err != nil {
		return err
	} else if regex, err := regexp.Compile(regexStr); err != nil {
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
		return acc.WriteToFile(toFile, true)
	}
}

func describe(l *logger.Logger, args []string) error {
	flagSet := flag.NewFlagSet(flag.CommandLine.Name()+" list", flag.ContinueOnError)
	flagSet.Usage = workflow.NewFlagUsage(flagSet, "Show descriptions elements")

	var regexStr, inDir string
	flagSet.StringVar(&regexStr, "regex", "", "Regular expression to match against element type path (optional)")
	flagSet.StringVar(&inDir, "indir", "", "Directory from which to read files (optional)")

	if err := flagSet.Parse(args); err != nil {
		return err
	} else if wfPaths, err := workflow.GatherPaths(flagSet, inDir); err != nil {
		return err
	} else if rs, err := workflow.ReadersFromPaths(wfPaths); err != nil {
		return err
	} else if wf, err := workflow.WorkflowFromReaders(rs...); err != nil {
		return err
	} else if regex, err := regexp.Compile(regexStr); err != nil {
		return err
	} else {
		for prefix, repo := range wf.Repositories {
			err := repo.Walk(func(f *workflow.File) error {
				if !workflow.IsAnthaFile(f.Name) {
					return nil
				}
				et := &workflow.ElementType{
					ElementPath:      workflow.ElementPath(filepath.Dir(f.Name)),
					RepositoryPrefix: prefix,
				}
				if !regex.MatchString(string(et.ElementPath)) {
					return nil
				}
				tet := composer.NewTranspilableElementType(et)
				if reader, err := f.Contents(); err != nil {
					return err
				} else {
					defer reader.Close()
					if bs, err := ioutil.ReadAll(reader); err != nil {
						return err
					} else if meta, err := tet.Meta(bs, f.Name); err != nil {
						return err
					} else {
						fmt.Printf("%v\n RepositoryPrefix: %v\n ElementPath: %v\n Description:\n  %s\n",
							et.Name(), et.RepositoryPrefix, et.ElementPath, strings.Replace(meta.Description, "\n", "\n  ", -1))
						return nil
					}
				}
			})
			if err != nil {
				return err
			}
		}
		return nil
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
