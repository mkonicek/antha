package wtype

import (
	"testing"
)

func TestComponentVectorEqual(t *testing.T) {
	cv := ComponentVector{
		&Liquid{},
		&Liquid{},
		nil,
		&Liquid{},
	}

	if !cv.Equal(cv) {
		t.Errorf("Vector must be equal to itself")
	}
}
