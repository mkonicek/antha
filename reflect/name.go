package reflect

import (
	"fmt"
	"reflect"
	"strings"
)

func FullTypeName(t reflect.Type) string {
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
		return fmt.Sprintf("[%d]%s", t.Len(), FullTypeName(t.Elem()))
	case reflect.Chan:
		switch t.ChanDir() {
		case reflect.BothDir:
			return fmt.Sprintf("chan %s", FullTypeName(t.Elem()))
		case reflect.RecvDir:
			return fmt.Sprintf("<-chan %s", FullTypeName(t.Elem()))
		case reflect.SendDir:
			return fmt.Sprintf("chan<- %s", FullTypeName(t.Elem()))
		}
	case reflect.Func:
		var ins []string
		for i, n := 0, t.NumIn(); i < n; i += 1 {
			ins = append(ins, FullTypeName(t.In(i)))
		}
		var outs []string
		for i, n := 0, t.NumOut(); i < n; i += 1 {
			outs = append(outs, FullTypeName(t.In(i)))
		}
		return joinWith(" ", fmt.Sprintf("func(%s)", strings.Join(ins, ",")), strings.Join(outs, ","))
	case reflect.Interface:
		return joinWith(".", t.PkgPath(), t.Name())
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", FullTypeName(t.Key()), FullTypeName(t.Elem()))
	case reflect.Ptr:
		return fmt.Sprintf("*%s", FullTypeName(t.Elem()))
	case reflect.Slice:
		return fmt.Sprintf("[]%s", FullTypeName(t.Elem()))
	case reflect.Struct:
		return joinWith(".", t.PkgPath(), t.Name())
	}
	return t.Name()
}
