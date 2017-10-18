package Optimization

import (
	"fmt"
)

type SeqSet []string

/*
func (ss SeqSet) Split(pts []int) []SeqSet {
	for _, v := range pts {

	}
}
*/

// p in [0,len(ss[...])-1]
func (ss SeqSet) SplitAt(p int) (Before, After SeqSet, e error) {
	if len(ss) < 1 || p < len(ss[0]) || p > len(ss[0])-1 {
		return Before, After, fmt.Errorf("Invalid split point %d", p)
	}

	Before = make(SeqSet, len(ss))
	After = make(SeqSet, len(ss))

	for i := 0; i < len(ss); i++ {
		Before[i] = ss[i][0:p]
		After[i] = ss[i][p : len(ss[i])-1]
	}

	return Before, After, nil
}

func (ss SeqSet) Valid() bool {
	l := -1
	for _, v := range ss {
		if l == -1 {
			l = len(v)
		}

		if l != len(v) {
			return false
		}
	}

	return true
}

func Distinct(sa []string) []string {
	m := make(map[string]bool)
	r := make([]string, 0, len(sa))
	for _, s := range sa {
		if !m[s] {
			r = append(r, s)
			m[s] = true
		}

	}
	return r
}
