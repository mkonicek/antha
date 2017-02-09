package execute

import (
	"reflect"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/meta"
	"github.com/antha-lang/antha/microArch/factory"
)

func unmarshal(obj interface{}, data []byte) error {
	var u unmarshaler
	if err := meta.UnmarshalJSON(meta.UnmarshalOpt{
		Struct: u.unmarshalStruct,
	}, data, obj); err != nil {
		return err
	}
	return nil
}

func TestString(t *testing.T) {
	type Value string
	var x Value
	golden := Value("hello")

	if err := unmarshal(&x, []byte(`"hello"`)); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(x, golden) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestInt(t *testing.T) {
	type Value int
	var x Value
	golden := Value(1)
	if err := unmarshal(&x, []byte(`1`)); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(x, golden) {
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

	if err := unmarshal(&x, []byte(`{"A": "hello", "B": 1}`)); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(x, golden) {
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
	if err := unmarshal(&x, []byte(`{"A": {"A": "hello", "B": 1}, "B": {"A": "hello", "B": 2} }`)); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(x, golden) {
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
	if err := unmarshal(&x, []byte(`[ {"A": "hello", "B": 1}, {"A": "hello", "B": 2} ]`)); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(x, golden) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestConstruct(t *testing.T) {
	var x wtype.LHTipbox
	if err := unmarshal(&x, []byte(`"CyBio250Tipbox"`)); err != nil {
		t.Fatal(err)
	}
}

func TestConstructMapFailure(t *testing.T) {
	type Elem struct {
		A string
		T *wtype.LHTipbox
	}
	type Value map[string]Elem
	var x Value
	if err := unmarshal(&x, []byte(`{"A": {"A": "hello", "T": "CyBio250Tipbox"} }`)); err == nil {
		t.Fatal("expecting failure but got success")
	}
}

func TestConstructMap(t *testing.T) {
	type Value map[string]interface{}
	x := Value{
		"A": &wtype.LHTipbox{},
		"B": 0,
		"C": "",
	}
	golden := Value{
		"A": factory.GetTipboxByType("CyBio250Tipbox"),
		"B": 1,
		"C": "hello",
	}
	if err := unmarshal(&x, []byte(`{"A": "CyBio250Tipbox", "B": 1, "C": "hello" }`)); err != nil {
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
	type Value []interface{}
	x := Value{
		&wtype.LHTipbox{},
		&wtype.LHPlate{},
	}
	golden := Value{
		factory.GetTipboxByType("CyBio250Tipbox"),
		factory.GetPlateByType("pcrplate_with_cooler"),
	}
	if err := unmarshal(&x, []byte(`[ "CyBio250Tipbox", "pcrplate_with_cooler" ]`)); err != nil {
		t.Fatal(err)
	} else if len(x) != 2 {
		t.Errorf("expecting %v but got %v instead", golden, x)
	} else if aa, ok := golden[0].(*wtype.LHTipbox); !ok {
		t.Errorf("expecting %v but got %v instead", golden, x)
	} else if bb, ok := x[0].(*wtype.LHTipbox); !ok {
		t.Errorf("expecting %v but got %v instead", golden, x)
	} else if aa.Type != bb.Type {
		t.Errorf("expecting %v but got %v instead", golden, x)
	} else if aa, ok := golden[1].(*wtype.LHPlate); !ok {
		t.Errorf("expecting %v but got %v instead", golden, x)
	} else if bb, ok := x[1].(*wtype.LHPlate); !ok {
		t.Errorf("expecting %v but got %v instead", golden, x)
	} else if aa.Type != bb.Type {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestTime(t *testing.T) {
	type Value map[string]interface{}
	x := Value{
		"A": wunit.Time{},
	}
	golden := Value{
		"A": wunit.NewTime(60.0, "s"),
	}
	if err := unmarshal(&x, []byte(`{ "A": "60s" }`)); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(x, golden) {
		t.Errorf("expecting %v but got %v instead", golden, x)
	}
}

func TestConstructFile(t *testing.T) {
	var x wtype.File

	if err := unmarshal(&x, []byte(`{"name":"mytest","bytes":{"bytes":"aGVsbG8="}}`)); err != nil {
		t.Fatal(err)
	} else if e, f := "mytest", x.Name; e != f {
		t.Errorf("expecting %v but got %v instead", e, f)
	} else if bs, err := x.ReadAll(); err != nil {
		t.Error(err)
	} else if e, f := "hello", string(bs); e != f {
		t.Errorf("expecting %v but got %v instead", e, f)
	}
}
