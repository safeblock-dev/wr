package main

import (
	"log"
	"time"

	"github.com/safeblock-dev/wr/syncgroup"
)

func main() {
	// Custom panic handler function that logs the panic value.
	panicHandler := syncgroup.PanicHandler(func(pc any) {
		log.Println("panic:", pc)
	})

	// Create a new WaitGroup instance with the custom panic handler.
	wg := syncgroup.New(panicHandler)

	// Run a function in a new goroutine that logs a message.
	wg.Go(func() {
		log.Println("value:", "foo")
	})

	// Introduce a brief delay to ensure the first goroutine starts.
	time.Sleep(time.Millisecond)

	// Run another function in a new goroutine that logs a message and then panics.
	wg.Go(func() {
		log.Println("send panic")
		panic("my panic")
	})

	// Wait for all goroutines in the WaitGroup to complete.
	wg.Wait()
}
