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
	return strings.Join(strs, "; ")
}

func (es ErrorSlice) Nub() ErrorSlice {
	res := make(ErrorSlice, 0, len(es))
	for _, err := range es {
		if err != nil {
			res = append(res, err)
		}
	}
	if len(res) > 0 {
		return res
	} else {
		return nil
	}
}
