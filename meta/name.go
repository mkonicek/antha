package meta

import (
	"fmt"
	"reflect"
	"strings"
)

// FullTypeName returns the full type name of an object
func FullTypeName(obj interface{}) string {
	return NameFromType(reflect.TypeOf(obj))
}

// NameFromType returns the full type name from a reflect.Type
func NameFromType(t reflect.Type) string {
	// Like join but ignore empty strings
	joinWith := func(sep string, args ...string) string {
		var vs []string
		for _, arg := range args {
			if len(arg) != 0 {
				vs = append(vs, arg)
			}
		}
		return strings.Join(vs, sep)
	}

	switch t.Kind() {

	case reflect.Array:
		return fmt.Sprintf("[%d]%s", t.Len(), NameFromType(t.Elem()))

	case reflect.Chan:
		switch t.ChanDir() {

		case reflect.BothDir:
			return fmt.Sprintf("chan %s", NameFromType(t.Elem()))

		case reflect.RecvDir:
			return fmt.Sprintf("<-chan %s", NameFromType(t.Elem()))

		case reflect.SendDir:
			return fmt.Sprintf("chan<- %s", NameFromType(t.Elem()))
		}

	case reflect.Func:
		var ins []string
		for i, n := 0, t.NumIn(); i < n; i++ {
			ins = append(ins, NameFromType(t.In(i)))
		}
		var outs []string
		for i, n := 0, t.NumOut(); i < n; i++ {
			outs = append(outs, NameFromType(t.In(i)))
		}
		return joinWith(" ", fmt.Sprintf("func(%s)", strings.Join(ins, ",")), strings.Join(outs, ","))

	case reflect.Interface:
		return joinWith(".", t.PkgPath(), t.Name())

	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", NameFromType(t.Key()), NameFromType(t.Elem()))

	case reflect.Ptr:
		return fmt.Sprintf("*%s", NameFromType(t.Elem()))

	case reflect.Slice:
		return fmt.Sprintf("[]%s", NameFromType(t.Elem()))

	}

	return joinWith(".", t.PkgPath(), t.Name())
}
