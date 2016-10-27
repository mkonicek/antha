package ast

// TODO make more specific
type Location interface{}

type Movement struct {
	From Location
	To   Location
}

type NameValue struct {
	Name  string
	Value string
}

type Request struct {
	MixVol   *Interval
	Temp     *Interval
	Time     *Interval
	Move     []Movement
	Selector []NameValue
}

func makeMovementMap(vs []Movement) map[interface{}]int {
	m := make(map[interface{}]int)
	for _, v := range vs {
		m[v] += 1
	}
	return m
}

func makeNameValueMap(vs []NameValue) map[interface{}]int {
	m := make(map[interface{}]int)
	for _, v := range vs {
		m[v] += 1
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

// A >= B?
func (reqA Request) Contains(reqB Request) bool {
	if !reqA.MixVol.Contains(reqB.MixVol) {
		return false
	}
	if !reqA.Temp.Contains(reqB.Temp) {
		return false
	}
	if !reqA.Time.Contains(reqB.Time) {
		return false
	}
	if !mapContains(makeMovementMap(reqA.Move), makeMovementMap(reqB.Move)) {
		return false
	}
	if !mapContains(makeNameValueMap(reqA.Selector), makeNameValueMap(reqB.Selector)) {
		return false
	}
	return true
}

// Compute greatest lower bound of a set of requests
func Meet(reqs ...Request) (req Request) {
	for _, r := range reqs {
		req.MixVol = req.MixVol.Meet(r.MixVol)
		req.Temp = req.Temp.Meet(r.Temp)
		req.Time = req.Time.Meet(r.Time)
		req.Move = append(req.Move, r.Move...)
		req.Selector = append(req.Selector, r.Selector...)
	}
	return
}

// Checks if the non-zero fields of A are a subset of the non-zero fields of B.
func (reqA Request) Matches(reqB Request) bool {
	if reqA.MixVol != nil && reqB.MixVol == nil {
		return false
	}
	if reqA.Temp != nil && reqB.Temp == nil {
		return false
	}
	if reqA.Time != nil && reqB.Time == nil {
		return false
	}
	if len(reqA.Move) != 0 && len(reqB.Move) == 0 {
		return false
	}
	if len(reqA.Selector) != 0 && len(reqB.Selector) == 0 {
		return false
	}
	return true
}
