package gostreamch

import (
	"sync"
	"sync/atomic"

	"github.com/safeblock-dev/werr"
	"github.com/safeblock-dev/wr/gostream"
)

// StreamCh is a wrapper around gostream.Stream with additional error handling.
type StreamCh struct {
	stream  *gostream.Stream
	err     error
	errCh   chan error
	errOnce sync.Once
	stopped atomic.Bool
}

// New creates a new StreamCh with the provided options.
// It initializes the stream with additional error and panic handlers.
func New(options ...gostream.Option) *StreamCh {
	var s StreamCh
	s.errCh = make(chan error, 1)

	// Combine default error handling options with user-provided options
	allOptions := append([]gostream.Option{
		gostream.ErrorHandler(s.errorHandler),
		gostream.PanicHandler(s.panicHandler),
	}, options...)

	// Initialize the stream with the combined options
	s.stream = gostream.New(allOptions...)

	return &s
}

// Go submits a task to the stream for execution.
func (s *StreamCh) Go(f gostream.Task) {
	if !s.stopped.Load() {
		s.stream.Go(f)
	}
}

// Wait waits for all tasks in the stream to complete and closes the error channel.
func (s *StreamCh) Wait() {
	if s.stopped.CompareAndSwap(false, true) {
		s.stream.Wait()
		close(s.errCh)
	}
}

// Reset reactivates the stream, allowing new tasks to be submitted.
func (s *StreamCh) Reset() {
	s.Wait()
	s.stream.Reset()
	s.stopped.Store(false)
	s.err = nil
	s.errCh = make(chan error, 1)
	s.errOnce = sync.Once{}
}

// ErrorChannel returns a channel that can be used to receive errors that occur in the stream.
func (s *StreamCh) ErrorChannel() <-chan error {
	return s.errCh
}

// Error returns the first error that occurred in the stream.
func (s *StreamCh) Error() error {
	return s.err
}

// HasError returns true if an error has occurred in the stream.
func (s *StreamCh) HasError() bool {
	return s.err != nil
}

// panicHandler handles panics that occur in the stream by converting them
// to errors and passing them to the error handler.
func (s *StreamCh) panicHandler(pc any) {
	s.errorHandler(werr.PanicToError(pc))
}

// errorHandler handles errors that occur in the stream by sending them
// to the error channel and setting the error.
func (s *StreamCh) errorHandler(err error) {
	s.errOnce.Do(func() {
		s.err = err
		s.errCh <- err
		s.stream.Cancel()
	})
}
