package meta

import (
	"testing"
	"time"
)

type B struct {
	I int
	S string
}

type A struct {
	S  string
	T  time.Time
	F  float64
	P  *A
	B  B
	Ss []string
}

func TestMergeDisjoint(t *testing.T) {
	if m, err := ShallowMerge(A{S: "Alpha"}, A{F: 1.0}); err != nil {
		t.Error(err)
	} else if a, ok := m.(A); !ok {
		t.Errorf("expecting %T found %T", a, m)
	} else if e, f := "Alpha", a.S; e != f {
		t.Errorf("expecting %q found %q", e, f)
	} else if e, f := 1.0, a.F; e != f {
		t.Errorf("expecting %f found %f", e, f)
	}

	if m, err := ShallowMerge(&A{S: "Alpha"}, &A{F: 1.0}); err != nil {
		t.Error(err)
	} else if a, ok := m.(*A); !ok {
		t.Errorf("expecting %T found %T", a, m)
	} else if e, f := "Alpha", a.S; e != f {
		t.Errorf("expecting %q found %q", e, f)
	} else if e, f := 1.0, a.F; e != f {
		t.Errorf("expecting %f found %f", e, f)
	}
}

func TestMergeOverride(t *testing.T) {
	if m, err := ShallowMerge(A{S: "Alpha", F: 1.0}, A{S: "Beta"}); err != nil {
		t.Error(err)
	} else if a, ok := m.(A); !ok {
		t.Errorf("expecting %T found %T", a, m)
	} else if e, f := "Beta", a.S; e != f {
		t.Errorf("expecting %q found %q", e, f)
	} else if e, f := 1.0, a.F; e != f {
		t.Errorf("expecting %f found %f", e, f)
	}

	if m, err := ShallowMerge(&A{S: "Alpha", F: 1.0}, &A{S: "Beta"}); err != nil {
		t.Error(err)
	} else if a, ok := m.(*A); !ok {
		t.Errorf("expecting %T found %T", a, m)
	} else if e, f := "Beta", a.S; e != f {
		t.Errorf("expecting %q found %q", e, f)
	} else if e, f := 1.0, a.F; e != f {
		t.Errorf("expecting %f found %f", e, f)
	}
}

func TestMergeEmbedded(t *testing.T) {
	if m, err := ShallowMerge(A{B: B{I: 1, S: "Alpha"}}, A{}); err != nil {
		t.Error(err)
	} else if a, ok := m.(A); !ok {
		t.Errorf("expecting %T found %T", a, m)
	} else if e, f := 1, a.B.I; e != f {
		t.Errorf("expecting %q found %q", e, f)
	} else if e, f := "Alpha", a.B.S; e != f {
		t.Errorf("expecting %q found %q", e, f)
	}

	if m, err := ShallowMerge(&A{B: B{I: 1, S: "Alpha"}}, &A{}); err != nil {
		t.Error(err)
	} else if a, ok := m.(*A); !ok {
		t.Errorf("expecting %T found %T", a, m)
	} else if e, f := 1, a.B.I; e != f {
		t.Errorf("expecting %q found %q", e, f)
	} else if e, f := "Alpha", a.B.S; e != f {
		t.Errorf("expecting %q found %q", e, f)
	}
}

func TestMergeSlices(t *testing.T) {
	if m, err := ShallowMerge(A{Ss: []string{"a"}}, A{}); err != nil {
		t.Error(err)
	} else if a, ok := m.(A); !ok {
		t.Errorf("expecting %T found %T", a, m)
	} else if e, f := []string{"a"}, a.Ss; len(e) != len(f) {
		t.Errorf("expecting %q found %q", e, f)
	} else if e[0] != f[0] {
		t.Errorf("expecting %q found %q", e, f)
	}
}

func TestMergeTime(t *testing.T) {
	someTime := time.Unix(1, 0)
	if m, err := ShallowMerge(A{T: someTime}, A{}); err != nil {
		t.Error(err)
	} else if a, ok := m.(A); !ok {
		t.Errorf("expecting %T found %T", a, m)
	} else if e, f := someTime, a.T; e != f {
		t.Errorf("expecting %q found %q", e, f)
	}

	if m, err := ShallowMerge(A{}, A{T: someTime}); err != nil {
		t.Error(err)
	} else if a, ok := m.(A); !ok {
		t.Errorf("expecting %T found %T", a, m)
	} else if e, f := someTime, a.T; e != f {
		t.Errorf("expecting %q found %q", e, f)
	}

	if m, err := ShallowMerge(&A{T: someTime}, &A{}); err != nil {
		t.Error(err)
	} else if a, ok := m.(*A); !ok {
		t.Errorf("expecting %T found %T", a, m)
	} else if e, f := someTime, a.T; e != f {
		t.Errorf("expecting %q found %q", e, f)
	}

	if m, err := ShallowMerge(&A{}, &A{T: someTime}); err != nil {
		t.Error(err)
	} else if a, ok := m.(*A); !ok {
		t.Errorf("expecting %T found %T", a, m)
	} else if e, f := someTime, a.T; e != f {
		t.Errorf("expecting %q found %q", e, f)
	}
}
