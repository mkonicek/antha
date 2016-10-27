package ast

import "testing"

func TestSelector(t *testing.T) {
	reqA := Request{
		Selector: []NameValue{
			NameValue{
				Name:  "alpha",
				Value: "alphavalue",
			},
		},
	}
	reqB := Request{
		Selector: []NameValue{
			NameValue{
				Name:  "beta",
				Value: "betavalue",
			},
		},
	}
	reqAB := Meet(reqA, reqB)

	if !reqA.Contains(reqA) {
		t.Errorf("%v should contain itself", reqA)
	}
	if reqA.Contains(reqB) {
		t.Errorf("%v should not contain %v", reqA, reqB)
	}
	if !reqAB.Contains(reqA) {
		t.Errorf("%v should contain %v", reqAB, reqA)
	}
	if !reqAB.Contains(reqB) {
		t.Errorf("%v should contain %v", reqAB, reqB)
	}
}

func TestMatches(t *testing.T) {
	reqA := Request{Time: NewPoint(1.0)}
	reqB := Request{}
	reqC := Request{Time: NewPoint(2.0)}

	if reqA.Matches(reqB) {
		t.Errorf("%v only %v", reqA, reqB)
	}
	if !reqA.Matches(reqC) {
		t.Errorf("%v not only %v", reqA, reqB)
	}
}
