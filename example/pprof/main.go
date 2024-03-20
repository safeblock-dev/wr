package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/safeblock-dev/wr/panics"
	"github.com/safeblock-dev/wr/wrpool"
)

const maxGoroutines = 3

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	example(ctx)

	<-ctx.Done()
}

func example(ctx context.Context) {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	pool := wrpool.New(
		wrpool.Context(ctx),
		wrpool.MaxGoroutines(maxGoroutines),
		wrpool.PanicHandler(panicHandler),
		wrpool.ErrorHandler(errorHandler),
	)
	defer pool.Wait()

	for {
		pool.Go(func() error {
			return nil
		})
	}
}

func panicHandler(p panics.Recovered) {
	log.Println("panic:", p.Value)
}

func errorHandler(err error) {
	log.Println("error:", err)
}
