package inject

import (
	"context"
	"errors"
	"sync"
)

var alreadyAdded = errors.New("already added")

type registry struct {
	lock   sync.Mutex
	parent context.Context
	reg    map[Name]Runner
	// TODO: add remote registries...
}

// Unique identifier for Runner
type Name struct {
	Host string // Host
	Repo string // Name
	Tag  string // Version
}

// Query for a Runner
type NameQuery struct {
	Repo string // Name
	Tag  string // Version
}

func (a *registry) Add(name Name, runner Runner) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.reg == nil {
		a.reg = make(map[Name]Runner)
	}
	if r := a.reg[name]; r != nil {
		return alreadyAdded
	}
	a.reg[name] = runner
	return nil
}

func (a *registry) Find(query NameQuery) ([]Runner, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	name := Name{Repo: query.Repo, Tag: query.Tag}
	if r := a.reg[name]; r == nil {
		return nil, nil
	} else {
		return []Runner{r}, nil
	}
}
