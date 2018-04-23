package wutil

import "sync"

type IntSet struct {
	contents []int
	chash    map[int]bool
	lock     sync.Mutex
}

func NewIntSet(s int) IntSet {
	c := make([]int, 0, s)
	ch := make(map[int]bool, s)

	return IntSet{c, ch, sync.Mutex{}}
}

func (is *IntSet) Add(i int) {
	is.lock.Lock()
	defer is.lock.Unlock()
	if !is.chash[i] {
		is.chash[i] = true
		is.contents = append(is.contents, i)
	}
}

func (is *IntSet) Remove(s int) {
	is.lock.Lock()
	defer is.lock.Unlock()
	if is.chash[s] {
		// maintain the order
		c := make([]int, 0, len(is.contents)-1)

		for _, v := range is.contents {
			if v != s {
				c = append(c, v)
			}
		}

		is.chash[s] = false
		is.contents = c
	}
}

func (is *IntSet) AsSlice() []int {
	s := make([]int, len(is.contents))
	copy(s, is.contents)
	return s
}
