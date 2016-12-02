package inject

import (
	"context"
	"testing"
)

func TestFuncRunner(t *testing.T) {
	ctx := NewContext(context.Background())
	var result int
	if err := Add(ctx, Name{Repo: "noop"}, &FuncRunner{
		RunFunc: func(context.Context, Value) (Value, error) {
			return Value{"Result": &result}, nil
		},
	}); err != nil {
		t.Fatal(err)
	}

	if out, err := Call(ctx, NameQuery{Repo: "noop"}, nil); err != nil {
		t.Fatal(err)
	} else if r := out["Result"]; r != &result {
		t.Errorf("expecting %p got %p", &result, r)
	}
}

func TestCheckedRunner(t *testing.T) {
	type AddInput struct {
		X, Y int
	}

	type AddOutput struct {
		Sum int
	}

	ctx := NewContext(context.Background())
	if err := Add(ctx, Name{Repo: "add"}, &CheckedRunner{
		RunFunc: func(_ context.Context, value Value) (Value, error) {
			input := AddInput{}
			output := AddOutput{}
			if err := Assign(value, &input); err != nil {
				t.Fatal(err)
			}
			output.Sum = input.X + input.Y
			return MakeValue(output), nil
		},
		In:  &AddInput{},
		Out: &AddOutput{},
	}); err != nil {
		t.Fatal(err)
	}

	type IncompatibleInput struct {
		Z float64
	}
	type IncompatibleOutput struct{}

	x := 1
	y := 2
	var output AddOutput
	var badOutput IncompatibleOutput
	if _, err := Call(ctx, NameQuery{Repo: "add"}, MakeValue(IncompatibleInput{})); err == nil {
		t.Errorf("expecting error on wrong input but got success instead")
	} else if out, err := Call(ctx, NameQuery{Repo: "add"}, MakeValue(AddInput{X: x, Y: y})); err != nil {
		t.Fatal(err)
	} else if err := Assign(out, &badOutput); err == nil {
		t.Errorf("expecting error on wrong output but got success instead: %s", out)
	} else if err := Assign(out, &output); err != nil {
		t.Fatal(err)
	} else if output.Sum != x+y {
		t.Fatalf("expecting %d but got %d", x+y, output.Sum)
	}
}

func TestCheckedRunnerWithInterface(t *testing.T) {
	type AddInput struct {
		X, Y int
	}

	type AddOutput struct {
		Sum   int
		Error error
	}

	ctx := NewContext(context.Background())
	if err := Add(ctx, Name{Repo: "add"}, &CheckedRunner{
		RunFunc: func(_ context.Context, value Value) (Value, error) {
			input := &AddInput{}
			output := &AddOutput{}
			if err := Assign(value, input); err != nil {
				t.Fatal(err)
			}
			output.Sum = input.X + input.Y
			return MakeValue(output), nil
		},
		In:  &AddInput{},
		Out: &AddOutput{},
	}); err != nil {
		t.Fatal(err)
	}

	x := 1
	y := 2
	var output AddOutput
	if out, err := Call(ctx, NameQuery{Repo: "add"}, MakeValue(AddInput{X: x, Y: y})); err != nil {
		t.Fatal(err)
	} else if err := Assign(out, &output); err != nil {
		t.Fatal(err)
	} else if output.Sum != x+y {
		t.Fatalf("expecting %d but got %d", x+y, output.Sum)
	}
}
