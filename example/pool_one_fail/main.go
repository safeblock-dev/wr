package main

import (
	"context"
	"errors"
	"log"

	"github.com/safeblock-dev/wr/wrpool"
)

const (
	maxGoroutines = 2
	dataSize      = 100
)

// result:
// 2024/03/21 07:02:17 value: 1
// 2024/03/21 07:02:17 value: 2
// 2024/03/21 07:02:17 value: 3
// 2024/03/21 07:02:17 value: 0
// 2024/03/21 07:02:17 value: 5
// 2024/03/21 07:02:17 value: 4
// 2024/03/21 07:02:17 error: my error

func main() {
	ctx, cancel := context.WithCancel(context.Background()) // Create a cancellable context.
	defer cancel()                                          // Ensure the context is cancelled when the function exits.

	// Create a new worker pool with a context, maximum number of goroutines, and custom panic and error handlers.
	pool := wrpool.New(
		wrpool.Context(ctx),
		wrpool.MaxGoroutines(maxGoroutines),
		wrpool.ErrorHandler(func(err error) {
			log.Printf("error: %v", err)
			cancel()
		}),
	)
	defer pool.Wait() // Ensure all tasks are completed before exiting.

	for i := 0; i < dataSize; i++ {
		pool.Go(func() error {
			log.Printf("value: %d", i)

			if i == 5 {
				return errors.New("my error")
			}

			return nil
		})
	}
}
