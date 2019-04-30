package main

import (
	"flag"

	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/workflow"
	"github.com/antha-lang/antha/workflow/v1_2"
)

func main() {
	flag.Usage = workflow.NewFlagUsage(nil, "Migrate workflow to latest schema version")

	var fromFile, outDir, gilsonDevice string
	var validate bool
	flag.StringVar(&outDir, "outdir", "", "Directory to write to (default: a temporary directory will be created)")
	flag.StringVar(&fromFile, "from", "", "File to migrate from (default: will be read from stdin)")
	flag.StringVar(&gilsonDevice, "gilson-device", "", "A gilson device name to use for migrated config. If not present, device specific configuration will not be migrated.")
	flag.BoolVar(&validate, "validate", true, "Validate input and output files.")
	flag.Parse()

	l := logger.NewLogger()

	m, err := v1_2.NewMigrater(l, flag.Args(), fromFile, outDir, gilsonDevice)
	if err != nil {
		logger.Fatal(l, err)
	}

	if err := m.ValidateOld(); err != nil {
		if validate {
			logger.Fatal(l, err)
		} else {
			l.Log("OriginalFileValidationError", err)
		}
	}

	if err := m.MigrateAll(); err != nil {
		logger.Fatal(l, err)
	}

	if err := m.ValidateCur(); err != nil {
		if validate {
			logger.Fatal(l, err)
		} else {
			l.Log("ValidationError", err)
		}
	}

	if err := m.SaveCur(); err != nil {
		logger.Fatal(l, err)
	}
}
