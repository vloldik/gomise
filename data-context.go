package gomise

import (
	"context"
)

// Used to share data between promises in promise group
type DataContext[D any] struct {
	IPromiseContext
	Data D
}

// Use it instead of default constrictor in promise or promisegroup
func NewDataContext[D any](ctx context.Context, cancel context.CancelFunc) *DataContext[D] {
	return &DataContext[D]{
		IPromiseContext: NewPromiseContext(ctx, cancel),
	}
}
