package meta

import (
	"bytes"
	"testing"
	"time"
)

const (
	kitchenZero = "0000-01-01T15:04:00Z"
)

func TestMarshalString(t *testing.T) {
	x := "hello"
	golden := []byte(`"hello"`)

	var m Marshaler

	if bs, err := m.Marshal(x); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(golden, bs) {
		t.Errorf("expecting %q but got %q instead", golden, bs)
	}
}

func TestMarshalInt(t *testing.T) {
	x := 1
	golden := []byte(`1`)

	var m Marshaler

	if bs, err := m.Marshal(x); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(golden, bs) {
		t.Errorf("expecting %q but got %q instead", golden, bs)
	}
}

func TestMarshalStruct(t *testing.T) {
	type Value struct {
		A string
		B int
	}
	x := Value{
		A: "hello",
		B: 1,
	}
	golden := []byte(`{"A":"hello","B":1}`)

	var m Marshaler

	if bs, err := m.Marshal(x); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(golden, bs) {
		t.Errorf("expecting %q but got %q instead", golden, bs)
	}
}

func TestMarshalPtrStruct(t *testing.T) {
	type Value struct {
		A string
		B int
	}
	x := &Value{
		A: "hello",
		B: 1,
	}
	golden := []byte(`{"A":"hello","B":1}`)

	var m Marshaler

	if bs, err := m.Marshal(x); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(golden, bs) {
		t.Errorf("expecting %q but got %q instead", golden, bs)
	}
}

func TestMarshalMap(t *testing.T) {
	type Elem struct {
		A string
		B int
	}
	type Value map[string]Elem
	x := Value{
		"A": Elem{
			A: "hello",
			B: 1,
		},
		"B": Elem{
			A: "hello",
			B: 2,
		},
	}
	golden := []byte(`{"A":{"A":"hello","B":1},"B":{"A":"hello","B":2}}`)

	var m Marshaler

	if bs, err := m.Marshal(x); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(golden, bs) {
		t.Errorf("expecting %q but got %q instead", golden, bs)
	}
}

func TestMarshalSlice(t *testing.T) {
	type Elem struct {
		A string
		B int
	}
	x := []Elem{
		{A: "hello", B: 1},
		{A: "hello", B: 2},
	}
	golden := []byte(`[{"A":"hello","B":1},{"A":"hello","B":2}]`)

	var m Marshaler

	if bs, err := m.Marshal(x); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(golden, bs) {
		t.Errorf("expecting %q but got %q instead", golden, bs)
	}
}

func TestMarshalMapInterface(t *testing.T) {
	type Value map[string]interface{}
	someTime, err := time.Parse(time.Kitchen, time.Kitchen)
	if err != nil {
		t.Fatal(err)
	}
	x := Value{
		"A": someTime,
	}
	golden := []byte(`{"A":"` + kitchenZero + `"}`)

	var m Marshaler

	if bs, err := m.Marshal(x); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(golden, bs) {
		t.Errorf("expecting %q but got %q instead", golden, bs)
	}
}

func TestMarshalSliceInterface(t *testing.T) {
	type Value []interface{}
	someTime, err := time.Parse(time.Kitchen, time.Kitchen)
	if err != nil {
		t.Fatal(err)
	}
	x := Value{someTime}
	golden := []byte(`["` + kitchenZero + `"]`)

	var m Marshaler

	if bs, err := m.Marshal(x); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(golden, bs) {
		t.Errorf("expecting %q but got %q instead", golden, bs)
	}
}

func TestMarshalInterface(t *testing.T) {
	type Value interface{}
	someTime, err := time.Parse(time.Kitchen, time.Kitchen)
	if err != nil {
		t.Fatal(err)
	}
	x := Value(someTime)
	golden := []byte(`"` + kitchenZero + `"`)

	var m Marshaler

	if bs, err := m.Marshal(x); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(golden, bs) {
		t.Errorf("expecting %q but got %q instead", golden, bs)
	}
}

func TestMarshalPtrInterface(t *testing.T) {
	type Value interface{}
	someTime, err := time.Parse(time.Kitchen, time.Kitchen)
	if err != nil {
		t.Fatal(err)
	}
	elem := Value(someTime)
	x := &elem
	golden := []byte(`"` + kitchenZero + `"`)

	var m Marshaler

	if bs, err := m.Marshal(x); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(golden, bs) {
		t.Errorf("expecting %q but got %q instead", golden, bs)
	}
}
