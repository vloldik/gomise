package interfaces

import "context"

type IPromiseContext interface {
	Context() context.Context
	Cancel()
}

type FnPromiseExecutable[C IPromiseContext] func(ctx C, resolve FnResolve, reject FnReject)
type FnResolve = func(...any)
type FnReject = func(error)
type FnContextConstructor[T IPromiseContext] func(context.Context, context.CancelFunc) T

type IPromise[T any] interface {
	Await(context.Context) (T, error)
}
