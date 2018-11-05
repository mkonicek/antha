package execute

import (
	"context"
	"fmt"
)

// An UserError reported by user code
type UserError struct {
	message string
}

// Error satisfies the error interface
func (err UserError) Error() string {
	return err.message
}

// Errorf reports an execution error. Does not return
func Errorf(ctx context.Context, format string, args ...interface{}) {
	userMsg := fmt.Sprintf(format, args...)
	elementName := getElementName(ctx)

	msg := userMsg
	if len(elementName) != 0 {
		msg = "element " + elementName + ": " + userMsg
	}

	panic(UserError{message: msg})
}
