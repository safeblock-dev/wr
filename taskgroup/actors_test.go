package taskgroup_test

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/safeblock-dev/wr/taskgroup"
	"github.com/stretchr/testify/require"
)

func TestContextHandler(t *testing.T) {
	t.Parallel()

	t.Run("context canceled", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		execute, interrupt := taskgroup.ContextHandler(ctx)

		// Cancel the context to trigger termination.
		cancel()

		// The execute function should return an error indicating the context was canceled.
		err := execute()
		require.ErrorIs(t, err, context.Canceled, "Execute should return context.Canceled error")

		// The interrupt function should not cause any error.
		interrupt(nil)
	})

	t.Run("context timeout", func(t *testing.T) {
		t.Parallel()

		// Create a context that will be canceled after a timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		execute, interrupt := taskgroup.ContextHandler(ctx)

		// The execute function should eventually return an error due to the timeout.
		err := execute()
		require.ErrorIs(t, err, context.DeadlineExceeded, "Execute should return context.DeadlineExceeded error")

		// The interrupt function should not cause any error.
		interrupt(nil)
	})

	t.Run("interrupt", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		execute, interrupt := taskgroup.ContextHandler(ctx)

		// Simulate interrupting the handler.
		interrupt(nil)

		// The execute function should return an error indicating the context was canceled.
		err := execute()
		require.ErrorIs(t, err, context.Canceled, "Execute should return context.Canceled error after interrupt")
	})
}

func TestSignalHandler(t *testing.T) {
	t.Parallel()

	t.Run("signal received", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		execute, _ := taskgroup.SignalHandler(ctx, syscall.SIGTERM)

		// Simulate receiving a signal
		go func() {
			time.Sleep(time.Millisecond * 100) // Ensure execute() is running
			_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		}()

		// Wait for signal to be processed
		err := execute()

		require.True(t, taskgroup.IsSignalError(err))
	})

	t.Run("context cancelled", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		execute, _ := taskgroup.SignalHandler(ctx, syscall.SIGTERM)

		// Cancel the context to trigger ctx.Done()
		cancel()

		err := execute()
		require.ErrorIs(t, err, context.Canceled)
	})

	t.Run("interrupt cancelled", func(t *testing.T) {
		t.Parallel()

		execute, interrupt := taskgroup.SignalHandler(context.Background(), syscall.SIGTERM)

		// Cancel the interrupt context
		interrupt(nil)

		err := execute()
		require.ErrorIs(t, err, context.Canceled)
	})
}
