package gomise

import (
	"context"

	"github.com/vloldik/gomise/interfaces"
)

type DataContext[D any] struct {
	interfaces.IPromiseContext
	Data D
}

func NewDataContext[D any](ctx context.Context, cancel context.CancelFunc) *DataContext[D] {
	return &DataContext[D]{
		IPromiseContext: NewPromiseContext(ctx, cancel),
	}
}
