package target

import (
	"fmt"
	"strings"
	"time"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	runner "github.com/antha-lang/antha/driver/antha_runner_v1"
)

type Runner struct {
	client runner.RunnerClient
}

func NewRunner(c runner.RunnerClient) *Runner {
	return &Runner{c}
}

func (a *Runner) Run(files Files) error {
	ctx := context.Background()
	reply, err := a.client.Run(ctx, &runner.RunRequest{
		Type: files.Type,
		Data: files.Tarball,
	})
	if err != nil {
		return err
	}

	// Proof of concept
	var errors []string
	for range time.Tick(5 * time.Second) {
		msgs, err := a.client.Messages(ctx, &runner.MessagesRequest{
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

func (a *Runner) types() ([]string, error) {
	ctx := context.Background()
	reply, err := a.client.SupportedRunTypes(ctx, &runner.SupportedRunTypesRequest{})
	if err != nil {
		return nil, err
	}
	return reply.Types, nil
}

func (a *Runner) String() string {
	return "Runner"
}

func (a *Runner) CanCompile(ast.Request) bool {
	return false
}

func (a *Runner) MoveCost(Device) int {
	return 0
}

func (a *Runner) Compile([]ast.Node) ([]Inst, error) {
	return nil, nil
}
