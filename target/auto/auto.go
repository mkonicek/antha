// Package auto provides methods for creating a simulation target based on
// auto discovery of drivers via gRPC
package auto

import (
	"context"

	"github.com/antha-lang/antha/ast"
	runner "github.com/antha-lang/antha/driver/antha_runner_v1"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
	"google.golang.org/grpc"
)

// An Endpoint is a network address of a device plugin (driver)
type Endpoint struct {
	URI string
	Arg interface{}
}

// An Opt are options for connecting to a set of device plugins (drivers)
type Opt struct {
	Endpoints []Endpoint
	MaybeArgs []interface{}
}

// An Auto contains the state of autodiscovery of device plugins
type Auto struct {
	Target  *target.Target
	Conns   []*grpc.ClientConn
	runners map[string][]runner.RunnerClient
	handler map[ast.Device]*grpc.ClientConn
}

// Close releases any resources like network connections associated
// with auto-discovery state.
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

// New makes target by inspecting a set of gRPC network services for a list
// of drivers.
func New(opt Opt) (ret *Auto, err error) {
	ret = &Auto{
		Target:  target.New(),
		runners: make(map[string][]runner.RunnerClient),
		handler: make(map[ast.Device]*grpc.ClientConn),
	}

	defer func() {
		if err == nil {
			return
		}
		err = ret.Close()
	}()

	tryer := &tryer{
		Auto:      ret,
		MaybeArgs: opt.MaybeArgs,
		HumanOpt:  human.Opt{CanMix: true, CanIncubate: true},
	}

	ctx := context.Background()
	for _, ep := range opt.Endpoints {
		var conn *grpc.ClientConn
		conn, err = grpc.Dial(ep.URI, grpc.WithInsecure())
		if err != nil {
			return
		}
		ret.Conns = append(ret.Conns, conn)

		if err = tryer.Try(ctx, conn, ep.Arg); err != nil {
			return
		}
	}

	ret.Target.AddDevice(human.New(tryer.HumanOpt))

	return
}
