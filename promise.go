package gomise

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/vloldik/gomise/interfaces"
)

// ErrNotFulfilled is returned when a promise is not resolved or rejected
var ErrNotFulfilled = errors.New("neither reject nor resolve has been called")

// gomise is a generic promise implementation
type gomise[T any] struct {
	isDone    bool        // indicates whether promise is completed
	doneMutex *sync.Mutex // mutex for thread-safe done state
	dataChan  chan T      // channel for successful result
	errChan   chan error  // channel for error result
}

// SetDone marks the promise as completed in a thread-safe manner
func (p *gomise[T]) SetDone() {
	p.doneMutex.Lock()
	defer p.doneMutex.Unlock()
	p.isDone = true
}

// GetIsDone checks if the promise is completed in a thread-safe manner
func (p *gomise[T]) GetIsDone() bool {
	p.doneMutex.Lock()
	defer p.doneMutex.Unlock()
	return p.isDone
}

// Await waits for promise resolution with context support
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

// close terminates the promise's channels
func (p *gomise[T]) close() {
	close(p.dataChan)
	close(p.errChan)
}

// fulfill completes the promise with a successful value
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

	// Type checking and value sending
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

// reject completes the promise with an error
func (p *gomise[T]) reject(err error) {
	if p.GetIsDone() {
		return
	}
	p.SetDone()
	defer p.close()
	p.errChan <- err
}

// rejectIfNotFulfilled ensures the promise is resolved if not already done
func (p *gomise[T]) rejectIfNotFulfilled() {
	if !p.GetIsDone() {
		p.reject(ErrNotFulfilled)
	}
}

// makePromise creates a new promise instance
func makePromise[T any]() *gomise[T] {
	return &gomise[T]{
		doneMutex: new(sync.Mutex),
		dataChan:  make(chan T),
		errChan:   make(chan error),
	}
}

// NewPromise creates a promise with default context
func NewPromise[R any](ctx context.Context, executable interfaces.FnPromiseExecutable[interfaces.IPromiseContext]) interfaces.IPromise[R] {
	return NewPromiseWithConstructor[R, interfaces.IPromiseContext](ctx, NewPromiseContext, executable)
}

// NewPromiseWithConstructor creates a promise with a custom context constructor
func NewPromiseWithConstructor[R any, C interfaces.IPromiseContext](ctx context.Context, constructor interfaces.FnContextConstructor[C], executable interfaces.FnPromiseExecutable[C]) interfaces.IPromise[R] {
	gomise := makePromise[R]()
	promiseContext := constructor(context.WithCancel(ctx))
	go func() {
		executable(promiseContext, gomise.fulfill, gomise.reject)
		gomise.rejectIfNotFulfilled()
	}()
	return gomise
}

// NewPromiseFromContext creates a promise from an existing context
func NewPromiseFromContext[R any, C interfaces.IPromiseContext](ctx C, executable interfaces.FnPromiseExecutable[C]) interfaces.IPromise[R] {
	gomise := makePromise[R]()
	go func() {
		executable(ctx, gomise.fulfill, gomise.reject)
		gomise.rejectIfNotFulfilled()
	}()
	return gomise
}
