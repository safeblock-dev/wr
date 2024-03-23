package gostream_test

import (
	"bytes"
	"context"
	"errors"
	"log"
	"sync/atomic"
	"testing"

	"github.com/safeblock-dev/wr/gostream"
	"github.com/safeblock-dev/wr/panics"
	"github.com/stretchr/testify/require"
)

func TestStream_Go(t *testing.T) {
	t.Parallel()

	t.Run("base", func(t *testing.T) {
		t.Parallel()

		stream := gostream.New()
		var workerPoolCounter, callbackCounter atomic.Int64
		for i := 0; i < 10; i++ {
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

		require.Equal(t, int64(10), workerPoolCounter.Load())
		require.Equal(t, int64(10), callbackCounter.Load())
	})

	t.Run("consistency is maintained", func(t *testing.T) {
		t.Parallel()

		syncer := make(chan struct{})
		counter := atomic.Int64{}

		stream := gostream.New(gostream.MaxGoroutines(10))

		// the first task that takes a long time to complete
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
		for i := 1; i < 10; i++ {
			stream.Go(func() (gostream.Callback, error) {
				return func() error {
					require.Equal(t, int64(i), counter.Load())
					counter.Add(1)

					return nil
				}, nil
			})
		}
		syncer <- struct{}{}
		stream.Wait()

		require.Equal(t, int64(10), counter.Load())
	})

	t.Run("starting after Wait", func(t *testing.T) {
		t.Parallel()

		stream := gostream.New()
		stream.Wait()

		isStarted := false
		stream.Go(func() (gostream.Callback, error) {
			isStarted = true

			return nil, nil
		})

		require.False(t, isStarted)
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
		for i := 1; i < 10; i++ {
			stream.Go(func() (gostream.Callback, error) {
				if i == 5 {
					cancel()
				}

				return func() error {
					callbackCounter.Add(1)

					return nil
				}, nil
			})
		}
		stream.Wait()

		require.LessOrEqual(t, callbackCounter.Load(), int64(5))
	})

	t.Run("callback cancel", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		stream := gostream.New(gostream.Context(ctx))

		var callbackCounter atomic.Int64
		for i := 1; i < 10; i++ {
			stream.Go(func() (gostream.Callback, error) {
				return func() error {
					if i == 5 {
						cancel()
					}
					callbackCounter.Add(1)

					return nil
				}, nil
			})
		}
		stream.Wait()

		require.Equal(t, int64(5), callbackCounter.Load())
	})

	t.Run("errorHandler cancel", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		stream := gostream.New(gostream.Context(ctx), gostream.ErrorHandler(func(_ error) {
			cancel()
		}))

		var callbackCounter atomic.Int64
		for i := 1; i < 10; i++ {
			stream.Go(func() (gostream.Callback, error) {
				if i == 5 {
					return nil, errors.New("test error")
				}

				return func() error {
					callbackCounter.Add(1)

					return nil
				}, nil
			})
		}
		stream.Wait()

		require.Equal(t, int64(4), callbackCounter.Load())
	})
}

func TestStream_ErrorHandler(t *testing.T) {
	t.Parallel()

	t.Run("worker pool error", func(t *testing.T) {
		t.Parallel()

		errExpected := errors.New("task error")
		var errorHandlerCalled atomic.Bool
		stream := gostream.New(gostream.ErrorHandler(func(err error) {
			require.Equal(t, errExpected, err)
			errorHandlerCalled.Store(true)
		}))

		stream.Go(func() (gostream.Callback, error) {
			return nil, errExpected
		})

		stream.Wait()
		require.True(t, errorHandlerCalled.Load())
	})

	t.Run("callback pool error", func(t *testing.T) {
		t.Parallel()

		errExpected := errors.New("callback error")
		var errorHandlerCalled atomic.Bool
		stream := gostream.New(gostream.ErrorHandler(func(err error) {
			require.Equal(t, errExpected, err)
			errorHandlerCalled.Store(true)
		}))

		stream.Go(func() (gostream.Callback, error) {
			return func() error {
				return errExpected
			}, nil
		})

		stream.Wait()
		require.True(t, errorHandlerCalled.Load())
	})
}

func TestStream_PanicHandler(t *testing.T) {
	t.Parallel()

	t.Run("worker pool panic", func(t *testing.T) {
		t.Parallel()

		var panicOccurred atomic.Bool
		stream := gostream.New(gostream.PanicHandler(func(_ panics.Recovered) {
			panicOccurred.Store(true)
		}))

		stream.Go(func() (gostream.Callback, error) {
			panic("test panic in worker pool")
		})

		stream.Wait()

		require.True(t, panicOccurred.Load())
	})

	t.Run("callback panic", func(t *testing.T) {
		t.Parallel()

		var panicOccurred atomic.Bool
		stream := gostream.New(gostream.PanicHandler(func(_ panics.Recovered) {
			panicOccurred.Store(true)
		}))

		stream.Go(func() (gostream.Callback, error) {
			return func() error {
				panic("test panic in callback")
			}, nil
		})

		stream.Wait()

		require.True(t, panicOccurred.Load())
	})

	t.Run("error handler panic", func(t *testing.T) {
		t.Parallel()

		var panicOccurred atomic.Bool
		stream := gostream.New(
			gostream.PanicHandler(func(_ panics.Recovered) {
				panicOccurred.Store(true)
			}), gostream.ErrorHandler(func(err error) {
				panic(err)
			}),
		)

		stream.Go(func() (gostream.Callback, error) {
			return func() error {
				return errors.New("trigger error handler")
			}, nil
		})

		stream.Wait()

		require.True(t, panicOccurred.Load())
	})
}

func TestStream_Wait(t *testing.T) {
	t.Parallel()

	stream := gostream.New()
	var callbackExecuted atomic.Bool
	stream.Go(func() (gostream.Callback, error) {
		callbackExecuted.Store(true)

		return nil, nil
	})

	stream.Wait()
	require.True(t, callbackExecuted.Load())

	// Test double wait
	stream.Wait()
}

func TestStream_IsStopped(t *testing.T) {
	t.Parallel()

	t.Run("when wait", func(t *testing.T) {
		t.Parallel()

		stream := gostream.New()
		stream.Go(func() (gostream.Callback, error) {
			return nil, nil
		})

		require.False(t, stream.IsStopped())
		stream.Wait()
		require.True(t, stream.IsStopped())
	})

	t.Run("when context cancel", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		stream := gostream.New(gostream.Context(ctx))
		stream.Go(func() (gostream.Callback, error) {
			return nil, nil
		})

		cancel()
		require.False(t, stream.IsStopped())
		stream.Wait()
		require.True(t, stream.IsStopped())
	})
}

func TestStream_CallbackNilFunction(t *testing.T) {
	t.Parallel()

	stream := gostream.New()
	stream.Go(func() (gostream.Callback, error) {
		return nil, nil // Passing nil as the callback function.
	})

	stream.Wait()
}

func TestDefaultPanicHandler(t *testing.T) {
	t.Parallel()

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
}
