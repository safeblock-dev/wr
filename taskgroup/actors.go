package taskgroup

import (
	"context"
	"errors"
	"os"
	"os/signal"
)

var ErrSignal = errors.New("signal error")

// ContextHandler returns an actor, i.e. an execute and interrupt func, that
// terminates when the provided context is canceled.
func ContextHandler(ctx context.Context) (ExecuteFn, InterruptFn) {
	ctx, cancel := context.WithCancel(ctx)

	return func() error {
			<-ctx.Done()

			return ctx.Err()
		}, func(error) {
			cancel()
		}
}

// SignalHandler creates an actor that terminates when a signal is received or the context is canceled.
func SignalHandler(ctx context.Context, signals ...os.Signal) (ExecuteFn, InterruptFn) {
	ctx, cancel := context.WithCancel(ctx)

	execute := func() error {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, signals...)
		defer signal.Stop(sig)
		select {
		case <-sig:
			return ErrSignal
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	interrupt := func(_ error) {
		cancel()
	}

	return execute, interrupt
}
