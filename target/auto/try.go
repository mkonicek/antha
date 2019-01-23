package auto

import (
	"context"

	driver "github.com/antha-lang/antha/driver/antha_driver_v1"
	runner "github.com/antha-lang/antha/driver/antha_runner_v1"
	lhclient "github.com/antha-lang/antha/driver/liquidhandling/client"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/target/handler"
	"github.com/antha-lang/antha/target/human"
	"github.com/antha-lang/antha/target/mixer"
	"github.com/antha-lang/antha/target/shakerincubator"
	"google.golang.org/grpc"
)

// Common state for tryers
type tryer struct {
	Auto      *Auto
	MaybeArgs []interface{}
	HumanOpt  human.Opt
}

// AddDriver queries a driver and adds the corresponding device to the target
// based on the query response
func (a *tryer) AddDriver(ctx context.Context, conn *grpc.ClientConn, arg interface{}) error {
	c := driver.NewDriverClient(conn)
	reply, err := c.DriverType(ctx, &driver.TypeRequest{})
	if err != nil {
		return err
	}

	switch reply.Type {

	case "antha.runner.v1.Runner":
		r := runner.NewRunnerClient(conn)
		reply, err := r.SupportedRunTypes(ctx, &runner.SupportedRunTypesRequest{})
		if err != nil {
			return err
		}
		for _, typ := range reply.Types {
			a.Auto.runners[typ] = append(a.Auto.runners[typ], r)
		}

	case "antha.shakerincubator.v1.ShakerIncubator":
		s := &shakerincubator.ShakerIncubator{}
		a.HumanOpt.CanIncubate = false
		a.Auto.handler[s] = conn
		a.Auto.Target.AddDevice(s)
		return nil

	default:
		h := handler.New(
			[]effects.NameValue{
				{
					Name:  "antha.driver.v1.TypeReply.type",
					Value: reply.Type,
				},
			},
		)
		a.HumanOpt.CanHandle = false
		a.Auto.handler[h] = conn
		a.Auto.Target.AddDevice(h)
		return nil
	}

	return nil
}

// AddMixer queries a mixer driver and adds the corresponding device to the target
func (a *tryer) AddMixer(ctx context.Context, conn *grpc.ClientConn, arg interface{}) error {
	err := a.addHighLevelMixer(ctx, conn, arg)

	if err == nil {
		return nil
	}

	err = a.addLowLevelMixer(ctx, conn, arg)

	if err == nil {
		return nil
	}

	err = a.addHighLevelMixer(ctx, conn, arg)

	return err
}

func (a *tryer) addHighLevelMixer(ctx context.Context, conn *grpc.ClientConn, arg interface{}) error {
	var candidates []interface{}
	candidates = append(candidates, arg)
	candidates = append(candidates, a.MaybeArgs...)

	d, err := mixer.New(getMixerOpt(candidates), lhclient.NewHighLevelClientFromConn(conn))
	if err != nil {
		return err
	}

	a.HumanOpt.CanMix = false
	a.Auto.Target.AddDevice(d)
	return nil
}
func (a *tryer) addLowLevelMixer(ctx context.Context, conn *grpc.ClientConn, arg interface{}) error {

	var candidates []interface{}
	candidates = append(candidates, arg)
	candidates = append(candidates, a.MaybeArgs...)

	d, err := mixer.New(getMixerOpt(candidates), lhclient.NewLowLevelClientFromConn(conn))
	if err != nil {
		return err
	}

	a.HumanOpt.CanMix = false
	a.Auto.Target.AddDevice(d)
	return nil
}

func getMixerOpt(maybeArgs []interface{}) (ret mixer.Opt) {
	for _, v := range maybeArgs {
		if o, ok := v.(mixer.Opt); ok {
			return o
		}
	}
	return
}

func (a *tryer) Try(ctx context.Context, conn *grpc.ClientConn, arg interface{}) error {
	var tries []func(context.Context, *grpc.ClientConn, interface{}) error
	tries = append(tries, a.AddDriver, a.AddMixer)

	for _, t := range tries {
		if err := t(ctx, conn, arg); err == nil {
			return nil
		}
	}

	return errNoMatch
}
