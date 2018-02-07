package inject

import (
	"context"
	"errors"
	"sync"

	api "github.com/antha-lang/antha/api/v1"
)

var errAlreadyAdded = errors.New("already added")

type registry struct {
	lock   sync.Mutex
	parent context.Context
	reg    map[Name]Runner
}

// Name uniquely identifiers a inject.Runner
type Name struct {
	Host  string // Host
	Repo  string // Name
	Tag   string // Version
	Stage api.ElementStage
}

// NameQuery is a query for a Runner
type NameQuery struct {
	Repo  string // Name
	Tag   string // Version
	Stage api.ElementStage
}

func (a *registry) Add(name Name, runner Runner) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.reg == nil {
		a.reg = make(map[Name]Runner)
	}
	if r := a.reg[name]; r != nil {
		return errAlreadyAdded
	}
	a.reg[name] = runner
	return nil
}

func (a *registry) Find(query NameQuery) ([]Runner, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	name := Name{
		Repo:  query.Repo,
		Tag:   query.Tag,
		Stage: query.Stage,
	}
	r := a.reg[name]
	if r == nil {
		return nil, nil
	}
	return []Runner{r}, nil
}
