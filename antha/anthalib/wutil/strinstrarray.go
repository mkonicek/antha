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

// NUniqueStringsInArray counts the unique strings in an array
// i.e. if 1, all strings in the array are the same.
func NUniqueStringsInArray(a []string, excludeBlanks bool) int {
	m := make(map[string]bool, len(a))

	for _, v := range a {
		if !(v == "" && excludeBlanks) {
			m[v] = true
		}
	}

	return len(reflect.ValueOf(m).MapKeys())
}

func StringArrayEqual(a1, a2 []string) bool {
	if len(a1) != len(a2) {
		return false
	}

	for i := 0; i < len(a1); i++ {
		if a1[i] != a2[i] {
			return false
		}
	}

	return true
}
