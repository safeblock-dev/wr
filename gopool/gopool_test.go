package gopool_test

import (
	"bytes"
	"errors"
	"log"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/safeblock-dev/wr/gopool"
	"github.com/stretchr/testify/require"
)

// TestPool_Go tests the Go method of the gopool.Pool.
func TestPool_Go(t *testing.T) {
	t.Parallel()

	t.Run("increments counter correctly", func(t *testing.T) {
		t.Parallel()

		var counter atomic.Uint64
		pool := gopool.New()
		numTasks := 10

		for i := 0; i < numTasks; i++ {
			pool.Go(func() error {
				counter.Add(1)
				return nil
			})
		}

		pool.Wait()

		require.EqualValues(t, numTasks, counter.Load(), "All tasks should be executed")
	})

	t.Run("with max goroutines", func(t *testing.T) {
		t.Parallel()

		// Test different values for maximum concurrent goroutines.
		for _, maxConcurrent := range []int{1, 10, 100} {
			t.Run(strconv.Itoa(maxConcurrent), func(t *testing.T) {
				t.Parallel()

				// Create a new pool with the specified maximum number of concurrent goroutines.
				g := gopool.New(gopool.MaxGoroutines(maxConcurrent))

				var currentConcurrent atomic.Int64 // Tracks the current number of concurrent goroutines.
				var errCount atomic.Int64          // Tracks the number of times the concurrency limit is exceeded.

				taskCount := maxConcurrent * 10 // Total number of tasks to submit.

				for i := 0; i < taskCount; i++ {
					g.Go(func() error {
						cur := currentConcurrent.Add(1) // Increment the concurrent counter.

						// Check if the concurrency limit is exceeded.
						if cur > int64(maxConcurrent) {
							errCount.Add(1) // Increment the error count.
						}

						time.Sleep(time.Millisecond) // Simulate some work.

						currentConcurrent.Add(-1) // Decrement the concurrent counter.

						return nil
					})
				}

				g.Wait() // Wait for all tasks to complete.

				// Verify that the concurrency limit was never exceeded.
				require.Zero(t, errCount.Load())

				// Verify that all goroutines have completed.
				require.Zero(t, currentConcurrent.Load())
			})
		}
	})

	t.Run("handles errors", func(t *testing.T) {
		t.Parallel()

		var counter atomic.Uint64
		pool := gopool.New()
		numTasks := 10

		for i := 0; i < numTasks; i++ {
			if i%3 == 0 {
				pool.Go(func() error {
					counter.Add(1)
					return errors.New("example error")
				})
				continue
			}
			pool.Go(func() error {
				counter.Add(1)
				return nil
			})
		}

		pool.Wait()

		require.EqualValues(t, numTasks, counter.Load())
	})
}

// TestErrorHandler tests the error handling capability of the gopool.Pool.
func TestErrorHandler(t *testing.T) {
	t.Parallel()

	t.Run("should handle errors", func(t *testing.T) {
		t.Parallel()

		var errCounter atomic.Uint64
		errorHandler := func(_ error) {
			errCounter.Add(1)
		}

		pool := gopool.New(gopool.ErrorHandler(errorHandler))
		pool.Go(func() error {
			return errors.New("test error")
		})

		pool.Wait()

		require.Equal(t, uint64(1), errCounter.Load(), "Error handler should be called")
	})
}

// TestPool_Reset tests the Reset method of the gopool.Pool.
func TestPool_Reset(t *testing.T) {
	t.Parallel()

	t.Run("should reset the pool", func(t *testing.T) {
		t.Parallel()

		pool := gopool.New()
		pool.Wait()
		pool.Reset()

		var completed bool
		task := func() error { completed = true; return nil }
		pool.Go(task)
		pool.Wait()

		require.True(t, completed)
	})

	t.Run("with max goroutines limit", func(t *testing.T) {
		t.Parallel()

		pool := gopool.New(gopool.MaxGoroutines(5))
		pool.Wait()
		pool.Reset()

		var completed bool
		task := func() error { completed = true; return nil }
		pool.Go(task)
		pool.Wait()

		require.True(t, completed)
	})
}

// TestPool_Errors tests error handling and panic recovery in the gopool.Pool.
func TestPool_Errors(t *testing.T) {
	t.Parallel()

	t.Run("should handle panics", func(t *testing.T) {
		t.Parallel()

		errorHandled := false
		errorHandler := gopool.ErrorHandler(func(_ error) {
			errorHandled = true
		})
		pool := gopool.New(errorHandler)

		pool.Go(func() error {
			return errors.New("error")
		})

		pool.Wait()

		require.True(t, errorHandled, "Error should be handled")
	})
}

// TestPool_Panics tests panic handling in the gopool.Pool.
func TestPool_Panics(t *testing.T) {
	t.Parallel()

	t.Run("logs default panic", func(t *testing.T) { //nolint: paralleltest
		// Set up a buffer to capture log output.
		var logBuffer bytes.Buffer
		log.SetOutput(&logBuffer)
		defer func() {
			// Reset the log output to its default (stderr) after the test.
			log.SetOutput(nil)
		}()

		wg := gopool.New()
		wg.Go(func() error {
			panic("test panic")
		})
		wg.Wait()

		// Check that the log output contains the expected panic information.
		require.Contains(t, logBuffer.String(), "test panic")
	})

	t.Run("should handle panics", func(t *testing.T) {
		t.Parallel()

		panicHandled := false
		panicHandler := gopool.PanicHandler(func(any) {
			panicHandled = true
		})
		pool := gopool.New(panicHandler)

		pool.Go(func() error {
			panic("test panic")
		})

		pool.Wait()

		require.True(t, panicHandled, "Panic should be handled")
	})
}

// TestPool_DoubleWait tests behavior of double Wait calls on gopool.Pool.
func TestPool_DoubleWait(t *testing.T) {
	t.Parallel()

	t.Run("should not panic on double Wait", func(t *testing.T) {
		t.Parallel()

		pool := gopool.New()

		pool.Wait()
		require.NotPanics(t, func() { pool.Wait() }, "Double Wait should not cause a panic")
	})
}

// TestPool_GoAfterWait tests behavior of Go method after Wait on gopool.Pool.
func TestPool_GoAfterWait(t *testing.T) {
	t.Parallel()

	t.Run("should not panic on Go after Wait", func(t *testing.T) {
		t.Parallel()

		pool := gopool.New()
		pool.Wait()

		require.NotPanics(t, func() {
			pool.Go(func() error { return nil })
		}, "Go after Wait should not cause a panic")
	})
}
