package gostream_test

import (
	"bytes"
	"context"
	"errors"
	"log"
	"sync/atomic"
	"testing"

	"github.com/safeblock-dev/wr/gostream"
	"github.com/stretchr/testify/require"
)

func TestStream_Go(t *testing.T) {
	t.Parallel()

	t.Run("base", func(t *testing.T) {
		t.Parallel()

		// Test basic execution of tasks
		stream := gostream.New()
		var workerPoolCounter, callbackCounter atomic.Int64
		const numTasks = 3
		for i := 0; i < numTasks; i++ {
			stream.Go(func() (gostream.Callback, error) {
				workerPoolCounter.Add(1)

				return func() error {
					isConsistent := callbackCounter.CompareAndSwap(int64(i), int64(i+1))
					require.True(t, isConsistent)

					return nil
				}, nil
			})
		}

		stream.Wait()

		require.Equal(t, int64(numTasks), workerPoolCounter.Load())
		require.Equal(t, int64(numTasks), callbackCounter.Load())
	})

	t.Run("consistency is maintained", func(t *testing.T) {
		t.Parallel()

		// Test that tasks maintain consistency in execution order
		const maxGoroutines = 10
		syncer := make(chan struct{})
		counter := atomic.Int64{}
		stream := gostream.New(gostream.MaxGoroutines(maxGoroutines))

		// Long-running first task
		stream.Go(func() (gostream.Callback, error) {
			syncer <- struct{}{}
			<-syncer

			return func() error {
				require.Zero(t, counter.Load())
				counter.Add(1)
				return nil
			}, nil
		})

		<-syncer
		for i := 1; i < maxGoroutines; i++ {
			index := i
			stream.Go(func() (gostream.Callback, error) {
				return func() error {
					require.Equal(t, int64(index), counter.Load())
					counter.Add(1)
					return nil
				}, nil
			})
		}

		syncer <- struct{}{}
		stream.Wait()

		require.Equal(t, int64(maxGoroutines), counter.Load())
	})

	t.Run("starting after Wait", func(t *testing.T) {
		t.Parallel()

		// Test behavior when starting tasks after Wait
		stream := gostream.New()
		stream.Wait()

		started := false
		stream.Go(func() (gostream.Callback, error) {
			started = true
			return nil, nil
		})

		require.False(t, started)
	})
}

func TestStream_ContextCancellation(t *testing.T) {
	t.Parallel()

	t.Run("worker pool cancel", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		stream := gostream.New(gostream.Context(ctx))

		var callbackCounter atomic.Int64
		const numTasks = 10
		const cancelIndex = 5
		for i := 1; i < numTasks; i++ {
			index := i
			stream.Go(func() (gostream.Callback, error) {
				if index == cancelIndex {
					cancel() // Cancel the context on the 5th task
				}

				return func() error {
					callbackCounter.Add(1)

					return nil
				}, nil
			})
		}
		stream.Wait()

		require.LessOrEqual(t, callbackCounter.Load(), int64(cancelIndex))
	})

	t.Run("callback cancel", func(t *testing.T) {
		t.Parallel()

		// Test cancellation of tasks via callback cancellation
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		stream := gostream.New(gostream.Context(ctx))

		var callbackCounter atomic.Int64
		const numTasks = 10
		const cancelIndex = 5
		for i := 1; i < numTasks; i++ {
			index := i
			stream.Go(func() (gostream.Callback, error) {
				return func() error {
					if index == cancelIndex {
						cancel()
					}
					callbackCounter.Add(1)

					return nil
				}, nil
			})
		}
		stream.Wait()

		require.Equal(t, int64(cancelIndex), callbackCounter.Load())
	})

	t.Run("errorHandler cancel", func(t *testing.T) {
		t.Parallel()

		// Test cancellation of tasks via errorHandler
		ctx, cancel := context.WithCancel(context.Background())
		errHandler := func(_ error) { cancel() }
		stream := gostream.New(gostream.Context(ctx), gostream.ErrorHandler(errHandler))

		var callbackCounter atomic.Int64
		const numTasks = 10
		const errorIndex = 5
		for i := 1; i < numTasks; i++ {
			index := i
			stream.Go(func() (gostream.Callback, error) {
				if index == errorIndex {
					return nil, errors.New("test error") // Trigger error to cancel
				}

				return func() error {
					callbackCounter.Add(1)

					return nil
				}, nil
			})
		}
		stream.Wait()

		require.Equal(t, int64(errorIndex-1), callbackCounter.Load())
	})
}

