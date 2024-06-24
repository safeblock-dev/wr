package syncgroup_test

import (
	"bytes"
	"log"
	"sync/atomic"
	"testing"

	"github.com/safeblock-dev/wr/syncgroup"
	"github.com/stretchr/testify/require"
)

// TestWaitGroup_Go tests the WaitGroup's Go method for concurrent task execution.
func TestWaitGroup_Go(t *testing.T) {
	t.Parallel()

	t.Run("increments counter correctly", func(t *testing.T) {
		t.Parallel()

		var counter atomic.Int64
		wg := syncgroup.New()

		// Launch 10 goroutines that each increment the counter.
		for i := 0; i < 10; i++ {
			wg.Go(func() {
				counter.Add(1)
			})
		}

		// Wait for all goroutines to complete.
		wg.Wait()

		// Assert that the counter value matches the number of goroutines launched.
		require.Equal(t, int64(10), counter.Load())
	})
}

// TestPanicHandler tests panic handling capabilities of the WaitGroup.
func TestPanicHandler(t *testing.T) {
	t.Parallel()

	t.Run("logs default panic", func(t *testing.T) { //nolint: paralleltest
		// Set up a buffer to capture log output.
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		defer func() {
			// Reset the log output to its default (stderr) after the test.
			log.SetOutput(nil)
		}()

		// Create a new WaitGroup.
		wg := syncgroup.New()

		// Launch a goroutine that panics.
		wg.Go(func() {
			panic("test panic")
		})

		// Wait for all goroutines to complete.
		wg.Wait()

		// Check that the log output contains the panic message.
		require.Contains(t, logBuffer.String(), "test panic")
	})

	t.Run("handles panic with custom handler", func(t *testing.T) {
		t.Parallel()

		panicHandled := false

		// Define a custom panic handler.
		panicHandler := func(_ any) {
			panicHandled = true
		}

		// Create a new WaitGroup with the custom panic handler.
		wg := syncgroup.New(syncgroup.PanicHandler(panicHandler))

		// Launch a goroutine that panics.
		wg.Go(func() {
			panic("test panic")
		})

		// Wait for all goroutines to complete.
		wg.Wait()

		// Assert that the custom panic handler was invoked.
		require.True(t, panicHandled, "Panic should be handled")
	})
}
