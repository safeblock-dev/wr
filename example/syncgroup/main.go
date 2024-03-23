package main

import (
	"log"
	"time"

	"github.com/safeblock-dev/wr/panics"
	"github.com/safeblock-dev/wr/syncgroup"
)

// result:
// 2024/03/20 23:20:26 value: foo
// 2024/03/20 23:20:26 send panic
// 2024/03/20 23:20:26 panic: my panic

func main() {
	// Create a custom panic handler that logs the panic value.
	panicHandler := syncgroup.PanicHandler(func(recovered panics.Recovered) {
		log.Println("panic:", recovered.Value)
	})

	// Create a new WaitGroup with the custom panic handler.
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
