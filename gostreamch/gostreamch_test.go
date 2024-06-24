package gostreamch_test

import (
	"errors"
	"sync/atomic"
	"testing"

	"github.com/safeblock-dev/wr/gostream"
	"github.com/safeblock-dev/wr/gostreamch"
	"github.com/stretchr/testify/require"
)

func TestStreamCh_Go(t *testing.T) {
	t.Parallel()

	t.Run("submit task and wait for completion", func(t *testing.T) {
		t.Parallel()

		stream := gostreamch.New()

		// Sample task that completes successfully
		task := func() (gostream.Callback, error) {
			return func() error {
				return nil
			}, nil
		}

		// Submit multiple tasks
		for i := 0; i < 5; i++ {
			stream.Go(task)
		}

		// Wait for all tasks to complete
		stream.Wait()

		require.NoError(t, stream.Error())
		require.False(t, stream.HasError())
		require.NoError(t, <-stream.ErrorChannel())
	})

	t.Run("mix of successful and failing tasks", func(t *testing.T) {
		t.Parallel()

		stream := gostreamch.New()

		var successCalled, skipCalled atomic.Bool
		stream.Go(func() (gostream.Callback, error) {
			return func() error {
				successCalled.Store(true)
				return nil
			}, nil
		})
		stream.Go(func() (gostream.Callback, error) {
			return func() error {
				return errors.New("example error")
			}, nil
		})
		for i := 0; i < 5; i++ {
			stream.Go(func() (gostream.Callback, error) {
				return func() error {
					skipCalled.Store(true)
					return nil
				}, nil
			})
		}

		require.Error(t, <-stream.ErrorChannel())
		require.True(t, stream.HasError())
		require.Error(t, stream.Error())
		require.True(t, successCalled.Load())
		require.False(t, skipCalled.Load())

		// Wait for all tasks to complete
		stream.Wait()
	})

	t.Run("multiple worker errors", func(t *testing.T) {
		t.Parallel()

		stream := gostreamch.New()
		expectedError := errors.New("task error")
		task := func() (gostream.Callback, error) { return nil, expectedError }
		for i := 0; i < 3; i++ {
			stream.Go(task)
		}

		require.True(t, stream.HasError())
		require.Equal(t, expectedError, stream.Error())
		require.ErrorIs(t, expectedError, <-stream.ErrorChannel())
	})

	t.Run("multiple callback errors", func(t *testing.T) {
		t.Parallel()

		stream := gostreamch.New()
		expectedError := errors.New("task error")
		task := func() (gostream.Callback, error) {
			return func() error {
				return expectedError
			}, nil
		}
		for i := 0; i < 3; i++ {
			stream.Go(task)
		}

		require.True(t, stream.HasError())
		require.Equal(t, expectedError, stream.Error())
		require.ErrorIs(t, expectedError, <-stream.ErrorChannel())
	})

	t.Run("error handling", func(t *testing.T) {
		t.Parallel()

		stream := gostreamch.New()

		// Sample task that returns an error
		expectedError := errors.New("task error")
		task := func() (gostream.Callback, error) {
			return nil, expectedError
		}

		// Submit a task that will fail
		stream.Go(task)

		// Wait for all tasks to complete
		stream.Wait()

		require.True(t, stream.HasError())
		require.Equal(t, expectedError, stream.Error())
		require.ErrorIs(t, expectedError, <-stream.ErrorChannel())
	})

	t.Run("callback error handling", func(t *testing.T) {
		t.Parallel()

		stream := gostreamch.New()

		// Sample task that returns an error
		expectedError := errors.New("task error")
		task := func() (gostream.Callback, error) {
			return func() error {
				return expectedError
			}, nil
		}

		// Submit a task that will fail
		stream.Go(task)

		// Wait for all tasks to complete
		stream.Wait()

		require.True(t, stream.HasError())
		require.Equal(t, expectedError, stream.Error())
		require.ErrorIs(t, expectedError, <-stream.ErrorChannel())
	})

	t.Run("panic handling", func(t *testing.T) {
		t.Parallel()

		stream := gostreamch.New()

		// Sample task that panics
		expectedPanic := errors.New("task panic")
		task := func() (gostream.Callback, error) {
			panic(expectedPanic)
		}

		// Submit a task that will panic
		stream.Go(task)

		// Wait for all tasks to complete
		stream.Wait()

		require.True(t, stream.HasError())
		require.ErrorIs(t, stream.Error(), expectedPanic)
		require.Equal(t, stream.Error(), <-stream.ErrorChannel())
	})

	t.Run("panic callback handling", func(t *testing.T) {
		t.Parallel()

		stream := gostreamch.New()

		// Sample task that panics
		expectedPanic := errors.New("task panic")
		task := func() (gostream.Callback, error) {
			return func() error {
				panic(expectedPanic)
			}, nil
		}

		// Submit a task that will panic
		stream.Go(task)

		// Wait for all tasks to complete
		stream.Wait()

		require.True(t, stream.HasError())
		require.ErrorIs(t, stream.Error(), expectedPanic)
		require.Equal(t, stream.Error(), <-stream.ErrorChannel())
	})
}

func TestStreamCh_Reset(t *testing.T) {
	t.Parallel()

	t.Run("should reset the stream", func(t *testing.T) {
		t.Parallel()

		stream := gostreamch.New()
		stream.Wait()
		stream.Reset()

		var completed bool
		stream.Go(func() (gostream.Callback, error) {
			return func() error {
				completed = true
				return nil
			}, nil
		})
		stream.Wait()

		require.True(t, completed)
	})
}

func TestStreamCh_GoAfterWait(t *testing.T) {
	t.Parallel()

	t.Run("should not panic on Go after Wait", func(t *testing.T) {
		t.Parallel()

		stream := gostreamch.New()
		stream.Wait()
		require.NotPanics(t, func() {
			stream.Go(func() (gostream.Callback, error) {
				return func() error {
					return nil
				}, nil
			})
		})
	})
}
