package wutil

import "reflect"

func StrInStrArray(s string, a []string) bool {
	for _, v := range a {
		if v == s {
			return true
		}
	}

	return false
}

func NUniqueStringsInArray(a []string) int {
	m := make(map[string]bool, len(a))

	for _, v := range a {
		m[v] = true
	}

	return len(reflect.ValueOf(m).MapKeys())
}
