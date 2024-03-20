package wrtask

import (
	"context"
	"errors"
	"os"
	"os/signal"
)

var ErrSignal = errors.New("signal error")

// ContextHandler creates an actor that terminates when the provided context is canceled.
func ContextHandler(ctx context.Context) (func() error, func(error)) {
	execute := func() error {
		<-ctx.Done()

		return ctx.Err()
	}

	return execute, SkipInterrupt()
}

// SignalHandler creates an actor that terminates when a signal is received or the context is canceled.
func SignalHandler(ctx context.Context, signals ...os.Signal) (func() error, func(error)) {
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

	return execute, SkipInterrupt()
}
