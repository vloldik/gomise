package gomise_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vloldik/gomise"
)

type PromiseContextData struct {
	IntData int
}

func TestPromiseGroupError(t *testing.T) {
	ctx := context.Background()
	group := gomise.NewPromiseGroup[int](
		gomise.NewDataContext[*PromiseContextData],
		func(ctx *gomise.DataContext[*PromiseContextData], resolve func(...any), reject func(error)) {
			ctx.Data = &PromiseContextData{IntData: 1}
			resolve()
		},
		func(ctx *gomise.DataContext[*PromiseContextData], resolve func(...any), reject func(error)) {
			println(ctx.Data.IntData)
			resolve()
		},
		func(ctx *gomise.DataContext[*PromiseContextData], resolve func(...any), reject func(error)) {
			println(ctx.Data.IntData)
			reject(errors.New("error"))
		},
		func(ctx *gomise.DataContext[*PromiseContextData], resolve func(...any), reject func(error)) {
			panic("i should never be called")
		},
	)
	_, err := group.Execute(ctx).Await(ctx)
	assert.Error(t, err)
}

func TestPromiseGroupResolving(t *testing.T) {
	ctx := context.Background()
	resolveValue := 12341
	group := gomise.NewPromiseGroup[int](
		gomise.NewDataContext[*PromiseContextData],
		func(ctx *gomise.DataContext[*PromiseContextData], resolve func(...any), reject func(error)) {
			ctx.Data = &PromiseContextData{IntData: 123}
			resolve()
		},
		func(ctx *gomise.DataContext[*PromiseContextData], resolve func(...any), reject func(error)) {
			ctx.Data = &PromiseContextData{IntData: 124}
			resolve()
		},
	).Add(func(ctx *gomise.DataContext[*PromiseContextData], resolve func(...any), reject func(error)) {
		if ctx.Data.IntData == 124 {
			resolve(resolveValue)
		}
	})
	val, err := group.Execute(ctx).Await(ctx)
	assert.NoError(t, err)
	assert.Equal(t, resolveValue, val)
}

type Doer struct{}

func (d *Doer) Do() Doer {
	return *d
}

func TestPointerInterface(t *testing.T) {
	ctx := context.Background()
	group := gomise.NewPromiseGroup[any](
		gomise.NewDataContext[interface{ Do() Doer }],
	)
	group.Add(func(ctx *gomise.DataContext[interface{ Do() Doer }], resolve func(...any), reject func(error)) {
		ctx.Data = &Doer{}
		resolve()
	}).Add(func(ctx *gomise.DataContext[interface{ Do() Doer }], resolve func(...any), reject func(error)) {
		resolve(ctx.Data.Do())
	})

	value, err := group.Execute(ctx).Await(ctx)

	assert.NoError(t, err)
	assert.Equal(t, value, Doer{})
}
