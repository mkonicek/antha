package utils

import (
	"strings"
)

type ErrorSlice []error

func (es ErrorSlice) Error() string {
	strs := make([]string, len(es))
	for i, err := range es {
		strs[i] = err.Error()
	}
	return strings.Join(strs, "\n")
}
