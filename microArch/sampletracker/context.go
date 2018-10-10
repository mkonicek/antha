package sampletracker

import "context"

type contextKey string

const theContextKey contextKey = "sampletracker"

// NewContext return a new context inheriting the parent with the sampletracker added
func NewContext(parent context.Context) context.Context {
	return context.WithValue(parent, theContextKey, NewSampleTracker())
}

// FromContext fetch the SampleTracker from the context, calls panic() if the SampleTracker has not been set
func FromContext(ctx context.Context) *SampleTracker {
	if st, ok := ctx.Value(theContextKey).(*SampleTracker); !ok {
		panic("no SampleTracker in context")
	} else {
		return st
	}
}
