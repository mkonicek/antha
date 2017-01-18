//
package auto

import (
	"context"
	"errors"

	"github.com/antha-lang/antha/bvendor/google.golang.org/grpc"
	runner "github.com/antha-lang/antha/driver/antha_runner_v1"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
)

var (
	noMatch = errors.New("no match")
)

type Endpoint struct {
	Uri string
	Arg interface{}
}

type Opt struct {
	Endpoints []Endpoint
	MaybeArgs []interface{}
}

type Auto struct {
	Target  *target.Target
	Conns   []*grpc.ClientConn
	runners map[string][]runner.RunnerClient
	handler map[target.Device]*grpc.ClientConn
}

func (a *Auto) Close() error {
	var err error
	for _, conn := range a.Conns {
		e := conn.Close()
		if err == nil {
			err = e
		}
	}
	return err
}

// Make target by inspecting a set of network services
func New(opt Opt) (*Auto, error) {
	var err error

	ret := &Auto{
		Target:  target.New(),
		runners: make(map[string][]runner.RunnerClient),
		handler: make(map[target.Device]*grpc.ClientConn),
	}

	defer func() {
		if err == nil {
			return
		}
		ret.Close()
	}()

	tryer := &tryer{
		Auto:      ret,
		MaybeArgs: opt.MaybeArgs,
		HumanOpt:  human.Opt{CanMix: true, CanIncubate: true, CanHandle: true},
	}

	ctx := context.Background()
	for _, ep := range opt.Endpoints {
		var conn *grpc.ClientConn
		conn, err = grpc.Dial(ep.Uri, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		ret.Conns = append(ret.Conns, conn)

		if err = tryer.Try(ctx, conn, ep.Arg); err != nil {
			return nil, err
		}
	}

	if err = ret.Target.AddDevice(human.New(tryer.HumanOpt)); err != nil {
		return nil, err
	}
	return ret, nil
}
