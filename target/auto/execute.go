package auto

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/antha-lang/antha/ast"
	runner "github.com/antha-lang/antha/driver/antha_runner_v1"
	"github.com/antha-lang/antha/target"
	"google.golang.org/grpc"
)

// Execute runs an instruction based on current target
func (a *Auto) Execute(ctx context.Context, inst ast.Inst) error {
	switch inst := inst.(type) {
	case *target.Mix:
		return a.executeMix(ctx, inst)
	case *target.Run:
		return a.executeRun(ctx, inst)
	case *target.Manual:
		return nil
	case *target.Wait:
		return nil
	case *target.Order:
		return nil
	case *target.PlatePrep:
		return nil
	case *target.SetupMixer:
		return nil
	case *target.Prompt:
		return nil
	case *target.TimedWait:
		return nil
	default:
		return fmt.Errorf("unknown instruction %T", inst)
	}
}

func (a *Auto) executeRun(ctx context.Context, inst *target.Run) error {
	conn, ok := a.handler[inst.Dev]
	if !ok {
		return fmt.Errorf("no handler for %s", inst.Label)
	}

	for _, c := range inst.Calls {
		if err := grpc.Invoke(ctx, c.Method, c.Args, c.Reply, conn); err != nil {
			return err
		}
	}
	return nil
}

func (a *Auto) executeMix(ctx context.Context, inst *target.Mix) error {
	rs := a.runners[inst.Files.Type]
	if len(rs) == 0 {
		return fmt.Errorf("no runner for %s", inst.Files.Type)
	}
	r := rs[0]
	reply, err := r.Run(ctx, &runner.RunRequest{
		Type: inst.Files.Type,
		Data: inst.Files.Tarball,
	})
	if err != nil {
		return err
	}

	// Proof of concept
	var errors []string
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		msgs, err := r.Messages(ctx, &runner.MessagesRequest{
			Id: reply.Id,
		})
		if err != nil {
			return err
		}
		for _, m := range msgs.Values {
			if m.Code == "error" {
				errors = append(errors, string(m.Data))
			}
			if m.Code == "fatal" {
				return fmt.Errorf("error running protocol: %s", strings.Join(errors, " "))
			}
			if m.Code == "stop" {
				return nil
			}
		}
	}

	return nil
}
