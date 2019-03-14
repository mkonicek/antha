package main

import (
	"errors"
	"flag"
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
		"find":          find,
		"makeWorkflows": makeWorkflows,
		"makeWorkflow":  makeWorkflow,
	}

	if cmd, found := subCmds[args[1]]; found {
		cmd(l, args[2:])
	} else {
		l.Fatal(fmt.Errorf("Unknown subcommand: %s", args[1]))
	}
}

func find(l *logger.Logger, paths []string) {
	findElements(l, paths, func(r *workflow.Repository, et *workflow.ElementType) error {
		l.Log("repositoryPrefix", et.RepositoryPrefix, "elementPath", et.ElementPath)
		return nil
	})
}

func makeWorkflow(l *logger.Logger, args []string) {
	acc := workflow.EmptyWorkflow()
	acc.JobId = workflow.JobId("underTest")
	findElements(l, args, func(r *workflow.Repository, et *workflow.ElementType) error {
		etCopy := *et
		wf := &workflow.Workflow{
			SchemaVersion: workflow.CurrentSchemaVersion,
			Repositories: workflow.Repositories{
				et.RepositoryPrefix: r,
			},
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
	if err := acc.WriteToFile("/tmp/underTest.json"); err != nil {
		l.Fatal(err)
	}
}

func makeWorkflows(l *logger.Logger, args []string) {
	outdir := ""
	flagset := flag.NewFlagSet("makeWorkflows", flag.ContinueOnError)
	flagset.StringVar(&outdir, "outdir", "", "Directory to write to")
	if err := flagset.Parse(args); err != nil {
		l.Fatal(err)
	}
	paths := flagset.Args()
	if err := os.MkdirAll(outdir, 0700); err != nil {
		l.Fatal(err)
	}
	findElements(l, paths, func(r *workflow.Repository, et *workflow.ElementType) error {
		etCopy := *et
		wf := &workflow.Workflow{
			SchemaVersion: workflow.CurrentSchemaVersion,
			JobId:         workflow.JobId("underTest"),
			Repositories: workflow.Repositories{
				et.RepositoryPrefix: r,
			},
			Elements: workflow.Elements{
				Types: workflow.ElementTypes{&etCopy},
				Instances: workflow.ElementInstances{
					workflow.ElementInstanceName(etCopy.Name()): &workflow.ElementInstance{
						ElementTypeName: etCopy.Name(),
					},
				},
			},
		}
		dir := filepath.Join(outdir, filepath.FromSlash(string(et.ElementPath)))
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		} else {
			return wf.WriteToFile(filepath.Join(dir, "underTest.json"))
		}
	})
}

func findElements(l *logger.Logger, paths []string, consumer func(*workflow.Repository, *workflow.ElementType) error) {
	if rs, err := workflow.ReadersFromPaths(paths); err != nil {
		l.Fatal(err)
	} else if wf, err := workflow.WorkflowFromReaders(rs...); err != nil {
		l.Fatal(err)
	} else {
		for rp, r := range wf.Repositories {
			var et workflow.ElementType
			err := r.Walk(func(f *workflow.File) error {
				if f == nil || !f.IsRegular {
					return nil
				}
				if filepath.Ext(f.Name) != ".an" {
					return nil
				}
				et.RepositoryPrefix = rp
				et.ElementPath = workflow.ElementPath(filepath.ToSlash(filepath.Dir(f.Name)))
				return consumer(r, &et)
			})
			if err != nil {
				l.Fatal(err)
			}
		}
	}
}
