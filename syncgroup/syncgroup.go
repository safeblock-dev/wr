package syncgroup

import (
	"sync"
)

// WaitGroup is a wrapper around sync.WaitGroup with a custom panic handler.
type WaitGroup struct {
	panicHandler func(pc any) // panicHandler is a function to handle panics.
	wg           sync.WaitGroup
}

// New creates a new WaitGroup with the provided options.
func New(options ...Option) *WaitGroup {
	wg := &WaitGroup{
		panicHandler: defaultPanicHandler, // Set default panic handler.
		wg:           sync.WaitGroup{},    // Initialize embedded WaitGroup.
	}

	// Apply all options.
	for _, opt := range options {
		opt(wg)
	}

	return wg
}

// Go runs the given function in a new goroutine and handles panics using the panicHandler.
func (wg *WaitGroup) Go(f func()) {
	wg.wg.Add(1) // Increment the WaitGroup counter.
	go func() {
		defer wg.wg.Done() // Decrement the WaitGroup counter when done.
		defer func() {
			if pc := recover(); pc != nil && wg.panicHandler != nil {
				wg.panicHandler(pc) // Call panic handler on recovery.
			}
		}()

		f() // Execute the provided function.
	}()
}

// Wait waits for all goroutines in the WaitGroup to complete.
func (wg *WaitGroup) Wait() {
	wg.wg.Wait() // Wait for all goroutines to finish.
}
