package meta

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

var defaultTestOpt = UnmarshalOpt{
	Struct: func(bs []byte, obj interface{}) error {
		return json.Unmarshal(bs, obj)
	},
}

func TestString(t *testing.T) {
	type Value string
	var x Value
	golden := Value("hello")

	if err := UnmarshalJSON(defaultTestOpt, []byte(`"hello"`), &x); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(golden, x) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestInt(t *testing.T) {
	type Value int
	var x Value
	golden := Value(1)
	if err := UnmarshalJSON(defaultTestOpt, []byte(`1`), &x); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(golden, x) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestStruct(t *testing.T) {
	type Value struct {
		A string
		B int
	}
	var x Value
	golden := Value{A: "hello", B: 1}

	if err := UnmarshalJSON(defaultTestOpt, []byte(`{"A": "hello", "B": 1}`), &x); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(golden, x) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestPtrStruct(t *testing.T) {
	type Value struct {
		A string
		B int
	}
	var x *Value
	golden := &Value{A: "hello", B: 1}

	if err := UnmarshalJSON(defaultTestOpt, []byte(`{"A": "hello", "B": 1}`), &x); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(golden, x) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestMap(t *testing.T) {
	type Elem struct {
		A string
		B int
	}
	type Value map[string]Elem
	var x Value
	golden := Value{
		"A": Elem{A: "hello", B: 1},
		"B": Elem{A: "hello", B: 2},
	}
	if err := UnmarshalJSON(defaultTestOpt, []byte(`{"A": {"A": "hello", "B": 1}, "B": {"A": "hello", "B": 2} }`), &x); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(golden, x) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestSlice(t *testing.T) {
	type Elem struct {
		A string
		B int
	}
	type Value []Elem
	var x Value
	golden := Value{
		Elem{A: "hello", B: 1},
		Elem{A: "hello", B: 2},
	}
	if err := UnmarshalJSON(defaultTestOpt, []byte(`[ {"A": "hello", "B": 1}, {"A": "hello", "B": 2} ]`), &x); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(golden, x) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestMapInterface(t *testing.T) {
	type Value map[string]interface{}
	someTime, err := time.Parse(time.Kitchen, time.Kitchen)
	if err != nil {
		t.Fatal(err)
	}
	x := Value{
		"A": time.Time{},
	}
	golden := Value{
		"A": someTime,
	}
	bs, err := json.Marshal(someTime)
	if err != nil {
		t.Fatal(err)
	}

	if err := UnmarshalJSON(defaultTestOpt, []byte(`{ "A": `+string(bs)+` }`), &x); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(golden, x) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestSliceInterface(t *testing.T) {
	type Value []interface{}
	someTime, err := time.Parse(time.Kitchen, time.Kitchen)
	if err != nil {
		t.Fatal(err)
	}
	x := Value{time.Time{}}
	golden := Value{someTime}
	bs, err := json.Marshal(someTime)
	if err != nil {
		t.Fatal(err)
	}

	if err := UnmarshalJSON(defaultTestOpt, []byte(`[`+string(bs)+` ]`), &x); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(golden, x) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestPtrInterface(t *testing.T) {
	someTime, err := time.Parse(time.Kitchen, time.Kitchen)
	if err != nil {
		t.Fatal(err)
	}
	var x interface{} = time.Time{}
	ptr := &x
	var golden interface{} = someTime
	goldenPtr := &golden
	bs, err := json.Marshal(someTime)
	if err != nil {
		t.Fatal(err)
	}

	if err := UnmarshalJSON(defaultTestOpt, bs, &ptr); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(goldenPtr, ptr) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestInterface(t *testing.T) {
	someTime, err := time.Parse(time.Kitchen, time.Kitchen)
	if err != nil {
		t.Fatal(err)
	}
	var x interface{} = time.Time{}
	var golden interface{} = someTime
	bs, err := json.Marshal(someTime)
	if err != nil {
		t.Fatal(err)
	}

	if err := UnmarshalJSON(defaultTestOpt, bs, &x); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(golden, x) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}
