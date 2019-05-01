package workflow

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ReadersFromPaths(paths []string) ([]io.ReadCloser, error) {
	if len(paths) == 0 {
		return nil, nil
	}

	stdinUsed := false
	rs := make([]io.ReadCloser, len(paths))
	for idx, wfPath := range paths {
		if wfPath == "-" {
			if stdinUsed {
				return nil, errors.New("Stdin may only be used once")
			} else {
				stdinUsed = true
				rs[idx] = os.Stdin
			}

		} else {
			if fh, err := os.Open(wfPath); err != nil {
				return nil, err
			} else {
				rs[idx] = fh
			}
		}
	}

	return rs, nil
}

func JsonPathsWithin(dir string) ([]string, error) {
	if dir == "" {
		return nil, nil
	} else if entries, err := ioutil.ReadDir(dir); err != nil && os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		results := make([]string, 0, len(entries))
		for _, entry := range entries {
			if entry.Mode().IsRegular() && filepath.Ext(entry.Name()) == ".json" {
				results = append(results, filepath.Join(dir, entry.Name()))
			}
		}
		return results, nil
	}
}

func NewFlagUsage(fs *flag.FlagSet, summary string) func() {
	if fs == nil {
		fs = flag.CommandLine
	}
	output := fs.Output()
	name := fs.Name()
	return func() {
		fmt.Fprintf(output, "%s: %s\nUsage of %s:\n", name, summary, name)
		fs.PrintDefaults()
		fmt.Fprintf(output, "All further args are interpreted as paths to workflows to be merged and composed. Use - to read a workflow from stdin.\n")
	}
}

func GatherPaths(fs *flag.FlagSet, inDir string) ([]string, error) {
	if fs == nil {
		fs = flag.CommandLine
	}
	wfPaths := fs.Args()
	if moreWfPaths, err := JsonPathsWithin(inDir); err != nil {
		return nil, err
	} else {
		return append(wfPaths, moreWfPaths...), nil
	}
}
