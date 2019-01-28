package main

import (
	"flag"
	"fmt"
	"io"
	"log"
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

	workflows := flag.Args()
	if len(workflows) == 0 {
		log.Fatal("No workflow files provided (use - to read from stdin).")
	}

	stdinUnused := true
	rs := make([]io.Reader, len(workflows))
	for idx, wfPath := range workflows {
		if wfPath == "-" {
			if stdinUnused {
				stdinUnused = false
				rs[idx] = os.Stdin
			} else {
				log.Fatal("Workflow can only be read from stdin once")
			}

		} else {
			if fh, err := os.Open(wfPath); err != nil {
				log.Fatal(err)
			} else {
				defer fh.Close()
				rs[idx] = fh
			}
		}
	}

	wf, err := composer.WorkflowFromReaders(rs...)
	if err != nil {
		log.Fatal(err)
	}

	if comp, err := composer.NewComposer(outdir, wf); err != nil {
		log.Fatal(err)
	} else if err := comp.FindWorkflowElementTypes(); err != nil {
		log.Fatal(err)
	} else if err := comp.Transpile(); err != nil {
		log.Fatal(err)
	} else if err := comp.GenerateMain(); err != nil {
		log.Fatal(err)
	} else if err := comp.SaveWorkflow(); err != nil {
		log.Fatal(err)
	}
}
