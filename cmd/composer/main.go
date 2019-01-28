package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/antha-lang/antha/composer"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "All further args are interpreted as paths to workflows to be merged and composed. Use - to read a workflow from stdin.\n")
	}

	var outdir string
	flag.StringVar(&outdir, "outdir", "", "Directory to write to (default: a temporary directory will be created)")
	flag.Parse()

	logger := composer.NewLogger()

	workflows := flag.Args()
	if len(workflows) == 0 {
		logger.Fatal(errors.New("No workflow files provided (use - to read from stdin)."))
	}

	stdinUnused := true
	rs := make([]io.Reader, len(workflows))
	for idx, wfPath := range workflows {
		if wfPath == "-" {
			if stdinUnused {
				stdinUnused = false
				rs[idx] = os.Stdin
			} else {
				logger.Fatal(errors.New("Workflow can only be read from stdin once"))
			}

		} else {
			if fh, err := os.Open(wfPath); err != nil {
				logger.Fatal(err)
			} else {
				defer fh.Close()
				rs[idx] = fh
			}
		}
	}

	if wf, err := composer.WorkflowFromReaders(rs...); err != nil {
		logger.Fatal(err)
	} else if comp, err := composer.NewComposer(logger, outdir, wf); err != nil {
		logger.Fatal(err)
	} else if err := comp.FindWorkflowElementTypes(); err != nil {
		logger.Fatal(err)
	} else if err := comp.Transpile(); err != nil {
		logger.Fatal(err)
	} else if err := comp.GenerateMain(); err != nil {
		logger.Fatal(err)
	} else if err := comp.SaveWorkflow(); err != nil {
		logger.Fatal(err)
	} else if err := comp.CompileWorkflow(); err != nil {
		logger.Fatal(err)
	}
}
