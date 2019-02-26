package execute

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
	"github.com/antha-lang/antha/meta"
)

func unmarshal(ctx context.Context, obj interface{}, data []byte) error {
	var u unmarshaler

	um := &meta.Unmarshaler{
		Struct: func(bs []byte, obj interface{}) error {
			return u.unmarshalStruct(ctx, bs, obj)
		},
	}
	return um.Unmarshal(data, obj)
}

func TestString(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	type Value string
	var x Value
	golden := Value("hello")

	if err := unmarshal(ctx, &x, []byte(`"hello"`)); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(x, golden) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestInt(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	type Value int
	var x Value
	golden := Value(1)
	if err := unmarshal(ctx, &x, []byte(`1`)); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(x, golden) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestStruct(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	type Value struct {
		A string
		B int
	}
	var x Value
	golden := Value{A: "hello", B: 1}

	if err := unmarshal(ctx, &x, []byte(`{"A": "hello", "B": 1}`)); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(x, golden) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestMap(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

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
	if err := unmarshal(ctx, &x, []byte(`{"A": {"A": "hello", "B": 1}, "B": {"A": "hello", "B": 2} }`)); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(x, golden) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestSlice(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

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
	if err := unmarshal(ctx, &x, []byte(`[ {"A": "hello", "B": 1}, {"A": "hello", "B": 2} ]`)); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(x, golden) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestConstruct(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	var x wtype.LHTipbox
	if err := unmarshal(ctx, &x, []byte(`"CyBio250Tipbox"`)); err != nil {
		t.Fatal(err)
	}
}

func TestConstructMapFailure(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	type Elem struct {
		A string
		T *wtype.LHTipbox
	}
	type Value map[string]Elem
	var x Value
	if err := unmarshal(ctx, &x, []byte(`{"A": {"A": "hello", "T": "CyBio250Tipbox"} }`)); err == nil {
		t.Fatal("expecting failure but got success")
	}
}

func TestConstructMap(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	type Value map[string]interface{}
	x := Value{
		"A": &wtype.LHTipbox{},
		"B": 0,
		"C": "",
	}
	tb, err := inventory.NewTipbox(ctx, "CyBio250Tipbox")
	if err != nil {
		t.Fatal(err)
	}
	golden := Value{
		"A": tb,
		"B": 1,
		"C": "hello",
	}
	if err := unmarshal(ctx, &x, []byte(`{"A": "CyBio250Tipbox", "B": 1, "C": "hello" }`)); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(x["B"], golden["B"]) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	} else if !reflect.DeepEqual(x["C"], golden["C"]) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	} else if aa, ok := golden["A"].(*wtype.LHTipbox); !ok {
		t.Errorf("expecting %v but got %v instead", golden, x)
	} else if bb, ok := x["A"].(*wtype.LHTipbox); !ok {
		t.Errorf("expecting %v but got %v instead", golden, x)
	} else if aa.Type != bb.Type {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestConstructSlice(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	type Value []interface{}
	x := Value{
		&wtype.LHTipbox{},
		&wtype.Plate{},
	}
	tb1, err := inventory.NewTipbox(ctx, "CyBio250Tipbox")
	if err != nil {
		t.Fatal(err)
	}
	tb2, err := inventory.NewPlate(ctx, "pcrplate_with_cooler")
	if err != nil {
		t.Fatal(err)
	}
	golden := Value{
		tb1,
		tb2,
	}
	if err := unmarshal(ctx, &x, []byte(`[ "CyBio250Tipbox", "pcrplate_with_cooler" ]`)); err != nil {
		t.Fatal(err)
	} else if len(x) != 2 {
		t.Errorf("expecting %v but got %v instead", golden, x)
	} else if aa, ok := golden[0].(*wtype.LHTipbox); !ok {
		t.Errorf("expecting %v but got %v instead", golden, x)
	} else if bb, ok := x[0].(*wtype.LHTipbox); !ok {
		t.Errorf("expecting %v but got %v instead", golden, x)
	} else if aa.Type != bb.Type {
		t.Errorf("expecting %v but got %v instead", golden, x)
	} else if aa, ok := golden[1].(*wtype.Plate); !ok {
		t.Errorf("expecting %v but got %v instead", golden, x)
	} else if bb, ok := x[1].(*wtype.Plate); !ok {
		t.Errorf("expecting %v but got %v instead", golden, x)
	} else if aa.Type != bb.Type {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestTime(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	type Value map[string]interface{}
	x := Value{
		"A": wunit.Time{},
	}
	golden := Value{
		"A": wunit.NewTime(60.0, "s"),
	}
	if bytes, err := json.Marshal(golden); err != nil {
		t.Error(err)
	} else if err := unmarshal(ctx, &x, bytes); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(x, golden) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestConstructFile(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())
	var x wtype.File

	if err := unmarshal(ctx, &x, []byte(`{"name":"mytest","bytes":{"bytes":"aGVsbG8="}}`)); err != nil {
		t.Fatal(err)
	} else if e, f := "mytest", x.Name; e != f {
		t.Errorf("expecting %v but got %v instead", e, f)
	} else if bs, err := x.ReadAll(); err != nil {
		t.Error(err)
	} else if e, f := "hello", string(bs); e != f {
		t.Errorf("expecting %v but got %v instead", e, f)
	}
}
