/*
Package gomise implements a type-safe Promise pattern for Go with context support.

The package provides a Promise implementation that allows for asynchronous operations
with proper error handling and cancellation support through Go's context package.

Basic usage:

	promise := gomise.NewPromise[string](context.Background(), func(ctx interfaces.IPromiseContext, resolve interfaces.FnResolve, reject interfaces.FnReject) {
	    // Async work
	    resolve("result")
	})

	result, err := promise.Await(context.Background())

Key Types and Interfaces:

  - IPromise[T] represents a Promise that will eventually contain a value of type T or an error
  - IPromiseContext provides context management with cancellation support
  - FnPromiseExecutable[C] is the function type for Promise executors
  - PromiseGroup allows sequential execution of multiple promises

Promise Creation:

The package provides several ways to create promises:

  - NewPromise creates a basic promise with default context
  - NewPromiseWithConstructor creates a promise with custom context type
  - NewPromiseFromContext creates a promise from an existing context

Promise Groups:

PromiseGroup allows sequential execution of multiple promises:

	group := gomise.NewDefaultPromiseGroup[string](
	    func(ctx interfaces.IPromiseContext, resolve interfaces.FnResolve, reject interfaces.FnReject) {
	        resolve("Step 1")
	    },
	    func(ctx interfaces.IPromiseContext, resolve interfaces.FnResolve, reject interfaces.FnReject) {
	        resolve("Step 2")
	    },
	)

	promise := group.Execute(context.Background())
	result, err := promise.Await(context.Background())

Data Context:

The package supports custom data contexts for passing data between promise stages:

	type MyData struct {
	    Value string
	}

	promise := gomise.NewPromiseWithConstructor[string, *gomise.DataContext[MyData]](
	    context.Background(),
	    gomise.NewDataContext[MyData],
	    func(ctx *gomise.DataContext[MyData], resolve interfaces.FnResolve, reject interfaces.FnReject) {
	        ctx.Data = MyData{Value: "Hello"}
	        resolve(ctx.Data.Value)
	    },
	)

Error Handling:

Promises can be rejected with errors using the reject function:

	reject(errors.New("operation failed"))

Context cancellation is also properly handled:

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Promise will be cancelled if not completed within timeout
	result, err := promise.Await(ctx)
*/
package gomise
