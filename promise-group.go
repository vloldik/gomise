package gomise

import (
	"context"

	"github.com/vloldik/gomise/interfaces"
)

type PromiseGroup[R any, C interfaces.IPromiseContext] struct {
	contextConstructor interfaces.FnContextConstructor[C]
	fns                []interfaces.FnPromiseExecutable[C]
}

func NewPromiseGroup[R any, C interfaces.IPromiseContext](contextConstructor interfaces.FnContextConstructor[C], fns ...interfaces.FnPromiseExecutable[C]) *PromiseGroup[R, C] {
	return &PromiseGroup[R, C]{
		contextConstructor: contextConstructor,
		fns:                fns,
	}
}

func NewDefaultPromiseGroup[R any](fns ...interfaces.FnPromiseExecutable[interfaces.IPromiseContext]) *PromiseGroup[R, interfaces.IPromiseContext] {
	return &PromiseGroup[R, interfaces.IPromiseContext]{
		contextConstructor: NewPromiseContext,
		fns:                fns,
	}
}

func (pg *PromiseGroup[R, C]) Add(fn interfaces.FnPromiseExecutable[C]) *PromiseGroup[R, C] {
	pg.fns = append(pg.fns, fn)
	return pg
}

func (pg *PromiseGroup[R, C]) Execute(ctx context.Context) interfaces.IPromise[R] {
	return NewPromiseWithConstructor[R, C](ctx, pg.contextConstructor, func(ctx C, resolve interfaces.FnResolve, reject interfaces.FnReject) {
		for i, fn := range pg.fns {
			pg.executeOne(ctx, i == len(pg.fns)-1, fn, resolve, reject)
		}
	})
}

func (pg *PromiseGroup[R, C]) executeOne(ctx C, isLast bool, fn interfaces.FnPromiseExecutable[C], resolve interfaces.FnResolve, reject interfaces.FnReject) {
	select {
	case <-ctx.Context().Done():
		reject(ctx.Context().Err())
	default:
		if isLast {
			innerPromise := NewPromiseFromContext[R](ctx, fn)
			data, err := innerPromise.Await(ctx.Context())
			if err != nil {
				reject(err)
				return
			}
			resolve(data)
		} else {
			innerPromise := NewPromiseFromContext[any](ctx, fn)
			if _, err := innerPromise.Await(ctx.Context()); err != nil {
				reject(err)
				ctx.Cancel()
				return
			}
		}
	}
}
