package main

import (
	"context"
	"errors"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/safeblock-dev/wr/wrtask"
)

// result:
// 2024/03/21 04:06:08 Actor 1 working...
// 2024/03/21 04:06:08 Actor 2 working...
// 2024/03/21 04:06:13 Actor 2 stopped
// 2024/03/21 04:06:13 Actor 1 interrupted
// 2024/03/21 04:06:13 Actor 2 interrupted
// 2024/03/21 04:06:13 Actor 1 stopped

func main() {
	group := wrtask.New()

	group.Add(wrtask.SignalHandler(context.TODO(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM))
	log.Println("We're waiting for 5 seconds, giving you an opportunity to gracefully exit the program.")

	// Actor 1: performs a long operation
	group.AddContext(func(ctx context.Context) error {
		log.Println("Actor 1 working...")
		<-ctx.Done()
		log.Println("Actor 1 stopped")
		return nil
	}, func(context.Context, error) {
		log.Println("Actor 1 interrupted")
	})

	// Actor 2: returns an error
	group.Add(func() error {
		log.Println("Actor 2 working...")
		time.Sleep(5 * time.Second)
		log.Println("Actor 2 stopped")
		return errors.New("error in actor 2")
	}, func(error) {
		log.Println("Actor 2 interrupted")
	})

	// Run all actors and wait for their completion
	if err := group.Run(); err != nil {
		log.Println("Error in group:", err)
	}
}
