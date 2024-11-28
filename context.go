package gomise

import (
	"context"
)

type promiseContext struct {
	context context.Context
	cancel  context.CancelFunc
}

// Returns basic implementation of
//
//	IPromiseContext
func NewPromiseContext(ctx context.Context, cancel context.CancelFunc) IPromiseContext {
	return &promiseContext{
		context: ctx,
		cancel:  cancel,
	}
}

func (p *promiseContext) Context() context.Context {
	return p.context
}

func (p *promiseContext) Cancel() {
	p.cancel()
}
