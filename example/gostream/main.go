package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/safeblock-dev/wr/gostream"
)

const maxGoroutines = 3

func main() {
	ctx, cancel := context.WithCancel(context.Background()) // Create a cancellable context.
	defer cancel()                                          // Ensure the context is cancelled when the function exits.

	// Create a new stream with a context, maximum number of goroutines, and custom panic and error handlers.
	stream := gostream.New(
		gostream.Context(ctx),
		gostream.MaxGoroutines(maxGoroutines),
		gostream.PanicHandler(panicHandler),
		gostream.ErrorHandler(errorHandler),
	)
	defer stream.Wait() // Ensure all tasks are completed before exiting.

	for i := 0; i < 100; i++ {
		stream.Go(func() (gostream.Callback, error) {
			// Simulate variable task execution time
			time.Sleep(time.Duration(rand.Uint32()) / 2)

			switch {
			case i%7 == 0:
				// Return an error for every 7th task
				return nil, errors.New(strconv.Itoa(i))
			case i%11 == 0:
				// Simulate a panic for every 11th task
				panic(strconv.Itoa(i))
			}

			// Return a callback function for successful tasks
			return func() error {
				switch {
				case i%3 == 0:
					// Return an error for every 3rd successful task
					return errors.New(strconv.Itoa(i))
				case i%8 == 0:
					// Simulate a panic for every 8th successful task
					panic(strconv.Itoa(i))
				}

				// Log success for tasks that don't error or panic
				log.Println("success:\t", i)

				return nil
			}, nil
		})
	}
}

// Custom panic handler that logs the panic value.
func panicHandler(pc any) {
	log.Println("panic:\t", pc)
}

// Custom error handler that logs the error.
func errorHandler(err error) {
	log.Println("error:\t", err)
}
