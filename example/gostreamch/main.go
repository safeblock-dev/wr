package main

import (
	"errors"
	"log"

	"github.com/safeblock-dev/wr/gostream"
	"github.com/safeblock-dev/wr/gostreamch"
)

func main() {
	// Create a new StreamCh instance.
	stream := gostreamch.New()
	defer stream.Wait() // Ensure all tasks are completed before exiting.

	// Submit task 1 to the stream.
	stream.Go(func() (gostream.Callback, error) {
		// Return a callback function for task 1.
		return func() error {
			log.Println("success task 1")
			return nil
		}, nil // No error for task 1.
	})

	// Submit task 2 to the stream.
	stream.Go(func() (gostream.Callback, error) {
		// Return a callback function for task 2.
		return func() error {
			log.Println("success task 2")
			return nil
		}, nil // No error for task 2.
	})

	// Submit task 3 to the stream, which intentionally returns an error.
	stream.Go(func() (gostream.Callback, error) {
		// Return nil for task 3 to indicate no callback function.
		return nil, errors.New("example error")
	})

	// Submit task 4 to the stream (will not be executed due to task 3 error).
	stream.Go(func() (gostream.Callback, error) {
		// This function will not be executed because task 3 returned an error.
		return func() error {
			log.Println("will not be executed")
			return nil
		}, nil // No error for task 4.
	})

	// Read and log the error from the stream's error channel.
	log.Println("error:", <-stream.ErrorChannel())
}
