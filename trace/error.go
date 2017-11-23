package trace

import (
	"fmt"
)

// An Error is an error that arises during trace execution
type Error struct {
	BaseError interface{}
	Stack     []byte
}

func (a *Error) Error() string {
	return fmt.Sprintf("%s at:\n%s", a.BaseError, string(a.Stack))
}
