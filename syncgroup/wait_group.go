package wrgroup

import (
	"sync"

	"github.com/safeblock-dev/wr/panics"
)

// WaitGroup is a wrapper around sync.WaitGroup with a custom panic handler.
type WaitGroup struct {
	panicHandler func(recovered panics.Recovered)
	wg           sync.WaitGroup
}

// New creates a new WaitGroup with the provided options.
func New(options ...Option) *WaitGroup {
	wg := &WaitGroup{
		panicHandler: DefaultPanicHandler,
		wg:           sync.WaitGroup{},
	}

	// Apply all options.
	for _, opt := range options {
		opt(wg)
	}

	return wg
}

// Go runs the given function in a new goroutine and handles panics using the panicHandler.
func (wg *WaitGroup) Go(f func()) {
	wg.wg.Add(1)
	go func() {
		defer wg.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				wg.panicHandler(panics.NewRecovered(1, r))
			}
		}()

		f()
	}()
}

// Wait waits for all goroutines in the WaitGroup to complete.
func (wg *WaitGroup) Wait() {
	wg.wg.Wait()
}
