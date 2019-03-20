package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/antha-lang/antha/composer"
	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/workflow"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "All further args are interpreted as paths to workflows to be merged and composed. Use - to read a workflow from stdin.\n")
	}

	var inDir, outDir string
	var keep, run, linkedDrivers bool
	flag.StringVar(&inDir, "indir", "", "Directory from which to read files")
	flag.StringVar(&outDir, "outdir", "", "Directory to write to (default: a temporary directory will be created)")
	flag.BoolVar(&keep, "keep", false, "Keep build environment if compilation is successful")
	flag.BoolVar(&run, "run", true, "Run the workflow if compilation is successful")
	flag.BoolVar(&linkedDrivers, "linkedDrivers", false, "Compile workflow with linked-in drivers")
	flag.Parse()

	logger := logger.NewLogger()

	wfPaths := flag.Args()
	if inDir != "" {
		if moreWfPaths, err := workflow.JsonPathsWithin(filepath.Join(inDir, "workflow")); err != nil {
			logger.Fatal(err)
		} else {
			wfPaths = append(wfPaths, moreWfPaths...)
		}
	}

	if rs, err := workflow.ReadersFromPaths(wfPaths); err != nil {
		logger.Fatal(err)
	} else if wf, err := workflow.WorkflowFromReaders(rs...); err != nil {
		logger.Fatal(err)
	} else if err := wf.Validate(); err != nil {
		logger.Fatal(err)
	} else if comp, err := composer.NewComposer(logger, wf, inDir, outDir, keep, run, linkedDrivers); err != nil {
		logger.Fatal(err)
	} else {
		defer comp.CloseLogs()
		if err := comp.ComposeAndRun(); err != nil {
			logger.Fatal(err)
		} else {
			logger.Log("progress", "complete")
		}
	}
}
