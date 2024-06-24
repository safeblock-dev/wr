package main

import (
	"errors"
	"log"
	"strings"

	"github.com/safeblock-dev/wr/gopoolch"
)

func main() {
	// Create a new goroutine pool with callback-based error handling.
	pool := gopoolch.New()
	defer pool.Wait() // Ensure all tasks are completed before exiting.

	// Submit a task to the pool that logs a value.
	pool.Go(func() error {
		log.Println("value:", "foo")
		return nil // Return nil to indicate no error.
	})

	// Submit a task to the pool that returns an error.
	pool.Go(func() error {
		return errors.New("my error") // Return an error.
	})

	// Submit a task to the pool that panics.
	pool.Go(func() error {
		panic("my panic") // Simulate a panic.
	})

	// Retrieve the first error from the error channel.
	err := <-pool.ErrorChannel()
	// Log the error, ensuring only the first line is logged to avoid excessive output.
	log.Println("error:", strings.Split(err.Error(), "\n")[0])
}
