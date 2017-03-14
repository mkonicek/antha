package wtype

import (
	"fmt"
	"testing"
)

func TestColVectorPlateIterator(t *testing.T) {
	p := makeplatefortest()
	it := NewColVectorIterator(p, 8)
	c := 0
	for wv := it.Curr(); it.Valid(); wv = it.Next() {
		wv = wv
		c += 1
	}

	if c != 12 {
		t.Errorf(fmt.Sprintf("Expected 12 cols, got %d", c))
	}
}

func TestRowVectorPlateIterator(t *testing.T) {
	p := makeplatefortest()
	it := NewRowVectorIterator(p, 12)
	c := 0
	for wv := it.Curr(); it.Valid(); wv = it.Next() {
		wv = wv
		c += 1
	}

	if c != 8 {
		t.Errorf(fmt.Sprintf("Expected 8 rows, got %d", c))
	}
}
