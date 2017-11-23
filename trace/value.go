package trace

// A Value is a wrapper around instruction inputs and outputs
type Value interface {
	Name() Name
	Get() interface{}
}

type basicValue struct {
	name Name
	v    interface{}
}

func (a *basicValue) Get() interface{} {
	return a.v
}

func (a *basicValue) Name() Name {
	return a.name
}

type fromValue struct {
	name Name
	v    interface{}
	from []Value
}

func (a *fromValue) Get() interface{} {
	return a.v
}

func (a *fromValue) Name() Name {
	return a.name
}

type promiseValue struct {
	name Name
	v    interface{}
	inst interface{}
}

func (a *promiseValue) Get() interface{} {
	return a.v
}

func (a *promiseValue) Name() Name {
	return a.name
}
