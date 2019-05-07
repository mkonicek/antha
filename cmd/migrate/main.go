package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/workflow"
	"github.com/antha-lang/antha/workflow/migrate"
	"github.com/antha-lang/antha/workflow/simulaterequestpb"
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

	args := flag.Args()
	rs, err := workflow.ReadersFromPaths(append(args, fromFile))
	if err != nil {
		logger.Fatal(l, err)
	}
	fromReader := rs[len(args)]
	rs = rs[:len(args)]

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

	var provider migrate.WorkflowProvider
	switch fromFormat {
	case fromFormatJSON:
		provider, err = v1_2.NewProvider(fromReader, fm, repoMap, gilsonDevice, l)
	case fromFormatProtobuf:
		provider, err = simulaterequestpb.NewProvider(fromReader, fm, repoMap, gilsonDevice, l)
	default:
		logger.Fatal(l, fmt.Errorf("Unknown format '%v', valid formats are: %v", fromFormat, validFromFormats))
	}
	if err != nil {
		logger.Fatal(l, err)
	}

	m := migrate.NewMigrator(provider)
	wf, err := m.Workflow()
	if err != nil {
		logger.Fatal(l, err)
	}

	// Merge the generated v2.0 workflow into the stuff we got from the snippets
	if err = cwf.Merge(wf); err != nil {
		logger.Fatal(l, err)
	}

	// validate the resulting workflow if the -validate flag was set
	if validate {
		if err = cwf.Validate(); err != nil {
			logger.Fatal(l, err)
		}
	}

	// Save to disk
	p := filepath.Join(outDir, "workflow", "workflow.json")
	if err = cwf.WriteToFile(p, true); err != nil {
		logger.Fatal(l, err)
	}
}
