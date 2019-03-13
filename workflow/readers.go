package workflow

import (
	"errors"
	"io"
	"os"
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
