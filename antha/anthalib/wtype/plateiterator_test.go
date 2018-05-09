package wtype

import (
	"fmt"
	"testing"
)

func gettipboxfortest() *LHTipbox {
	shp := NewShape("cylinder", "mm", 7.3, 7.3, 51.2)
	w := NewLHWell("ul", 250.0, 10.0, shp, FlatWellBottom, 7.3, 7.3, 51.2, 0.0, "mm")
	tiptype := NewLHTip("me", "mytype", 0.5, 1000.0, "ul", false, shp)
	tb := NewLHTipbox(8, 12, Coordinates{127.76, 85.48, 120.0}, "me", "mytype", tiptype, w, 0.0, 0.0, 0.0, 0.0, 0.0)
	return tb
}

func TestColVectorTipboxIterator(t *testing.T) {
	tb := gettipboxfortest()
	it := NewColVectorIterator(tb, 8)
	c := 0
	for it.Curr(); it.Valid(); it.Next() {
		c += 1
	}

	if c != 12 {
		t.Errorf(fmt.Sprintf("Expected 12 cols, got %d", c))
	}
}

func TestRowVectorTipboxIterator(t *testing.T) {
	tb := gettipboxfortest()
	it := NewRowVectorIterator(tb, 12)
	c := 0
	for it.Curr(); it.Valid(); it.Next() {
		c += 1
	}

	if c != 8 {
		t.Errorf(fmt.Sprintf("Expected 8 rows, got %d", c))
	}
}

func TestColVectorPlateIterator(t *testing.T) {
	p := makeplatefortest()
	it := NewColVectorIterator(p, 8)
	c := 0
	for it.Curr(); it.Valid(); it.Next() {
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
	for it.Curr(); it.Valid(); it.Next() {
		c += 1
	}

	if c != 8 {
		t.Errorf(fmt.Sprintf("Expected 8 rows, got %d", c))
	}
}

func TestTickingPlateIterator(t *testing.T) {
	p := makeplatefortest()
	it := NewTickingColVectorIterator(p, 8, 1, 1)
	c := 0

	for it.Curr(); it.Valid(); it.Next() {
		c += 1
	}

	if c != 12 {
		t.Errorf(fmt.Sprintf("Expected 12 cols, got %d", c))
	}
}
func TestTickingPlateIterator2(t *testing.T) {
	p := make384platefortest()
	it := NewTickingColVectorIterator(p, 8, 1, 2)
	c := 0
	for it.Curr(); it.Valid(); it.Next() {
		c += 1
	}

	if c != 48 {
		t.Errorf(fmt.Sprintf("Expected 48 cols, got %d", c))
	}
}

func TestTickingPlateIterator3(t *testing.T) {
	p := make1536platefortest()
	it := NewTickingColVectorIterator(p, 8, 1, 4)
	c := 0
	for it.Curr(); it.Valid(); it.Next() {
		c += 1
	}

	if c != 192 {
		t.Errorf(fmt.Sprintf("Expected 192 cols, got %d", c))
	}
}

func TestTickingPlateIterator4(t *testing.T) {
	p := make24platefortest()

	it := NewTickingColVectorIterator(p, 8, 2, 1)
	c := 0
	for it.Curr(); it.Valid(); it.Next() {
		c += 1
	}
	if c != 6 {
		t.Errorf(fmt.Sprintf("Expected 6 cols, got %d ", c))
	}
}

func TestTickingPlateIterator5(t *testing.T) {
	p := make6platefortest()

	it := NewTickingColVectorIterator(p, 8, 4, 1)
	c := 0
	for it.Curr(); it.Valid(); it.Next() {
		c += 1
	}
	if c != 3 {
		t.Errorf(fmt.Sprintf("Expected 3 cols, got %d ", c))
	}
}

func TestTickingPlateIterator6(t *testing.T) {
	p := makeplatefortest()
	it := NewTickingColVectorIterator(p, 1, 1, 1)
	c := 0

	for it.Curr(); it.Valid(); it.Next() {
		c += 1
	}

	if c != 96 {
		t.Errorf(fmt.Sprintf("Expected 96 wells, got %d", c))
	}
}

func TestTickingPlateIterator7(t *testing.T) {
	t.Skip()
	p := makeplatefortest()
	it := NewTickingColVectorIterator(p, 2, 1, 1)
	c := 0

	for it.Curr(); it.Valid(); it.Next() {
		c += 1
	}

	if c != 48 {
		t.Errorf(fmt.Sprintf("Expected 48 wells, got %d", c))
	}
}
