package main

import (
	"context"
	"errors"
	"log"

	"github.com/safeblock-dev/wr/panics"
	"github.com/safeblock-dev/wr/wrpool"
)

const maxGoroutines = 2

// result:
// 2024/03/20 23:17:54 error: my error
// 2024/03/20 23:17:54 value: foo
// 2024/03/20 23:17:54 panic: my panic

func main() {
	ctx, cancel := context.WithCancel(context.Background()) // Create a cancellable context.
	defer cancel()                                          // Ensure the context is cancelled when the function exits.

	example(ctx) // Run the example function.
}

func example(ctx context.Context) {
	// Create a new worker pool with a context, maximum number of goroutines, and custom panic and error handlers.
	pool := wrpool.New(
		wrpool.Context(ctx),
		wrpool.MaxGoroutines(maxGoroutines),
		wrpool.PanicHandler(panicHandler),
		wrpool.ErrorHandler(errorHandler),
	)
	defer pool.Wait() // Ensure all tasks are completed before exiting.

	// Submit a task to the pool that logs a value.
	pool.Go(func() error {
		log.Println("value:", "foo")
		return nil
	})

	// Submit a task to the pool that returns an error.
	pool.Go(func() error {
		return errors.New("my error")
	})

	// Submit a task to the pool that panics.
	pool.Go(func() error {
		panic("my panic")
	})
}

// Custom panic handler that logs the panic value.
func panicHandler(p panics.Recovered) {
	log.Println("panic:", p.Value)
}

// Custom error handler that logs the error.
func errorHandler(err error) {
	log.Println("error:", err)
}
