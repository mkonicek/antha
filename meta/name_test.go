package meta

import "testing"

type myStruct struct{}
type myString string

func TestFullyQualifiedType(t *testing.T) {
	type Case struct {
		Obj      interface{}
		Expected string
	}

	cases := []Case{
		Case{Obj: 0, Expected: "int"},
		Case{Obj: 0.0, Expected: "float64"},
		Case{Obj: make(map[int]string), Expected: "map[int]string"},
		Case{Obj: make([]string, 0), Expected: "[]string"},
		Case{Obj: make(chan string), Expected: "chan string"},
		Case{Obj: TestFullyQualifiedType, Expected: "func(*testing.T)"},
		Case{Obj: myStruct{}, Expected: "github.com/antha-lang/antha/meta.myStruct"},
		Case{Obj: myString(""), Expected: "github.com/antha-lang/antha/meta.myString"},
		Case{Obj: func(error) {}, Expected: "func(error)"},
	}

	for _, c := range cases {
		if f, e := FullTypeName(c.Obj), c.Expected; f != e {
			t.Errorf("found %q expected %q", f, e)
		}
	}
}
