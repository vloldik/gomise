package gomise

import (
	"context"
)

// PromiseGroup represents a group of executable promises with shared context
// R is the result type, C is the context type
type PromiseGroup[R any, C IPromiseContext] struct {
	// Function to construct the promise context
	contextConstructor FnContextConstructor[C]
	// List of executable promise functions
	fns []FnPromiseExecutable[C]
}

// NewPromiseGroup creates a new promise group with custom context and functions
// Allows specifying a custom context constructor and initial set of executable functions
func NewPromiseGroup[R any, C IPromiseContext](contextConstructor FnContextConstructor[C], fns ...FnPromiseExecutable[C]) *PromiseGroup[R, C] {
	return &PromiseGroup[R, C]{
		contextConstructor: contextConstructor,
		fns:                fns,
	}
}

// NewDefaultPromiseGroup creates a promise group with default context
// Uses the standard promise context and allows adding executable functions
func NewDefaultPromiseGroup[R any](fns ...FnPromiseExecutable[IPromiseContext]) *PromiseGroup[R, IPromiseContext] {
	return &PromiseGroup[R, IPromiseContext]{
		contextConstructor: NewPromiseContext,
		fns:                fns,
	}
}

// Add appends a new executable function to the promise group
// Allows chaining additional functions dynamically
func (pg *PromiseGroup[R, C]) Add(fn FnPromiseExecutable[C]) *PromiseGroup[R, C] {
	pg.fns = append(pg.fns, fn)
	return pg
}

// Execute runs the promise group with the given context
// Executes all functions sequentially, with the last function determining the final result
func (pg *PromiseGroup[R, C]) Execute(ctx context.Context) IPromise[R] {
	return NewPromiseWithConstructor[R, C](ctx, pg.contextConstructor, func(ctx C, resolve FnResolve, reject FnReject) {
		for i, fn := range pg.fns {
			pg.executeOne(ctx, i == len(pg.fns)-1, fn, resolve, reject)
		}
	})
}

// executeOne handles the execution of a single promise function
// Manages context cancellation, error handling, and result resolution
func (pg *PromiseGroup[R, C]) executeOne(ctx C, isLast bool, fn FnPromiseExecutable[C], resolve FnResolve, reject FnReject) {
	select {
	case <-ctx.Context().Done():
		// Reject if context is cancelled
		reject(ctx.Context().Err())
	default:
		if isLast {
			// For the last function, resolve with its result
			innerPromise := NewPromiseFromContext[R](ctx, fn)
			data, err := innerPromise.Await(ctx.Context())
			if err != nil {
				reject(err)
				return
			}
			resolve(data)
		} else {
			// For intermediate functions, just execute without resolving
			innerPromise := NewPromiseFromContext[any](ctx, fn)
			if _, err := innerPromise.Await(ctx.Context()); err != nil {
				reject(err)
				ctx.Cancel()
				return
			}
		}
	}
}
