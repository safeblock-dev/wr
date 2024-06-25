package gostream

import (
	"context"
	"sync/atomic"

	"github.com/safeblock-dev/wr/gopool"
)

// Stream manages the execution of tasks and their corresponding callbacks.
type Stream struct {
	ctx             context.Context      // ctx is the current context for the stream.
	parentCtx       context.Context      // parentCtx is the parent context of the stream.
	cancelFunc      context.CancelFunc   // cancelFunc cancels the stream context.
	callbackQueueCh chan callbackChannel // callbackQueueCh is a channel for callback channels.
	panicHandler    func(any)            // panicHandler handles panics that occur in tasks.
	errorHandler    func(err error)      // errorHandler handles errors that occur in tasks.
	workerPool      *gopool.Pool         // workerPool manages the goroutines executing tasks.
	maxGoroutines   int                  // maxGoroutines is the maximum number of concurrent goroutines.
	stopped         atomic.Bool          // stopped indicates if the stream has been stopped.
}

// Task is a function that returns a Callback and an error.
type Task func() (Callback, error)

// Callback is a function that is executed after a Task completes.
type Callback func() error

// New creates a new Stream with the provided options.
func New(options ...Option) *Stream {
	stream := &Stream{ //nolint: exhaustruct
		panicHandler: defaultPanicHandler,
	}

	// Apply all options.
	for _, opt := range options {
		opt(stream)
	}

	// Initialize base context (if not already set).
	if stream.ctx == nil {
		Context(context.Background())(stream)
	}

	if stream.maxGoroutines > 0 {
		stream.maxGoroutines++
	}

	stream.workerPool = gopool.New(gopool.MaxGoroutines(stream.maxGoroutines))
	stream.callbackQueueCh = make(chan callbackChannel, stream.maxGoroutines+1)

	// Start the callback reader with panic protection.
	stream.workerPool.Go(func() error { stream.callbackReader(); return nil }) //nolint: nlreturn

	return stream
}

// Go submits a Task to the Stream for execution.
func (s *Stream) Go(f Task) {
	if s.ctx.Err() != nil {
		return
	}

	queueCh := getCallbackChannel()
	s.callbackQueueCh <- queueCh

	// Submit the task for execution with panic protection.
	s.workerPool.Go(func() error {
		defer func() {
			// Recover from any potential panic in the task function and send a
			// callbackData with the panic information to the callback reader. This
			// ensures that the callback reader is not blocked waiting for a callback
			// that will never come due to the panic.
			if r := recover(); r != nil {
				defer func() {
					queueCh <- callbackData{fn: nil, err: nil}
				}()
				s.panicHandler(r)
			}
		}()

		// Execute the task function and send its result or error (if any) to the
		// callback reader through the queue channel.
		callbackFn, err := f()
		queueCh <- callbackData{fn: callbackFn, err: err}

		return nil
	})
}

// Reset reactivates the stream, allowing new tasks to be submitted.
func (s *Stream) Reset() {
	s.Wait()
	Context(s.parentCtx)(s)
	s.workerPool.Reset()
	s.callbackQueueCh = make(chan callbackChannel, s.maxGoroutines+1)
	s.workerPool.Go(func() error { s.callbackReader(); return nil }) //nolint: nlreturn
	s.stopped.Store(false)
}

// Wait blocks until all tasks and their callbacks have been executed.
func (s *Stream) Wait() {
	if s.stopped.CompareAndSwap(false, true) {
		close(s.callbackQueueCh)
		s.workerPool.Wait()
		s.cancelFunc()
	}
}

// Cancel cancels the stream, stopping all pending tasks.
func (s *Stream) Cancel() {
	s.cancelFunc()
}

// callbackReader reads callbacks from the callbackQueueCh and executes them.
func (s *Stream) callbackReader() {
	for queueCh := range s.callbackQueueCh {
		data := <-queueCh

		putCallbackChannel(queueCh)
		s.callbackHandler(data)
	}
}

// callbackHandler executes a callback and handles errors.
func (s *Stream) callbackHandler(data callbackData) {
	defer func() {
		if r := recover(); r != nil {
			s.panicHandler(r)
		}
	}()

	if s.ctx.Err() != nil {
		return
	}
	if data.err != nil {
		s.errorHandler(data.err)
	}
	if data.fn == nil {
		return
	}

	if err := data.fn(); err != nil {
		s.errorHandler(err)
	}
}
