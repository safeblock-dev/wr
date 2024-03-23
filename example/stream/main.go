package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/safeblock-dev/wr/panics"
	"github.com/safeblock-dev/wr/wrstream"
)

const maxGoroutines = 3

func main() {
	ctx, cancel := context.WithCancel(context.Background()) // Create a cancellable context.
	defer cancel()                                          // Ensure the context is cancelled when the function exits.

	example(ctx) // Run the example function.
}

func example(ctx context.Context) {
	// Create a new worker pool with a context, maximum number of goroutines, and custom panic and error handlers.
	stream := wrstream.New(
		wrstream.Context(ctx),
		wrstream.MaxGoroutines(maxGoroutines),
		wrstream.PanicHandler(panicHandler),
		wrstream.ErrorHandler(errorHandler),
	)
	defer stream.Wait() // Ensure all tasks are completed before exiting.

	for i := 0; i < 100; i++ {
		stream.Go(func() (wrstream.Callback, error) {
			time.Sleep(time.Duration(rand.Uint32()) / 2)

			switch {
			case i%7 == 0:
				return nil, errors.New(strconv.Itoa(i))
			case i%11 == 0:
				panic(strconv.Itoa(i))
			}

			return func() error {
				switch {
				case i%3 == 0:
					return errors.New(strconv.Itoa(i))
				case i%8 == 0:
					panic(strconv.Itoa(i))
				}

				log.Println("success:\t", i)

				return nil
			}, nil
		})
	}
}

// Custom panic handler that logs the panic value.
func panicHandler(p panics.Recovered) {
	log.Println("panic:\t", p.Value)
}

// Custom error handler that logs the error.
func errorHandler(err error) {
	log.Println("error:\t", err)
}
