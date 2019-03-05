package auto

import (
	"context"
	"errors"
	"fmt"

	driver "github.com/antha-lang/antha/driver/antha_driver_v1"
	runner "github.com/antha-lang/antha/driver/antha_runner_v1"
	lhclient "github.com/antha-lang/antha/driver/liquidhandling/client"
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

// Try queries a driver and adds the corresponding device to the target
// based on the query response
func (a *tryer) Try(ctx context.Context, conn *grpc.ClientConn, arg interface{}) error {
	c := driver.NewDriverClient(conn)
	reply, err := c.DriverType(ctx, &driver.TypeRequest{})
	if err != nil {
		return err
	}

	switch reply.GetType() {

	case "antha.runner.v1.Runner":
		r := runner.NewRunnerClient(conn)
		reply, err := r.SupportedRunTypes(ctx, &runner.SupportedRunTypesRequest{})
		if err != nil {
			return err
		}
		for _, typ := range reply.Types {
			a.Auto.runners[typ] = append(a.Auto.runners[typ], r)
		}
		return nil

	case "antha.shakerincubator.v1.ShakerIncubator":
		s := &shakerincubator.ShakerIncubator{}
		a.HumanOpt.CanIncubate = false
		a.Auto.handler[s] = conn
		a.Auto.Target.AddDevice(s)
		return nil

	case "antha.mixer.v1.Mixer":
		return a.AddMixer(ctx, conn, arg, reply.GetSubtypes())

	default:
		return nil
	}
}

var mixerMap = map[string]func(*tryer, context.Context, *grpc.ClientConn, interface{}) error{
	"GilsonPipetmax": (*tryer).addLowLevelMixer,
	"CyBio":          (*tryer).addLowLevelMixer,
	"TecanEvo":       (*tryer).addLowLevelMixer,
	"LabCyteEcho":    (*tryer).addHighLevelMixer,
	"Hamilton":       (*tryer).addLowLevelMixer,
}

// AddMixer queries a mixer driver and adds the corresponding device to the target
func (a *tryer) AddMixer(ctx context.Context, conn *grpc.ClientConn, arg interface{}, subtypes []string) error {
	if len(subtypes) == 0 {
		return errors.New("Cannot add mixer: no subtypes provided")
	} else if fun, found := mixerMap[subtypes[0]]; !found {
		return fmt.Errorf("Unknown mixer device: %v", subtypes)
	} else {
		return fun(a, ctx, conn, arg)
	}
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
