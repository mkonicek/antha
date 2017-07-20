package execute

import (
	"context"
	"fmt"
)

// An Error reported by user code
type Error struct {
	Message string
}

// Error returns the error message
func (a *Error) Error() string {
	return a.Message
}

// Errorf reports an execution error. Does not return
func Errorf(ctx context.Context, format string, args ...interface{}) {
	panic(&Error{Message: fmt.Sprintf(format, args...)})
}
