package main

import (
	"flag"
	"path/filepath"

	"github.com/antha-lang/antha/composer"
	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/workflow"
)

func main() {
	flag.Usage = workflow.NewFlagUsage(nil, "Parse, compile and run a workflow")

	var inDir, outDir string
	var keep, run, linkedDrivers bool
	flag.StringVar(&inDir, "indir", "", "Directory from which to read files (optional)")
	flag.StringVar(&outDir, "outdir", "", "Directory to write to (default: a temporary directory will be created)")
	flag.BoolVar(&keep, "keep", false, "Keep build environment if compilation is successful")
	flag.BoolVar(&run, "run", true, "Run the workflow if compilation is successful")
	flag.BoolVar(&linkedDrivers, "linkedDrivers", false, "Compile workflow with linked-in drivers")
	flag.Parse()

	logger := logger.NewLogger()

	if wfPaths, err := workflow.GatherPaths(nil, filepath.Join(inDir, "workflow")); err != nil {
		logger.Fatal(err)
	} else if rs, err := workflow.ReadersFromPaths(wfPaths); err != nil {
		logger.Fatal(err)
	} else if wf, err := workflow.WorkflowFromReaders(rs...); err != nil {
		logger.Fatal(err)
	} else if err := wf.Validate(); err != nil {
		logger.Fatal(err)
	} else if cb, err := composer.NewComposerBase(logger, inDir, outDir); err != nil {
		logger.Fatal(err)
	} else {
		err := cb.ComposeMainAndRun(wf, keep, run, linkedDrivers)
		defer cb.CloseLogs()
		if err == nil {
			logger.Log("progress", "complete")
		} else {
			logger.Fatal(err)
		}
	}
}
