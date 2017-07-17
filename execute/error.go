package execute

import (
	"context"
	"fmt"
)

// An Error returned during execution
type Error struct {
	Message string
}

func (a *Error) Error() string {
	return a.Message
}

// Errorf reports an execution error. Does not return
func Errorf(ctx context.Context, format string, args ...interface{}) {
	panic(&Error{Message: fmt.Sprintf(format, args...)})
}
