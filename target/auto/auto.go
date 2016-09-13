package auto

import (
	"errors"

	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/bvendor/google.golang.org/grpc"
	driver "github.com/antha-lang/antha/driver/antha_driver_v1"
	runner "github.com/antha-lang/antha/driver/antha_runner_v1"
	lhclient "github.com/antha-lang/antha/driver/lh"
	"github.com/antha-lang/antha/driver/pb/lh"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/mixer"
)

var (
	noMatch = errors.New("no match")
)

var (
	tryers []tryer = []tryer{
		tryRunner,
		tryMixer,
	}
)

type tryer func(*grpc.ClientConn, []interface{}) (target.Device, error)

type Opt struct {
	Uri  string
	Opts []interface{}
}

func tryRunner(conn *grpc.ClientConn, opts []interface{}) (target.Device, error) {
	c := driver.NewDriverClient(conn)
	reply, err := c.DriverType(context.Background(), &driver.TypeRequest{})
	if err != nil {
		return nil, err
	}
	if reply.Type != "antha.runner.v1.Runner" {
		return nil, noMatch
	}

	return target.NewRunner(runner.NewRunnerClient(conn)), nil
}

func getMixerOpt(opt []interface{}) (ret mixer.Opt) {
	for _, v := range opt {
		if o, ok := v.(mixer.Opt); ok {
			return o
		}
	}
	return
}

func tryMixer(conn *grpc.ClientConn, opts []interface{}) (target.Device, error) {
	c := lh.NewExtendedLiquidhandlingDriverClient(conn)
	return mixer.New(getMixerOpt(opts), &lhclient.Driver{C: c})
}

func New(opt Opt) (target.Device, error) {
	conn, err := grpc.Dial(opt.Uri, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

	var d target.Device
	for _, t := range tryers {
		if d, err = t(conn, opt.Opts); err == nil {
			return d, nil
		}
	}
	return nil, noMatch
}
