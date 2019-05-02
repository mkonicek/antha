package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/workflow"
	"github.com/antha-lang/antha/workflow/migrate"
	"github.com/antha-lang/antha/workflow/migrate/provider"
	"github.com/antha-lang/antha/workflow/v1_2"
)

const (
	fromFormatJSON     = "json"
	fromFormatProtobuf = "protobuf"
)

var validFromFormats = []string{fromFormatJSON, fromFormatProtobuf}

func main() {
	flag.Usage = workflow.NewFlagUsage(nil, "Migrate workflow to latest schema version")

	var fromFile, outDir, gilsonDevice, fromFormat string
	var validate bool
	flag.StringVar(&outDir, "outdir", "", "Directory to write to (default: a temporary directory will be created)")
	flag.StringVar(&fromFile, "from", "-", "File to migrate from (default: will be read from stdin)")
	flag.StringVar(&fromFormat, "format", fromFormatJSON, fmt.Sprintf("Format of the from file. One of: %v", validFromFormats))
	flag.StringVar(&gilsonDevice, "gilson-device", "", "A gilson device name to use for migrated config. If not present, device specific configuration will not be migrated.")
	flag.BoolVar(&validate, "validate", true, "Validate input and output files.")
	flag.Parse()

	l := logger.NewLogger()

	// Read in the workflow snippets
	rs, err := workflow.ReadersFromPaths(flag.Args())
	if err != nil {
		logger.Fatal(l, err)
	}

	cwf, err := workflow.WorkflowFromReaders(rs...)
	if err != nil {
		logger.Fatal(l, err)
	}

	// Get the repo map
	repoMap, err := cwf.Repositories.FindAllElementTypes()
	if err != nil {
		logger.Fatal(l, err)
	}

	// Prepare the output directory
	if outDir == "" {
		if outDir, err = ioutil.TempDir("", "antha-migrater"); err != nil {
			logger.Fatal(l, err)
		}
	}
	for _, leaf := range []string{"workflow", "data"} {
		if err := os.MkdirAll(filepath.Join(outDir, leaf), 0700); err != nil {
			logger.Fatal(l, err)
		}
	}

	dataDir := filepath.Join(outDir, "data")
	fm, err := effects.NewFileManager(dataDir, dataDir)
	if err != nil {
		logger.Fatal(l, err)
	}

	// Sanity check: we can only read from STDIN once
	inputPaths := append(flag.Args(), fromFile)
	stdinCount := 0
	for _, path := range inputPaths {
		if path == "-" {
			stdinCount++
		}
	}
	if stdinCount > 1 {
		logger.Fatal(l, errors.New("Input '-' specified more than once: can only read from STDIN once"))
	}

	var fromReader io.ReadCloser
	if fromFile == "-" {
		fromReader = os.Stdin
	} else {
		if fromReader, err = os.Open(fromFile); err != nil {
			logger.Fatal(l, err)
		}
	}
	defer fromReader.Close()

	var provider provider.WorkflowProvider
	switch fromFormat {
	case fromFormatJSON:
		{
			provider, err = v1_2.NewProvider(fromReader, fm, repoMap, gilsonDevice, l)
			if err != nil {
				logger.Fatal(l, err)
			}
		}
	case fromFormatProtobuf:
		{
			logger.Fatal(l, fmt.Errorf("Format not implemented"))
		}
	default:
		{
			logger.Fatal(l, fmt.Errorf("Unknown format '%v', valid formats are: %v", fromFormat, validFromFormats))
		}
	}

	m := migrate.NewMigrator(l, provider)
	wf, err := m.Workflow()
	if err != nil {
		logger.Fatal(l, err)
	}

	// Merge the generated v2.0 workflow into the stuff we got from the snippets
	if err = cwf.Merge(wf); err != nil {
		logger.Fatal(l, err)
	}

	// validate the resulting workflow
	if err = cwf.Validate(); err != nil {
		logger.Fatal(l, err)
	}

	// Save to disk
	p := filepath.Join(outDir, "workflow", "workflow.json")
	if err = cwf.WriteToFile(p, true); err != nil {
		logger.Fatal(l, err)
	}
}
