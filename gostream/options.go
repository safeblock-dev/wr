package gostream

import (
	"context"
	"log"

	"github.com/safeblock-dev/wr/gopool"
	"github.com/safeblock-dev/wr/panics"
)

// Constants for the default capacity of options slices.
const (
	workerPoolOptsCount = 2
)

// Option represents an option that can be passed when instantiating a Stream to customize it.
type Option func(stream *Stream)

// DefaultPanicHandler is the default panic handler that prints the panic information.
func DefaultPanicHandler(recovered panics.Recovered) {
	log.Println(recovered.String())
}

// PanicHandler sets the panic handler function for the stream.
func PanicHandler(panicHandler func(recovered panics.Recovered)) Option {
	return func(stream *Stream) {
		stream.panicHandler = panicHandler
	}
}

// ErrorHandler sets the error handler function for the stream.
func ErrorHandler(errorHandler func(err error)) Option {
	return func(stream *Stream) {
		stream.errorHandler = errorHandler
	}
}

// Context sets a parent context for the stream to stop all workers when it is cancelled.
func Context(ctx context.Context) Option {
	return func(stream *Stream) {
		stream.ctx = ctx
		stream.workerPoolOpts = append(stream.workerPoolOpts, gopool.Context(ctx))
	}
}

// MaxGoroutines sets the maximum number of goroutines allowed in the worker pool.
func MaxGoroutines(limit int) Option {
	return func(stream *Stream) {
		stream.workerPoolOpts = append(stream.workerPoolOpts, gopool.MaxGoroutines(limit))
	}
}
