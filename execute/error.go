package execute

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

// An Error reported by user code
type Error struct {
	message string
}

// Error returns the error message
func (a *Error) Error() string {
	return a.message
}

// Errorf reports an execution error. Does not return
func Errorf(ctx context.Context, format string, args ...interface{}) {
	userMsg := fmt.Sprintf(format, args...)
	elementName := getElementName(ctx)

	msg := userMsg
	if len(elementName) != 0 {
		msg = "element " + elementName + ": " + userMsg
	}

	var err error = &Error{message: msg}
	err = errors.WithStack(err)
	panic(err)
}

// unwrapError unpacks the result of Errorf
func unwrapError(obj interface{}) (error, bool) { // nolint
	err, ok := obj.(error)
	if !ok {
		return nil, false
	}

	if _, ok := errors.Cause(err).(*Error); !ok {
		return nil, false
	}

	return err, true
}
