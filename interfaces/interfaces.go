package interfaces

import "context"

type IPromiseContext interface {
	// Context of current promise (for promise.Await another context used)
	Context() context.Context
	// Context cancel func
	Cancel()
}

type FnPromiseExecutable[C IPromiseContext] func(ctx C, resolve FnResolve, reject FnReject)
type FnResolve = func(...any)
type FnReject = func(error)
type FnContextConstructor[T IPromiseContext] func(context.Context, context.CancelFunc) T

// Base promise interface. It gets context and returns resolved value or error if something went wrong
type IPromise[T any] interface {
	Await(context.Context) (T, error)
}
