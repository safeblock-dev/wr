package wrpool_test

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/safeblock-dev/wr/panics"
	"github.com/safeblock-dev/wr/wrpool"
	"github.com/stretchr/testify/require"
)

func TestPool_Go_Wait(t *testing.T) {
	t.Parallel()

	t.Run("should execute tasks", func(t *testing.T) {
		t.Parallel()

		var counter atomic.Uint64
		pool := wrpool.New()
		numTasks := 10

		for i := 0; i < numTasks; i++ {
			pool.Go(func() error {
				counter.Add(1)

				return nil
			})
		}

		pool.Wait()

		require.Equal(t, uint64(numTasks), counter.Load(), "All tasks should be executed")
	})

	t.Run("should handle errors", func(t *testing.T) {
		t.Parallel()

		var errCounter atomic.Uint64
		errorHandler := func(_ error) {
			errCounter.Add(1)
		}

		pool := wrpool.New(wrpool.ErrorHandler(errorHandler))
		pool.Go(func() error {
			return errors.New("test error")
		})

		pool.Wait()

		require.Equal(t, uint64(1), errCounter.Load(), "Error handler should be called")
	})
}

func TestPool_Reset(t *testing.T) {
	t.Parallel()

	t.Run("should reset the pool", func(t *testing.T) {
		t.Parallel()

		pool := wrpool.New()
		pool.Go(func() error {
			time.Sleep(100 * time.Millisecond)

			return nil
		})

		pool.Wait()
		pool.Reset()

		require.False(t, pool.IsStopped(), "Pool should be active after reset")
	})
}

func TestPool_Errors(t *testing.T) {
	t.Parallel()

	t.Run("should handle panics", func(t *testing.T) {
		t.Parallel()

		errorHandled := false
		errorHandler := wrpool.ErrorHandler(func(_ error) {
			errorHandled = true
		})
		pool := wrpool.New(errorHandler)

		pool.Go(func() error {
			return errors.New("error")
		})

		pool.Wait()

		require.True(t, errorHandled, "Error should be handled")
	})
}

func TestPool_Panics(t *testing.T) {
	t.Parallel()

	t.Run("should handle panics", func(t *testing.T) {
		t.Parallel()

		panicHandled := false
		panicHandler := wrpool.PanicHandler(func(_ panics.Recovered) {
			panicHandled = true
		})
		pool := wrpool.New(panicHandler)

		pool.Go(func() error {
			panic("test panic")
		})

		pool.Wait()

		require.True(t, panicHandled, "Panic should be handled")
	})
}

func TestPool_DoubleWait(t *testing.T) {
	t.Parallel()

	t.Run("should not panic on double Wait", func(t *testing.T) {
		t.Parallel()

		pool := wrpool.New()

		pool.Wait()
		require.NotPanics(t, func() { pool.Wait() }, "Double Wait should not cause a panic")
	})
}

func TestPool_GoAfterWait(t *testing.T) {
	t.Parallel()

	t.Run("should not execute Go after Wait", func(t *testing.T) {
		t.Parallel()

		pool := wrpool.New()

		pool.Wait()

		executed := false
		pool.Go(func() error {
			executed = true

			return nil
		})

		time.Sleep(100 * time.Millisecond) // Give some time for the task to potentially execute

		require.False(t, executed, "Go should not execute a task after Wait has been called")
	})
}
