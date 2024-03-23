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

	ctx, cancel := context.WithCancel(context.Background())
	execute, interrupt := taskgroup.ContextHandler(ctx)

	// Simulate context cancellation
	cancel()

	err := execute()
	require.ErrorIs(t, err, context.Canceled)

	// Interrupt should not cause any error
	interrupt(nil)
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

		err := execute()
		require.ErrorIs(t, err, taskgroup.ErrSignal)
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
}
