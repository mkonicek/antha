package workflow

import (
	"errors"
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
		if wfPath == "-" || wfPath == "" {
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
	if entries, err := ioutil.ReadDir(dir); err != nil && os.IsNotExist(err) {
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
