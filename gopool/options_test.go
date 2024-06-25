package gopool //nolint: testpackage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPool_MaxGoroutines(t *testing.T) {
	t.Parallel()

	const maxConcurrent = 5

	t.Run("should set the maximum number of goroutines", func(t *testing.T) {
		t.Parallel()

		pool := New(MaxGoroutines(maxConcurrent))
		require.Equal(t, maxConcurrent, cap(pool.limiter), "limiter capacity should be set to maxConcurrent")
	})

	t.Run("should initialize limiter correctly after reset", func(t *testing.T) {
		t.Parallel()

		pool := New()
		require.Nil(t, pool.limiter, "limiter should be nil initially")

		// Reconfigure the pool with a new limit
		pool.Reset()
		MaxGoroutines(maxConcurrent)(pool)

		require.Equal(t, maxConcurrent, cap(pool.limiter), "limiter capacity should be set to maxConcurrent after reset")
	})

	t.Run("should close old limiter channel", func(t *testing.T) {
		t.Parallel()

		pool := New(MaxGoroutines(maxConcurrent))
		oldLimiter := pool.limiter
		MaxGoroutines(maxConcurrent)(pool)

		select {
		case _, ok := <-oldLimiter:
			require.False(t, ok, "old limiter channel should be closed")
		default:
			require.Fail(t, "old limiter channel should be closed but it is not")
		}
	})
}
