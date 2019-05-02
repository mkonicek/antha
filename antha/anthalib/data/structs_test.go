package data

import (
	"reflect"
	"testing"
)

func TestFromStructs(t *testing.T) {
	type foo struct {
		Bar        int64
		Baz        []string
		unexported float64 //nolint
	}
	tab := Must().NewTableFromStructs([]foo{})
	assertEqual(t, NewTable(
		makeNativeSeries("Bar", []int64{}, nil),
		makeNativeSeries("Baz", [][]string{}, nil),
	), tab, "empty")

	tab2 := Must().NewTableFromStructs([]struct{ A, B int64 }{{11, 2}, {33, 4}})
	assertEqual(t, NewTable(
		makeNativeSeries("A", []int64{11, 33}, nil),
		makeNativeSeries("B", []int64{2, 4}, nil),
	), tab2, "anonymous struct type, filled")

	s := &foo{Bar: 1}
	tab3 := Must().NewTableFromStructs([]*foo{s})
	assertEqual(t, NewTable(
		makeNativeSeries("Bar", []int64{1}, nil),
		makeNativeSeries("Baz", [][]string{nil}, nil),
	), tab3, "ptr struct type")

	tab4 := Must().NewTableFromStructs([]*foo{nil, {Bar: 1}}).
		Must().Project("Bar")
	assertEqual(t, NewTable(
		makeNativeSeries("Bar", []int64{0, 1}, nil),
	), tab4, "nil struct -> zeros")

	_, err := NewTableFromStructs(1)
	if err == nil {
		t.Error("kind check")
	}
	_, err = NewTableFromStructs([]int{1})
	if err == nil {
		t.Error("struct check")
	}
}
func TestToStructs(t *testing.T) {
	tab := NewTable(
		makeNativeSeries("A", []int64{1, 1000}, nil),
		makeNativeSeries("B", []string{"abcdef", "abcd"}, nil),
		makeNativeSeries("unexported", []string{"xx", "xx"}, nil),
		makeNativeSeries("Unmapped", []string{"xx", "xx"}, nil),
		makeNativeSeries("C", []float64{0, 0}, nil),
	)
	type destT struct {
		B          string
		A          int64
		unexported string //nolint
	}
	dest := []destT{}
	err := tab.ToStructs(&dest)
	if err != nil {
		t.Fatal(err)
	}
	// notice zeros set on the unexported field
	if !reflect.DeepEqual(dest, []destT{
		{A: 1, B: "abcdef"},
		{A: 1000, B: "abcd"},
	}) {
		t.Errorf("actual: %+v", dest)
	}
	roundtrip := Must().NewTableFromStructs(dest)
	// notice column order is set by struct field order
	expected := tab.Must().Project("B", "A")
	assertEqual(t, expected, roundtrip, "roundtrip")

	type destTWithUnbound struct {
		A       int64
		Unbound int // there is no such column, this will cause an error
	}
	destX := []destTWithUnbound{}
	err = tab.ToStructs(&destX)
	if err == nil {
		t.Error("unbound col check")
	}

	err = tab.Must().Project("B").Rename("B", "A").ToStructs(&dest)
	if err == nil {
		t.Error("assignability check")
	}

	err = tab.ToStructs(dest)
	if err == nil {
		t.Error("ptr check")
	}

	err = tab.ToStructs(1)
	if err == nil {
		t.Error("kind check")
	}

	destPtrs := []*destT{}
	err = tab.ToStructs(&destPtrs)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(destPtrs, []*destT{
		{A: 1, B: "abcdef"},
		{A: 1000, B: "abcd"},
	}) {
		t.Errorf("actual ptrs: %+v", destPtrs)
	}
}
func TestRow_ToStruct(t *testing.T) {
	tab := NewTable(
		makeNativeSeries("A", []int64{1, 1000}, nil),
		makeNativeSeries("ignored", []int64{2, 2}, nil),
	)
	for row := range tab.IterAll() {
		s := &struct{ A int64 }{}
		err := row.ToStruct(s)
		if err != nil {
			t.Fatal(err)
		}
		if s.A != row.ValueAt(0).MustInt64() {
			t.Errorf("at row %+v, %+v", row, s)
		}
	}
}
