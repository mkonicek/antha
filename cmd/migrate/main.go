package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/workflow/v1_2"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "All further args are interpreted as paths to up to date workflows to be merged and used as a basis for the upgrade.\n")
	}

	var fromFile, toFile, gilsonDevice string
	var validate bool
	flag.StringVar(&toFile, "to", "", "File to write to (default: will write to stdout)")
	flag.StringVar(&fromFile, "from", "", "File to migrate from (default: will be read from stdin)")
	flag.StringVar(&gilsonDevice, "gilson-device", "", "A gilson device name to use for migrated config. If not present, device specific configuration will not be migrated.")
	flag.BoolVar(&validate, "validate", true, "Validate input and output files.")
	flag.Parse()

	logger := logger.NewLogger()

	m, err := v1_2.NewMigrater(logger, flag.Args(), fromFile, gilsonDevice)
	if err != nil {
		logger.Fatal(err)
	}

	if err := m.ValidateOld(); err != nil {
		if validate {
			logger.Fatal(err)
		} else {
			logger.Log("OriginalFileValidationError", err)
		}
	}

	if err := m.MigrateAll(); err != nil {
		logger.Fatal(err)
	}

	if err := m.ValidateCur(); err != nil {
		if validate {
			logger.Fatal(err)
		} else {
			logger.Log("ValidationError", err)
		}
	}

	if err = m.WriteToPath(toFile); err != nil {
		logger.Fatal(err)
	}
}
