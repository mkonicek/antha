package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"

	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/workflow"
)

func defaults(l *logger.Logger, args []string) error {
	flagSet := flag.NewFlagSet(flag.CommandLine.Name()+" defaults", flag.ContinueOnError)
	flagSet.Usage = workflow.NewFlagUsage(flagSet, "Gather defaults for an element set from metadata.json files in the repo")

	var regexStr, inDir, outputFormat string
	flagSet.StringVar(&regexStr, "regex", "", "Regular expression to match against element type path (optional)")
	flagSet.StringVar(&inDir, "indir", "", "Directory from which to read files (optional)")
	flagSet.StringVar(&outputFormat, "format", "human", "Format to output data in. One of [human, json, protobuf]")

	if err := flagSet.Parse(args); err != nil {
		return err
	} else if wfPaths, err := workflow.GatherPaths(flagSet, inDir); err != nil {
		return err
	} else if rs, err := workflow.ReadersFromPaths(wfPaths); err != nil {
		return err
	} else if wf, err := workflow.WorkflowFromReaders(rs...); err != nil {
		return err
	} else if regex, err := regexp.Compile(regexStr); err != nil {
		return err
	} else {
		// the map keys are the dir paths of the element so that it's the same for the antha file and the metadata
		defaults := make(map[string]interface{})

		for _, repo := range wf.Repositories {
			err := repo.Walk(func(f *workflow.File) error {
				dir := filepath.Dir(f.Name)
				if (!workflow.IsAnthaMetadata(f.Name)) || !regex.MatchString(dir) {
					return nil
				}

				bs, err := ioutil.ReadFile(path.Join(repo.Directory, f.Name))
				if err != nil {
					return err
				}

				var doc map[string]interface{}
				err = json.Unmarshal(bs, &doc)
				if err != nil {
					return err
				}

				name, ok := doc["name"].(string)
				if !ok {
					return fmt.Errorf("Got unexpected data in name field: expected string, got %v", reflect.TypeOf(doc["name"]))
				}
				defaults[name] = doc["defaults"]

				return nil
			})
			if err != nil {
				return err
			}
		}

		bs, err := json.Marshal(defaults)
		if err != nil {
			return err
		}
		w := bufio.NewWriter(os.Stdout)
		defer w.Flush()
		_, err = w.Write(bs)
		if err != nil {
			return err
		}

		return nil
	}
}
