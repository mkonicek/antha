package injectutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/antha-lang/antha/inject"
)

func TestProduct(t *testing.T) {
	type Pair struct {
		A int
		B string
		C float64
	}

	space := ProductSpace{
		"A": []interface{}{1, 2},
		"B": []interface{}{"A", "B"},
		"C": []interface{}{1.5, 2.5},
	}

	expected := make(map[Pair]bool)
	expected[Pair{A: 1, B: "A", C: 1.5}] = true
	expected[Pair{A: 1, B: "B", C: 1.5}] = true
	expected[Pair{A: 2, B: "A", C: 1.5}] = true
	expected[Pair{A: 2, B: "B", C: 1.5}] = true
	expected[Pair{A: 1, B: "A", C: 2.5}] = true
	expected[Pair{A: 1, B: "B", C: 2.5}] = true
	expected[Pair{A: 2, B: "A", C: 2.5}] = true
	expected[Pair{A: 2, B: "B", C: 2.5}] = true

	seen := make(map[Pair]bool)
	for _, elem := range space.Elements() {
		var pair Pair

		if l := len(elem); l != 3 {
			t.Fatalf("expected 3 components but found %d", l)
		} else if err := inject.Assign(elem, &pair); err != nil {
			t.Fatal(err)
		} else if _, ok := expected[pair]; !ok {
			t.Fatalf("unexpected element %q", pair)
		} else if _, ok := seen[pair]; ok {
			t.Fatalf("already seen element %q", pair)
		} else {
			seen[pair] = true
		}
	}
	if le, ls := len(expected), len(seen); le != ls {
		t.Fatalf("expected %d elements but found %d", le, ls)
	}
}

func TestCartesianProduct(t *testing.T) {
	type Pair struct {
		A string
		B int
		C int
	}

	ctx := context.Background()

	runner := &inject.FuncRunner{
		RunFunc: func(_ context.Context, value inject.Value) (inject.Value, error) {
			var pair Pair
			if err := inject.Assign(value, &pair); err != nil {
				return nil, err
			} else {
				return inject.Value{"Output": fmt.Sprintf("%s%d", pair.A, pair.B*pair.C)}, nil
			}
		},
	}

	space := ProductSpace{
		"A": []interface{}{"A", "B"},
		"B": []interface{}{1, 2, 3},
		"C": []interface{}{4},
	}
	expected := make(map[string]bool)
	expected["A4"] = true
	expected["A8"] = true
	expected["A12"] = true
	expected["B4"] = true
	expected["B8"] = true
	expected["B12"] = true

	products, err := CallCartesianProduct(ctx, runner, nil, space)
	if err != nil {
		t.Fatal(err)
	}

	seen := make(map[string]bool)
	for _, p := range products {
		t.Log(p)
		if out, ok := p.Output["Output"]; !ok {
			t.Fatal("missing output")
		} else if s, ok := out.(string); !ok {
			t.Fatalf("expected %T but found %T", s, out)
		} else if _, ok := expected[s]; !ok {
			t.Fatalf("already seen element %q", s)
		} else {
			seen[s] = true
		}
	}

	if le, ls := len(expected), len(seen); le != ls {
		t.Fatalf("expected %d elements but found %d", le, ls)
	}
}

func TestCartesianProductWithFixed(t *testing.T) {
	type Pair struct {
		A string
		B int
		C int
	}

	ctx := context.Background()

	runner := &inject.FuncRunner{
		RunFunc: func(_ context.Context, value inject.Value) (inject.Value, error) {
			var pair Pair
			if err := inject.Assign(value, &pair); err != nil {
				return nil, err
			} else {
				return inject.Value{"Output": fmt.Sprintf("%s%d", pair.A, pair.B*pair.C)}, nil
			}
		},
	}

	space := ProductSpace{
		"B": []interface{}{1, 2, 3},
		"C": []interface{}{4, 5},
	}
	expected := make(map[string]bool)
	expected["A4"] = true
	expected["A5"] = true
	expected["A8"] = true
	expected["A10"] = true
	expected["A12"] = true
	expected["A15"] = true

	products, err := CallCartesianProduct(ctx, runner, inject.Value{"A": "A"}, space)
	if err != nil {
		t.Fatal(err)
	}

	seen := make(map[string]bool)
	for _, p := range products {
		t.Log(p)
		if out, ok := p.Output["Output"]; !ok {
			t.Fatal("missing output")
		} else if s, ok := out.(string); !ok {
			t.Fatalf("expected %T but found %T", s, out)
		} else if _, ok := expected[s]; !ok {
			t.Fatalf("already seen element %q", s)
		} else {
			seen[s] = true
		}
	}

	if le, ls := len(expected), len(seen); le != ls {
		t.Fatalf("expected %d elements but found %d", le, ls)
	}
}
