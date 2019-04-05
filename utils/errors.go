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

// Pack returns a new ErrorSlice containing only the non-nill errors, or nil if there are none
// suitable for returning directly
func (es ErrorSlice) Pack() error {
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

type ErrorFunc func() error

type ErrorFuncs []ErrorFunc

// Run the ErrorFuncs in the supplied order until a non-nil error is
// encountered and return that error. Returns nil iff all funcs return
// nil errors.
func (efs ErrorFuncs) Run() error {
	for _, ef := range efs {
		if err := ef(); err != nil {
			return err
		}
	}
	return nil
}
