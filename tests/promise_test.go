package gomise_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vloldik/gomise"
)

func TestContextCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	promise := gomise.NewPromise[string](ctx, func(ctx gomise.IPromiseContext, resolve gomise.FnResolve, reject gomise.FnReject) {
		time.Sleep(2 * time.Second)
		resolve("Too late")
	})

	_, err := promise.Await(ctx)

	assert.ErrorIs(t, err, context.DeadlineExceeded)
}
