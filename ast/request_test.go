package ast

import "testing"

func TestSelector(t *testing.T) {
	reqA := Request{
		Selector: []NameValue{
			{
				Name:  "alpha",
				Value: "alphavalue",
			},
		},
	}
	reqB := Request{
		Selector: []NameValue{
			{
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
