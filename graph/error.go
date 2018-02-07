package graph

import (
	"errors"
)

var (
	// ErrTraversalDone is a predefined error for expected early termination of a traversal
	ErrTraversalDone = errors.New("traversal done")
	// ErrNextNode is a predefined error for continuing a traversal
	ErrNextNode = errors.New("next node")
)
