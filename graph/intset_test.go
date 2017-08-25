package graph

import (
	"testing"
)

func TestIntSetAdd(t *testing.T) {
	is := NewIntSet()

	if a := is.Add(); a != nil {
		t.Errorf("expected %v found %v", nil, a)
	} else if a, b := is.Add(1), is.Add(1); a == nil || b == nil || a != b {
		t.Errorf("expected %q == %q", a, b)
	} else if a, b := is.Add(1, 2), is.Add(2, 1); a == nil || b == nil || a != b {
		t.Errorf("expected %q == %q", a, b)
	} else if a, b := is.Add(1, 2, 3, 4), is.Add(4, 3, 2, 1); a == nil || b == nil || a != b {
		t.Errorf("expected %q == %q", a, b)
	} else if a, b := is.Add(1, 2, 3, 4, 5), is.Add(5, 4, 3, 2, 1); a == nil || b == nil || a != b {
		t.Errorf("expected %q == %q", a, b)
	}

	if a, b := is.Add(5), is.Add(6); a == nil || b == nil || a == b {
		t.Errorf("expected %q != %q", a, b)
	} else if a, b := is.Add(1, 2), is.Add(2, 3); a == nil || b == nil || a == b {
		t.Errorf("expected %q != %q", a, b)
	} else if a, b := is.Add(1, 2, 3), is.Add(2, 3, 4); a == nil || b == nil || a == b {
		t.Errorf("expected %q != %q", a, b)
	} else if a, b := is.Add(1, 2, 3, 4), is.Add(2, 3); a == nil || b == nil || a == b {
		t.Errorf("expected %q != %q", a, b)
	} else if a, b := is.Add(1, 2, 3, 4), is.Add(2, 3, 4); a == nil || b == nil || a == b {
		t.Errorf("expected %q != %q", a, b)
	} else if a, b := is.Add(4, 5, 6, 7), is.Add(7, 8, 9); a == nil || b == nil || a == b {
		t.Errorf("expected %q != %q", a, b)
	}
}
