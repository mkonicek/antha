package ast

// A NameValue is a name-value pair
type NameValue struct {
	Name  string
	Value string
}

// A Request is set of required device capabilities
type Request struct {
	Selector []NameValue
}

func makeNameValueMap(vs []NameValue) map[interface{}]int {
	m := make(map[interface{}]int)
	for _, v := range vs {
		m[v]++
	}
	return m
}

func mapContains(a, b map[interface{}]int) bool {
	for k, v := range b {
		if v > a[k] {
			return false
		}
	}
	return true
}

// Contains returns if request A is greater than or equal to request B
func (reqA Request) Contains(reqB Request) bool {
	return mapContains(makeNameValueMap(reqA.Selector), makeNameValueMap(reqB.Selector))
}

// Meet computes greatest lower bound of a set of requests
func Meet(reqs ...Request) (req Request) {
	for _, r := range reqs {
		req.Selector = append(req.Selector, r.Selector...)
	}
	return
}
