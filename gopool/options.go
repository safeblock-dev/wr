package gopool

import (
	"context"
	"log"
)

// Option represents an option that can be passed when instantiating a Pool to customize it.
type Option func(pool *Pool)

// Context sets a parent context for the pool to stop all workers when it is cancelled.
func Context(ctx context.Context) Option {
	return func(pool *Pool) {
		pool.parentCtx = ctx
		pool.ctx, pool.cancelFunc = context.WithCancel(ctx)
	}
}

// PanicHandler sets the panic handler function for the pool.
// It allows customizing how panics are handled within the pool.
func PanicHandler(panicHandler func(any)) Option {
	return func(pool *Pool) {
		pool.panicHandler = panicHandler
	}
}

// ErrorHandler sets the error handler function for the pool.
// It allows customizing how errors are handled within the pool.
func ErrorHandler(errorHandler func(err error)) Option {
	return func(pool *Pool) {
		pool.errorHandler = errorHandler
	}
}

// MaxGoroutines sets the maximum number of goroutines allowed in the pool.
// It limits the number of concurrent tasks that can be executed simultaneously.
func MaxGoroutines(limit int) Option {
	return func(pool *Pool) {
		if pool.limiter != nil {
			pool.limiter.close()
		}
		switch {
		case limit < 1:
			pool.limiter = nil
		default:
			pool.limiter = make(limiter, limit)
		}
	}
}

// defaultPanicHandler is the default panic handler that prints the panic information.
// It logs the panic message with a distinctive error formatting.
func defaultPanicHandler(pc any) {
	const red = "\u001B[31m"
	log.Printf("[%[1]sERROR%[1]s] %v", red, pc)
}
