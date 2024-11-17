package gomise

import (
	"context"

	"github.com/vloldik/gomise/interfaces"
)

type promiseContext struct {
	context context.Context
	cancel  context.CancelFunc
}

func NewPromiseContext(ctx context.Context, cancel context.CancelFunc) interfaces.IPromiseContext {
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
