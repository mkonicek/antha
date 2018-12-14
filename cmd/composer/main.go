package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/antha-lang/antha/composer"
)

func main() {
	var workflowPath, configPath string
	flag.StringVar(&workflowPath, "workflow", "-", "Path to workflow. (use '-' to read from stdin)")
	flag.StringVar(&configPath, "config", "-", "Path to config file.")
	flag.Parse()

	var cfg *composer.Config
	if configFH, err := os.Open(configPath); err != nil {
		log.Fatal(err)
	} else {
		defer configFH.Close()
		if cfg, err = composer.ConfigFromReader(configFH); err != nil {
			log.Fatal(err)
		}
	}

	var r io.Reader
	if workflowPath == "-" {
		r = os.Stdin
	} else {
		if fh, err := os.Open(workflowPath); err != nil {
			log.Fatal(err)
		} else {
			defer fh.Close()
			r = fh
		}
	}

	wf, err := composer.WorkflowFromReader(r)
	if err != nil {
		log.Fatal(err)
	}

	comp := composer.NewComposer(cfg, wf)
	if err := comp.LocateElementClasses(); err != nil {
		log.Fatal(err)
	} else if err := comp.Transpile(); err != nil {
		log.Fatal(err)
		//} else if err := comp.GenerateMain(os.Stdout); err != nil {
		//		log.Fatal(err)
	}
}
