# Gomise - Go Promise Implementation

Gomise is a type-safe Promise implementation for Go that provides asynchronous programming capabilities with context support. It offers features similar to JavaScript Promises but with Go's type safety and context-aware concurrency patterns.

## Features

- Generic type-safe promises
- Context-aware promise execution
- Promise groups for sequential execution
- Custom context support with data passing
- Cancellation support
- Error handling with type safety

## Installation

```bash
go get github.com/vloldik/gomise
```

## Usage

### Basic Promise

```go
import (
    "context"
    "github.com/vloldik/gomise"
)

// Create a simple promise that resolves with a string
promise := gomise.NewPromise[string](context.Background(), func(ctx interfaces.IPromiseContext, resolve interfaces.FnResolve, reject interfaces.FnReject) {
    // Async work here
    resolve("Hello, World!")
})

// Await the result
result, err := promise.Await(context.Background())
if err != nil {
    // Handle error
}
fmt.Println(result) // "Hello, World!"
```

### Promise with Data Context

```go
type MyData struct {
    Value string
}

// Create a promise with custom data context
promise := gomise.NewPromiseWithConstructor[string, *gomise.DataContext[MyData]](
    context.Background(),
    gomise.NewDataContext[MyData],
    func(ctx *gomise.DataContext[MyData], resolve interfaces.FnResolve, reject interfaces.FnReject) {
        ctx.Data = MyData{Value: "Hello"}
        resolve(ctx.Data.Value)
    },
)
```

### Promise Groups

```go
// Create a promise group for sequential execution
group := gomise.NewDefaultPromiseGroup[string](
    func(ctx interfaces.IPromiseContext, resolve interfaces.FnResolve, reject interfaces.FnReject) {
        resolve("Step 1")
    },
    func(ctx interfaces.IPromiseContext, resolve interfaces.FnResolve, reject interfaces.FnReject) {
        resolve("Step 2")
    },
)

// Execute the group
promise := group.Execute(context.Background())
result, err := promise.Await(context.Background())
```

## Error Handling

Promises can be rejected with errors:

```go
promise := gomise.NewPromise[string](context.Background(), func(ctx interfaces.IPromiseContext, resolve interfaces.FnResolve, reject interfaces.FnReject) {
    reject(errors.New("something went wrong"))
})

result, err := promise.Await(context.Background())
if err != nil {
    // Handle error
}
```

## Context Cancellation

Promises respect context cancellation:

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()

promise := gomise.NewPromise[string](ctx, func(ctx interfaces.IPromiseContext, resolve interfaces.FnResolve, reject interfaces.FnReject) {
    time.Sleep(2 * time.Second) // Will be cancelled
    resolve("Too late")
})

_, err := promise.Await(ctx)
// err will be context.DeadlineExceeded
```

## License

This project is open source and available under the [MIT License](LICENSE).
