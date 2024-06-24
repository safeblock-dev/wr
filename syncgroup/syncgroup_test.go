package syncgroup_test

import (
	"bytes"
	"log"
	"sync/atomic"
	"testing"

	"github.com/safeblock-dev/wr/syncgroup"
	"github.com/stretchr/testify/require"
)

func TestWaitGroup_Go(t *testing.T) {
	t.Parallel()

	t.Run("simple", func(t *testing.T) {
		t.Parallel()

		var counter atomic.Int64
		wg := syncgroup.New()

		for i := 0; i < 10; i++ {
			wg.Go(func() {
				counter.Add(1)
			})
		}

		wg.Wait()
		require.Equal(t, int64(10), counter.Load())
	})
}

func TestPanicHandler(t *testing.T) {
	t.Parallel()

	t.Run("default panic", func(t *testing.T) { //nolint: paralleltest
		// Set up a buffer to capture log output.
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		defer func() {
			// Reset the log output to its default (stderr) after the test.
			log.SetOutput(nil)
		}()

		wg := syncgroup.New()
		wg.Go(func() {
			panic("test panic")
		})
		wg.Wait()

		// Check that the log output contains the expected panic information.
		require.Contains(t, logBuffer.String(), "test panic")
	})

	t.Run("custom panic handler", func(t *testing.T) {
		t.Parallel()

		panicHandled := false
		panicHandler := func(_ any) {
			panicHandled = true
		}

		wg := syncgroup.New(syncgroup.PanicHandler(panicHandler))
		wg.Go(func() {
			panic("test panic")
		})
		wg.Wait()

		require.True(t, panicHandled, "Panic should be handled")
	})
}
