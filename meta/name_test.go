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
		{Obj: 0, Expected: "int"},
		{Obj: 0.0, Expected: "float64"},
		{Obj: make(map[int]string), Expected: "map[int]string"},
		{Obj: make([]string, 0), Expected: "[]string"},
		{Obj: make(chan string), Expected: "chan string"},
		{Obj: TestFullyQualifiedType, Expected: "func(*testing.T)"},
		{Obj: myStruct{}, Expected: "github.com/antha-lang/antha/meta.myStruct"},
		{Obj: myString(""), Expected: "github.com/antha-lang/antha/meta.myString"},
		{Obj: func(error) {}, Expected: "func(error)"},
	}

	for _, c := range cases {
		if f, e := FullTypeName(c.Obj), c.Expected; f != e {
			t.Errorf("found %q expected %q", f, e)
		}
	}
}
