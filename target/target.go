// Package target provides the construction of a target machine from a
// collection of devices
package target

import (
	"errors"
	"fmt"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
)

var (
	noLh         = errors.New("no liquid handler found")
	noTarget     = errors.New("no target configuration found")
	alreadyAdded = errors.New("already added")
)

type targetKey int

const theTargetKey targetKey = 0

func GetTarget(ctx context.Context) (*Target, error) {
	v, ok := ctx.Value(theTargetKey).(*Target)
	if !ok {
		return nil, noTarget
	}
	return v, nil
}

func WithTarget(parent context.Context, t *Target) context.Context {
	return context.WithValue(parent, theTargetKey, t)
}

// Target machine for execution.
type Target struct {
	devices []Device
}

func New() *Target {
	return &Target{}
}

func (a *Target) String() string {
	var r []string
	for _, d := range a.devices {
		r = append(r, fmt.Sprint(d))
	}
	return fmt.Sprint(r)
}

func (a *Target) canCompile(d Device, reqs ...ast.Request) bool {
	for _, req := range reqs {
		if !d.CanCompile(req) {
			return false
		}
	}
	return true
}

func (a *Target) CanCompile(reqs ...ast.Request) (r []Device) {
	for _, d := range a.devices {
		if a.canCompile(d, reqs...) {
			r = append(r, d)
		}
	}
	return
}

func (a *Target) AddDevice(d Device) error {
	a.devices = append(a.devices, d)
	return nil
}
