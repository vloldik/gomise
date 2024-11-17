package gomise

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/vloldik/gomise/interfaces"
)

var ErrNotFulfilled = errors.New("neither reject nor resolve has been called")

type gomise[T any] struct {
	isDone    bool
	doneMutex *sync.Mutex
	dataChan  chan T
	errChan   chan error
}

func (p *gomise[T]) SetDone() {
	p.doneMutex.Lock()
	defer p.doneMutex.Unlock()
	p.isDone = true
}

func (p *gomise[T]) GetIsDone() bool {
	p.doneMutex.Lock()
	defer p.doneMutex.Unlock()
	return p.isDone
}

func (p *gomise[T]) Await(ctx context.Context) (val T, err error) {
	select {
	case <-ctx.Done():
		return val, ctx.Err()
	case val = <-p.dataChan:
		return
	case err = <-p.errChan:
		return
	}
}

func (p *gomise[T]) close() {
	close(p.dataChan)
	close(p.errChan)
}

func (p *gomise[T]) fulfill(datas ...any) {
	if p.GetIsDone() {
		return
	}
	p.SetDone()
	defer p.close()

	var data any = nil
	if len(datas) > 0 {
		data = datas[0]
	}

	switch typedData := data.(type) {
	case nil:
		var defaultVal T
		p.dataChan <- defaultVal
	case T:
		p.dataChan <- typedData
	default:
		p.errChan <- fmt.Errorf("type of data is %T, illegal for %T", typedData, p)
	}
}

func (p *gomise[T]) reject(err error) {
	if p.GetIsDone() {
		return
	}
	p.SetDone()
	defer p.close()
	p.errChan <- err
}

func (p *gomise[T]) rejectIfNotFulfilled() {
	if !p.GetIsDone() {
		p.reject(ErrNotFulfilled)
	}
}

func makePromise[T any]() *gomise[T] {
	return &gomise[T]{
		doneMutex: new(sync.Mutex),
		dataChan:  make(chan T),
		errChan:   make(chan error),
	}
}

func NewPromise[R any](ctx context.Context, executable interfaces.FnPromiseExecutable[interfaces.IPromiseContext]) interfaces.IPromise[R] {
	return NewPromiseWithConstructor[R, interfaces.IPromiseContext](ctx, NewPromiseContext, executable)
}

func NewPromiseWithConstructor[R any, C interfaces.IPromiseContext](ctx context.Context, constructor interfaces.FnContextConstructor[C], executable interfaces.FnPromiseExecutable[C]) interfaces.IPromise[R] {
	gomise := makePromise[R]()
	promiseContext := constructor(context.WithCancel(ctx))
	go func() {
		executable(promiseContext, gomise.fulfill, gomise.reject)
		gomise.rejectIfNotFulfilled()
	}()
	return gomise
}

func NewPromiseFromContext[R any, C interfaces.IPromiseContext](ctx C, executable interfaces.FnPromiseExecutable[C]) interfaces.IPromise[R] {
	gomise := makePromise[R]()
	go func() {
		executable(ctx, gomise.fulfill, gomise.reject)
		gomise.rejectIfNotFulfilled()
	}()
	return gomise
}
