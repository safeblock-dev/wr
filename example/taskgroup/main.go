package main

import (
	"context"
	"errors"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/safeblock-dev/wr/taskgroup"
)

func main() {
	ctx := context.Background()

	// Create a new task group (TaskGroup).
	tg := taskgroup.New()

	// Add signal handler to gracefully exit the program on interrupt signals.
	tg.Add(taskgroup.SignalHandler(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM))
	log.Println("We're waiting for 5 seconds, giving you an opportunity to gracefully exit the program.")

	// Actor 1: Performs a long operation and stops on context cancellation.
	tg.AddContext(func(ctx context.Context) error {
		log.Println("Actor 1 working...")
		<-ctx.Done() // Wait for context cancellation
		log.Println("Actor 1 stopped")
		return nil
	}, func(context.Context, error) {
		log.Println("Actor 1 interrupted")
	})

	// Actor 2: Returns an error after a long operation.
	tg.Add(func() error {
		log.Println("Actor 2 working...")
		time.Sleep(5 * time.Second) // Simulate a long operation
		log.Println("Actor 2 stopped")
		return errors.New("error in actor 2")
	}, func(error) {
		log.Println("Actor 2 interrupted")
	})

	// Run all actors and wait for their completion.
	if err := tg.Run(); err != nil {
		log.Println("Error in group:", err)
	}
}
