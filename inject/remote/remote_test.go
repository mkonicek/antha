package remote

import (
	"context"
	"net"
	"testing"

	"github.com/antha-lang/antha/inject"
)

func haveNetwork() error {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}
	defer ln.Close()

	if _, err := net.Dial("tcp", ln.Addr().String()); err != nil {
		return err
	}
	return nil
}

func TestRunner(t *testing.T) {
	if err := haveNetwork(); err != nil {
		t.Skip("no network")
	}
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	type value struct {
		X int
	}

	fn := func(ctx context.Context, in inject.Value) (inject.Value, error) {
		var v value
		if err := inject.Assign(in, &v); err != nil {
			return nil, err
		}
		return inject.MakeValue(value{X: v.X + 1}), nil
	}

	s := &Server{
		Name:     "Test",
		Listener: ln,
		Runner: &inject.CheckedRunner{
			In:      &value{},
			Out:     &value{},
			RunFunc: fn,
		},
	}

	go func() {
		s.Serve()
	}()

	runner := &Runner{
		Name:    "Test",
		Address: ln.Addr().String(),
		Runner: &inject.CheckedRunner{
			In:  &value{},
			Out: &value{},
		},
	}

	var v value
	if out, err := runner.Run(context.Background(), inject.MakeValue(value{X: 2})); err != nil {
		t.Error(err)
	} else if err := inject.Assign(out, &v); err != nil {
		t.Error(err)
	} else if v.X != 3 {
		t.Errorf("expecting %d found %d", 3, v.X)
	}
}
