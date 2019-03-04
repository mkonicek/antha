package workflow

import (
	"errors"
	"io"
	"os"
)

func ReadersFromPaths(paths []string) ([]io.ReadCloser, error) {
	if len(paths) == 0 {
		return nil, errors.New("No workflow files provided (use - to read from stdin).")
	}

	stdinUnused := true
	rs := make([]io.ReadCloser, len(paths))
	for idx, wfPath := range paths {
		if wfPath == "-" {
			if stdinUnused {
				stdinUnused = false
				rs[idx] = os.Stdin
			} else {
				return nil, errors.New("Workflow can only be read from stdin once")
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
