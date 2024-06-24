package gostream

import (
	"context"
	"log"
)

// Option represents an option that can be passed when instantiating a Stream to customize it.
type Option func(stream *Stream)

// PanicHandler sets the panic handler function for the stream.
func PanicHandler(panicHandler func(pc any)) Option {
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
		stream.parentCtx = ctx
		stream.ctx, stream.cancelFunc = context.WithCancel(ctx)
	}
}

// MaxGoroutines sets the maximum number of goroutines allowed in the worker pool.
func MaxGoroutines(limit int) Option {
	return func(stream *Stream) {
		stream.maxGoroutines = limit
	}
}

// defaultPanicHandler is the default panic handler that prints the panic information.
func defaultPanicHandler(pc any) {
	const red = "\u001B[31m"
	log.Printf("[%[1]sERROR%[1]s] %v", red, pc)
}
