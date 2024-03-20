package wrtask_test

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/safeblock-dev/wr/wrtask"
	"github.com/stretchr/testify/require"
)

func TestContextHandler(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	execute, interrupt := wrtask.ContextHandler(ctx)

	// Simulate context cancellation
	cancel()

	err := execute()
	require.ErrorIs(t, err, context.Canceled)

	// Interrupt should not cause any error
	interrupt(nil)
}

func TestSignalHandler(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	execute, interrupt := wrtask.SignalHandler(ctx, syscall.SIGTERM)

	// Simulate receiving a signal
	go func() {
		time.Sleep(time.Millisecond * 100) // Ensure execute() is running
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	}()

	err := execute()
	require.ErrorIs(t, err, wrtask.ErrSignal)

	// Interrupt should not cause any error
	interrupt(nil)
}
