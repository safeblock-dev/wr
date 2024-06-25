package gopoolch_test

import (
	"errors"
	"testing"

	"github.com/safeblock-dev/wr/gopoolch"
	"github.com/stretchr/testify/require"
)

func TestPoolCh_Go(t *testing.T) {
	t.Parallel()

	t.Run("submit task and wait for completion", func(t *testing.T) {
		t.Parallel()

		pool := gopoolch.New()

		// Sample task that completes successfully
		task := func() error {
			return nil
		}

		// Submit multiple tasks
		for i := 0; i < 5; i++ {
			pool.Go(task)
		}

		// Wait for all tasks to complete
		pool.Wait()

		// Assertions
		require.NoError(t, pool.Error())
		require.False(t, pool.HasError())
		require.NoError(t, <-pool.ErrorChannel())
	})

	t.Run("error handling", func(t *testing.T) {
		t.Parallel()

		pool := gopoolch.New()

		// Sample task that returns an error
		expectedError := errors.New("task error")
		task := func() error {
			return expectedError
		}

		// Submit a success task and a failing task
		pool.Go(func() error { return nil })
		pool.Go(task)

		// Wait for all tasks to complete
		pool.Wait()

		// Assertions
		require.True(t, pool.HasError())
		require.Equal(t, expectedError, pool.Error())
		require.Equal(t, expectedError, <-pool.ErrorChannel())
	})

	t.Run("panic handling", func(t *testing.T) {
		t.Parallel()

		pool := gopoolch.New()

		// Sample task that panics
		expectedPanic := errors.New("task panic")
		task := func() error {
			panic(expectedPanic)
		}

		// Submit a success task and a task that panics
		pool.Go(func() error { return nil })
		pool.Go(task)

		// Wait for all tasks to complete
		pool.Wait()

		// Assertions
		require.True(t, pool.HasError())
		require.ErrorIs(t, pool.Error(), expectedPanic)
		require.Equal(t, pool.Error(), <-pool.ErrorChannel())
	})
}

func TestPoolCh_Reset(t *testing.T) {
	t.Parallel()

	t.Run("should reset the pool", func(t *testing.T) {
		t.Parallel()

		pool := gopoolch.New()
		pool.Wait()
		pool.Reset()

		var completed bool
		task := func() error { completed = true; return nil }
		pool.Go(task)
		pool.Wait()

		require.True(t, completed)
	})
}

func TestPoolCh_GoAfterWait(t *testing.T) {
	t.Parallel()

	t.Run("don't panic on go after wait", func(t *testing.T) {
		t.Parallel()

		pool := gopoolch.New()
		pool.Wait()
		require.NotPanics(t, func() {
			pool.Go(func() error { return nil })
		})
	})
}