func TestStream_ErrorHandler(t *testing.T) {
	t.Parallel()

	t.Run("worker pool error", func(t *testing.T) {
		t.Parallel()

		// Test handling of errors in the worker pool
		expectedError := errors.New("task error")
		var errorCalled atomic.Bool
		var errorCounter atomic.Int64
		errHandler := func(err error) {
			require.Equal(t, expectedError, err)
			errorCalled.Store(true)
			errorCounter.Add(1)
		}

		stream := gostream.New(gostream.ErrorHandler(errHandler))
		stream.Go(func() (gostream.Callback, error) {
			return nil, expectedError
		})
		stream.Wait()

		require.True(t, errorCalled.Load())
		require.Equal(t, int64(1), errorCounter.Load())
	})

	t.Run("callback pool error", func(t *testing.T) {
		t.Parallel()

		// Test handling of errors in the callback
		expectedError := errors.New("callback error")
		var errorCalled atomic.Bool
		var errorCounter atomic.Int64
		errHandler := func(err error) {
			require.Equal(t, expectedError, err)
			errorCalled.Store(true)
			errorCounter.Add(1)
		}

		stream := gostream.New(gostream.ErrorHandler(errHandler))
		stream.Go(func() (gostream.Callback, error) {
			return func() error {
				return expectedError
			}, nil
		})
		stream.Wait()

		require.True(t, errorCalled.Load())
		require.Equal(t, int64(1), errorCounter.Load())
	})
}

func TestStream_PanicHandler(t *testing.T) {
	t.Parallel()

	t.Run("default worker panic", func(t *testing.T) { //nolint: paralleltest
		// Set up a buffer to capture log output.
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		defer func() {
			// Reset the log output to its default (stderr) after the test.
			log.SetOutput(nil)
		}()

		wg := gostream.New()
		wg.Go(func() (gostream.Callback, error) {
			panic("test panic")
		})
		wg.Wait()

		// Check that the log output contains the expected panic information.
		require.Contains(t, logBuffer.String(), "test panic")
	})

	t.Run("default callback panic", func(t *testing.T) { //nolint: paralleltest
		// Set up a buffer to capture log output.
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		defer func() {
			// Reset the log output to its default (stderr) after the test.
			log.SetOutput(nil)
		}()

		wg := gostream.New()
		wg.Go(func() (gostream.Callback, error) {
			return func() error { panic("test panic") }, nil
		})
		wg.Wait()

		// Check that the log output contains the expected panic information.
		require.Contains(t, logBuffer.String(), "test panic")
	})

	t.Run("worker panic", func(t *testing.T) {
		t.Parallel()

		var panicCalled atomic.Bool
		var panicCounter atomic.Int64
		errHandler := func(_ any) {
			panicCalled.Store(true)
			panicCounter.Add(1)
		}

		stream := gostream.New(gostream.PanicHandler(errHandler))
		stream.Go(func() (gostream.Callback, error) {
			panic("test panic in worker pool")
		})
		stream.Wait()

		require.True(t, panicCalled.Load())
		require.Equal(t, int64(1), panicCounter.Load())
	})

	t.Run("callback panic", func(t *testing.T) {
		t.Parallel()

		var panicCalled atomic.Bool
		var panicCounter atomic.Int64
		panicHandler := func(_ any) {
			panicCalled.Store(true)
			panicCounter.Add(1)
		}

		stream := gostream.New(gostream.PanicHandler(panicHandler))
		stream.Go(func() (gostream.Callback, error) {
			return func() error {
				panic("test panic in callback")
			}, nil
		})
		stream.Wait()

		require.True(t, panicCalled.Load())
		require.Equal(t, int64(1), panicCounter.Load())
	})

	t.Run("error handler call panic", func(t *testing.T) {
		t.Parallel()

		var panicCalled atomic.Bool
		var panicCounter atomic.Int64
		panicHandler := func(_ any) {
			panicCalled.Store(true)
			panicCounter.Add(1)
		}
		errHandler := func(err error) {
			panic(err)
		}

		stream := gostream.New(
			gostream.PanicHandler(panicHandler),
			gostream.ErrorHandler(errHandler),
		)
		stream.Go(func() (gostream.Callback, error) {
			return func() error {
				return errors.New("trigger error handler")
			}, nil
		})
		stream.Wait()

		require.True(t, panicCalled.Load())
		require.Equal(t, int64(1), panicCounter.Load())
	})
}

func TestStream_Wait(t *testing.T) {
	t.Parallel()

	stream := gostream.New()
	var callbackExecuted atomic.Bool
	stream.Go(func() (gostream.Callback, error) {
		callbackExecuted.Store(true)

		return nil, nil // Passing nil as the callback function.
	})

	stream.Wait()
	require.True(t, callbackExecuted.Load())

	// Test double Wait
	stream.Wait()
}

func TestStream_CallbackNilFunction(t *testing.T) {
	t.Parallel()

	// Test when a nil function is passed as the callback
	stream := gostream.New()
	stream.Go(func() (gostream.Callback, error) {
		return nil, nil // Passing nil as the callback function.
	})

	stream.Wait()
}

func TestPool_Reset(t *testing.T) {
	t.Parallel()

	t.Run("should reset the stream", func(t *testing.T) {
		t.Parallel()

		// Test resetting the stream
		stream := gostream.New()
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
