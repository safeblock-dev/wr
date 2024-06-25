package main

import (
	"errors"
	"log"

	"github.com/safeblock-dev/wr/gopool"
)

const maxGoroutines = 2

func main() {
	// Create a new worker pool with a context, maximum number of goroutines, and custom panic and error handlers.
	pool := gopool.New(
		gopool.MaxGoroutines(maxGoroutines), // Limit the pool to a maximum of maxGoroutines concurrent goroutines.
		gopool.PanicHandler(panicHandler),   // Set custom panic handler to log panics.
		gopool.ErrorHandler(errorHandler),   // Set custom error handler to log errors.
	)
	defer pool.Wait() // Ensure all tasks are completed before exiting by waiting for the pool.

	// Submit a task to the pool that logs a value.
	pool.Go(func() error {
		log.Println("value:", "foo")
		return nil // Return nil to indicate no error.
	})

	// Submit a task to the pool that returns an error.
	pool.Go(func() error {
		return errors.New("my error") // Return an error to simulate an error condition.
	})

	// Submit a task to the pool that panics.
	pool.Go(func() error {
		panic("my panic") // Intentionally panic to demonstrate panic handling.
	})
}

// Custom panic handler that logs the panic value.
func panicHandler(pc any) {
	log.Println("panic:", pc)
}

// Custom error handler that logs the error.
func errorHandler(err error) {
	log.Println("error:", err)
}
