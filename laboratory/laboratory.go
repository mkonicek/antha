package laboratory

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Element interface {
	Setup(*Laboratory)
	Steps(*Laboratory)
	Analysis(*Laboratory)
	Validation(*Laboratory)
}

type Laboratory struct {
	elemLock   sync.Mutex
	elemsUnrun int64
	elements   map[Element]*ElementBase

	errLock sync.Mutex
	errors  []error
	Errored chan struct{}

	Completed chan struct{}
}

func NewLaboratory() *Laboratory {
	return &Laboratory{
		elements:  make(map[Element]*ElementBase),
		Errored:   make(chan struct{}),
		Completed: make(chan struct{}),
	}
}

// Only use this before you call run.
func (lab *Laboratory) InstallElement(e Element) {
	eb := NewElementBase(e)
	lab.elemLock.Lock()
	defer lab.elemLock.Unlock()
	lab.elements[e] = eb
	lab.elemsUnrun++
}

func (lab *Laboratory) AddLink(src, dst Element, fun func()) error {
	if ebSrc, found := lab.elements[src]; !found {
		return fmt.Errorf("Unknown src element: %v", src)
	} else if ebDst, found := lab.elements[dst]; !found {
		return fmt.Errorf("Unknown dst element: %v", dst)
	} else {
		ebDst.AddBlockedInput()
		ebSrc.AddOnExit(func() { fun(); ebDst.InputReady() })
		return nil
	}
}

// Run all the installed elements.
func (lab *Laboratory) Run() error {
	lab.elemLock.Lock()
	if lab.elemsUnrun == 0 {
		lab.elemLock.Unlock()
		close(lab.Completed)
		return nil

	} else {
		for _, eb := range lab.elements {
			eb.AddOnExit(lab.elementCompleted)
		}
		for _, eb := range lab.elements {
			go eb.Run(lab)
		}
		<-lab.Completed

		select {
		case <-lab.Errored:
			return lab.errors[0]
		default:
			return nil
		}
	}
}

// Only for use when you're in an element and want to call another element.
func (lab *Laboratory) CallSteps(e Element) error {
	eb := NewElementBase(e)

	finished := make(chan struct{})
	eb.AddOnExit(func() { close(finished) })

	go eb.Run(lab, eb.element.Steps)
	<-finished

	select {
	case <-lab.Errored:
		return lab.errors[0]
	default:
		return nil
	}
}

func (lab *Laboratory) elementCompleted() {
	lab.elemLock.Lock()
	defer lab.elemLock.Unlock()
	lab.elemsUnrun--
	if lab.elemsUnrun == 0 {
		close(lab.Completed)
	}
}

func (lab *Laboratory) Error(err error) {
	lab.errLock.Lock()
	defer lab.errLock.Unlock()
	lab.errors = append(lab.errors, err)
	select { // we keep the lock here to avoid a race to close
	case <-lab.Errored:
	default:
		close(lab.Errored)
	}
}

func (lab *Laboratory) Errorf(fmtStr string, args ...interface{}) {
	err := fmt.Errorf(fmtStr, args...)
	lab.Error(err)
}

// ElementBase
type ElementBase struct {
	// count of inputs that are not yet ready (plus 1)
	pendingCount int64
	// this gets closed when all inputs become ready
	InputsReady chan struct{}
	// funcs to run when this element is completed
	onExit []func()
	// the actual element
	element Element
}

func NewElementBase(e Element) *ElementBase {
	return &ElementBase{
		pendingCount: 1,
		InputsReady:  make(chan struct{}),
		element:      e,
	}
}

func (eb *ElementBase) Run(lab *Laboratory, funs ...func(*Laboratory)) {
	defer eb.Completed()
	eb.InputReady()

	if len(funs) == 0 {
		funs = []func(*Laboratory){
			eb.element.Setup,
			eb.element.Steps,
			eb.element.Analysis,
			eb.element.Validation,
		}
	}

	select {
	case <-eb.InputsReady:
		for _, fun := range funs {
			select {
			case <-lab.Errored:
				return
			default:
				fun(lab)
			}
		}
	case <-lab.Errored:
		return
	}
}

func (eb *ElementBase) Completed() {
	funs := eb.onExit
	eb.onExit = nil
	for _, fun := range funs {
		fun()
	}
}

func (eb *ElementBase) InputReady() {
	if atomic.AddInt64(&eb.pendingCount, -1) == 0 {
		// we've done the transition from 1 -> 0. By definition, we're
		// the only routine that can do that, so we don't need to be
		// careful.
		close(eb.InputsReady)
	}
}

func (eb *ElementBase) AddBlockedInput() {
	atomic.AddInt64(&eb.pendingCount, 1)
}

func (eb *ElementBase) AddOnExit(fun func()) {
	eb.onExit = append(eb.onExit, fun)
}
