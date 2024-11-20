package gomise_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vloldik/gomise"
	"github.com/vloldik/gomise/interfaces"
)

func TestContextCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	promise := gomise.NewPromise[string](ctx, func(ctx interfaces.IPromiseContext, resolve interfaces.FnResolve, reject interfaces.FnReject) {
		time.Sleep(2 * time.Second)
		resolve("Too late")
	})

	_, err := promise.Await(ctx)

	assert.ErrorIs(t, err, context.DeadlineExceeded)
}
