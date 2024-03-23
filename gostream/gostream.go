package gostream

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/safeblock-dev/wr/gopool"
	"github.com/safeblock-dev/wr/panics"
	"github.com/safeblock-dev/wr/syncgroup"
)

// Stream manages the execution of tasks and their corresponding callbacks.
type Stream struct {
	context         context.Context
	callbackQueueCh chan callbackChannel
	panicHandler    func(recovered panics.Recovered)
	errorHandler    func(err error) // Handler for errors occurring during task execution.
	contextCancel   context.CancelFunc
	workerPoolOpts  []gopool.Option
	callbackGroup   syncgroup.WaitGroup
	workerPool      gopool.Pool
	initOnce        sync.Once
	stopped         atomic.Bool
}

// Task is a function that returns a Callback and an error.
type Task func() (Callback, error)

// Callback is a function that is executed after a Task completes.
type Callback func() error

// New creates a new Stream with the provided options.
func New(options ...Option) *Stream {
	stream := &Stream{
		panicHandler:   DefaultPanicHandler,
		workerPoolOpts: make([]gopool.Option, 0, workerPoolOptsCount),
		callbackGroup:  *syncgroup.New(),
	}

	// Apply all options.
	for _, opt := range options {
		opt(stream)
	}

	// Initialize base context (if not already set).
	if stream.context == nil {
		Context(context.Background())(stream)
	}

	stream.workerPool = *gopool.New(stream.workerPoolOpts...)
	stream.callbackQueueCh = make(chan callbackChannel, stream.workerPool.MaxGoroutines())

	return stream
}

// Go submits a Task to the Stream for execution.
func (s *Stream) Go(f Task) {
	if s.IsStopped() {
		return
	}

	s.initOnce.Do(func() {
		s.callbackGroup.Go(s.callbackReader) // Start the callback reader with panic protection.
	})

	queueCh := getCallbackChannel()

	// Submit the task for execution with panic protection.
	ok := s.workerPool.Go(func() error {
		defer func() {
			// Recover from any potential panic in the task function and send a
			// callbackData with the panic information to the callback reader. This
			// ensures that the callback reader is not blocked waiting for a callback
			// that will never come due to the panic.
			if r := recover(); r != nil {
				p := panics.NewRecovered(1, r)
				queueCh <- callbackData{fn: nil, err: nil, panic: &p}
			}
		}()

		// Execute the task function and send its result or error (if any) to the
		// callback reader through the queue channel.
		callbackFn, err := f()
		queueCh <- callbackData{fn: callbackFn, err: err, panic: nil}

		return nil
	})
	if ok {
		s.callbackQueueCh <- queueCh
	}
}

// Wait blocks until all tasks and their callbacks have been executed.
func (s *Stream) Wait() {
	if s.stopped.CompareAndSwap(false, true) {
		defer func() {
			close(s.callbackQueueCh)
			s.callbackGroup.Wait()
		}()

		s.workerPool.Wait()
	}
}

// callbackReader reads callbacks from the callbackQueueCh and executes them.
func (s *Stream) callbackReader() {
	for queueCh := range s.callbackQueueCh {
		data := <-queueCh

		s.callbackHandler(data)
		putCallbackChannel(queueCh)
	}
}

// callbackHandler executes a callback and handles errors.
func (s *Stream) callbackHandler(data callbackData) {
	defer func() {
		if r := recover(); r != nil {
			s.panicHandler(panics.NewRecovered(1, r))
		}
	}()

	switch {
	case data.panic != nil:
		s.panicHandler(*data.panic)
	case data.err != nil:
		s.errorHandler(data.err)
	}
	if s.context.Err() != nil {
		return
	}

	if data.fn != nil {
		err := data.fn()
		if err != nil {
			s.errorHandler(err)
		}
	}
}

// IsStopped returns true if the stream is stopped.
func (s *Stream) IsStopped() bool {
	return s.stopped.Load()
}
