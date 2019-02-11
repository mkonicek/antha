package effects

// A NameValue is a name-value pair
type NameValue struct {
	Name  string
	Value string
}

// A Request is set of required device capabilities
type Request struct {
	Selector []NameValue
}

func asSet(vs []NameValue) map[NameValue]struct{} {
	m := make(map[NameValue]struct{})
	for _, v := range vs {
		m[v] = struct{}{}
	}
	return m
}

// returns true iff a is a superset of b (i.e. everything in b must be in a)
func isSuperset(a, b map[NameValue]struct{}) bool {
	for k := range b {
		if _, found := a[k]; !found {
			return false
		}
	}
	return true
}

// Contains returns if request A is greater than or equal to request B
func (reqA Request) Contains(reqB Request) bool {
	return isSuperset(asSet(reqA.Selector), asSet(reqB.Selector))
}

// Meet computes greatest lower bound of a set of requests
func Meet(reqs ...Request) (req Request) {
	for _, r := range reqs {
		req.Selector = append(req.Selector, r.Selector...)
	}
	return
}
