package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	api "github.com/antha-lang/antha/api/v1"
	"github.com/antha-lang/antha/inject"
)

func createContext() (context.Context, error) {
	ctx := inject.NewContext(context.Background())

	if err := inject.Add(ctx, inject.Name{Repo: "Equals", Stage: api.ElementStage_STEPS}, &inject.FuncRunner{
		RunFunc: func(_ context.Context, value inject.Value) (inject.Value, error) {
			if a, ok := value["A"].(string); !ok {
				return nil, fmt.Errorf("cannot read parameter A")
			} else if b, ok := value["B"].(string); !ok {
				return nil, fmt.Errorf("cannot read parameter B")
			} else {
				return map[string]interface{}{"Out": a == b}, nil
			}
		},
	}); err != nil {
		return nil, err
	}
	if err := inject.Add(ctx, inject.Name{Repo: "Cond", Stage: api.ElementStage_STEPS}, &inject.FuncRunner{
		RunFunc: func(_ context.Context, value inject.Value) (inject.Value, error) {
			if a, ok := value["True"].(string); !ok {
				return nil, fmt.Errorf("cannot read parameter True")
			} else if b, ok := value["False"].(string); !ok {
				return nil, fmt.Errorf("cannot read parameter False")
			} else if c, ok := value["Cond"].(bool); !ok {
				return nil, fmt.Errorf("cannot read parameter Cond")
			} else if c {
				return map[string]interface{}{"Out": a}, nil
			} else {
				return map[string]interface{}{"Out": b}, nil
			}
		},
	}); err != nil {
		return nil, err
	}
	if err := inject.Add(ctx, inject.Name{Repo: "Copy", Stage: api.ElementStage_STEPS}, &inject.FuncRunner{
		RunFunc: func(_ context.Context, value inject.Value) (inject.Value, error) {
			a, ok := value["In"].(string)
			if !ok {
				return nil, fmt.Errorf("cannot read parameter In")
			}
			return map[string]interface{}{"Out": a}, nil
		},
	}); err != nil {
		return nil, err
	}
	return ctx, nil
}

func TestRunFromFile(t *testing.T) {
	var desc *Desc
	if err := json.Unmarshal([]byte(condCopyEqualsJSON), &desc); err != nil {
		t.Fatal(err)
	}

	w, err := New(Opt{FromDesc: desc})
	if err != nil {
		t.Fatal(err)
	}

	ctx, err := createContext()
	if err != nil {
		t.Error(err)
	}

	if err := w.SetParam(Port{Process: "Equals", Port: "A"}, "A"); err != nil {
		t.Error(err)
	}
	if err := w.SetParam(Port{Process: "Equals", Port: "B"}, "B"); err != nil {
		t.Error(err)
	}
	if err := w.SetParam(Port{Process: "Cond", Port: "True"}, "True"); err != nil {
		t.Error(err)
	}
	if err := w.SetParam(Port{Process: "Cond", Port: "False"}, "False"); err != nil {
		t.Error(err)
	}
	if err := w.Run(ctx); err != nil {
		t.Error(err)
	}

	if out, ok := w.Outputs[Port{Process: "Copy", Port: "Out"}].(string); !ok {
		t.Errorf("cannot read parameter Out")
	} else if out != "False" {
		t.Errorf("expecting output %q but got %q", "False", out)
	}
}

func TestRun(t *testing.T) {
	w, err := New(Opt{})
	if err != nil {
		t.Fatal(err)
	}

	ctx, err := createContext()
	if err != nil {
		t.Error(err)
	}

	if err := w.AddNode("Equals", "Equals"); err != nil {
		t.Error(err)
	}
	if err := w.AddNode("Cond", "Cond"); err != nil {
		t.Error(err)
	}
	if err := w.AddNode("Copy", "Copy"); err != nil {
		t.Error(err)
	}

	if err := w.AddEdge(Port{Process: "Equals", Port: "Out"}, Port{Process: "Cond", Port: "Cond"}); err != nil {
		t.Error(err)
	}
	if err := w.AddEdge(Port{Process: "Cond", Port: "Out"}, Port{Process: "Copy", Port: "In"}); err != nil {
		t.Error(err)
	}

	if err := w.SetParam(Port{Process: "Equals", Port: "A"}, "A"); err != nil {
		t.Error(err)
	}
	if err := w.SetParam(Port{Process: "Equals", Port: "B"}, "B"); err != nil {
		t.Error(err)
	}
	if err := w.SetParam(Port{Process: "Cond", Port: "True"}, "True"); err != nil {
		t.Error(err)
	}
	if err := w.SetParam(Port{Process: "Cond", Port: "False"}, "False"); err != nil {
		t.Error(err)
	}

	if err := w.Run(ctx); err != nil {
		t.Error(err)
	}

	if out, ok := w.Outputs[Port{Process: "Copy", Port: "Out"}].(string); !ok {
		t.Errorf("cannot read parameter Out")
	} else if out != "False" {
		t.Errorf("expecting output %q but got %q", "False", out)
	}
}

func TestDuplicateIns(t *testing.T) {
	w, err := New(Opt{})
	if err != nil {
		t.Fatal(err)
	}

	if err := w.AddNode("A", "A"); err != nil {
		t.Error(err)
	}
	if err := w.AddNode("B", "B"); err != nil {
		t.Error(err)
	}

	if err := w.AddEdge(Port{Process: "A", Port: "Out"}, Port{Process: "B", Port: "In"}); err != nil {
		t.Error(err)
	}
	if err := w.AddEdge(Port{Process: "A", Port: "Out"}, Port{Process: "B", Port: "In"}); err == nil {
		t.Errorf("expecting error setting in port")
	}
}
