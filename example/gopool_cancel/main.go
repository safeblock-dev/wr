package main

import (
	"context"
	"errors"
	"log"

	"github.com/safeblock-dev/wr/gopool"
)

const (
	maxGoroutines = 2
	dataSize      = 100
)

func main() {
	// Create a cancellable context and defer cancellation to ensure resources are cleaned up.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new worker pool with the cancellable context, maximum number of goroutines,
	// and a custom error handler that logs errors and cancels the context on encountering an error.
	pool := gopool.New(
		gopool.Context(ctx),                 // Set the context for the worker pool.
		gopool.MaxGoroutines(maxGoroutines), // Limit the pool to maxGoroutines concurrent goroutines.
		gopool.ErrorHandler(func(err error) {
			log.Printf("error: %v", err) // Log the error encountered.
			cancel()                     // Cancel the context to stop further task execution.
		}),
	)
	defer pool.Wait() // Ensure all tasks are completed before exiting by waiting for the pool.

	// Submit dataSize tasks to the pool.
	for i := 0; i < dataSize; i++ {
		pool.Go(func() error {
			log.Printf("value: %d", i) // Log the current value of i.

			if i == 5 {
				return errors.New("my error") // Return an error for a specific value of i.
			}

			return nil // Return nil to indicate no error.
		})
	}
}
