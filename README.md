# WR

`wr` is a Go library that provides convenient wrappers and utilities for concurrent programming, including a custom `WaitGroup` implementation with panic handling and a flexible worker pool.

## Features

- **wrgroup**: A wrapper around `sync.WaitGroup` with a custom panic handler.
- **wrpool**: A worker pool for managing and executing tasks concurrently with optional error and panic handling.
- **wrtask**: Package provides a TaskGroup structure for running and managing multiple concurrent tasks with support for context and interruption.

## Installation

```sh
go get github.com/safeblock-dev/wr
```

### Examples:

- [WaitingGroup](example/waitgroup/main.go)
- [Pool](example/pool/main.go)

#### WaitingGroup

```go
wg := wr.NewWaitingGroup()
wg.Go(func() {
    // Your task logic here
    log.Println("Task executed")
})
wg.Wait()
```

#### Pool

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

pool := wr.NewPool(
    wrpool.Context(ctx),
    wrpool.MaxGoroutines(5),
)
defer pool.Wait()

pool.Go(func() error {
    // Your task logic here
    log.Println("Task executed")
    return nil
})
```

#### Tasks

```go
group := wrtask.New()

group.Add(wrtask.SignalHandler(context.TODO(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM))
log.Println("We're waiting for 5 seconds, giving you an opportunity to gracefully exit the program.")

// Actor 1: performs a long operation
group.AddContext(func(ctx context.Context) error {
	log.Println("Actor 1 working...")
	<-ctx.Done()
	log.Println("Actor 1 stopped")
	return nil
}, func(context.Context, error) {
	log.Println("Actor 1 interrupted")
})

// Actor 2: returns an error
group.Add(func() error {
	log.Println("Actor 2 working...")
	time.Sleep(5 * time.Second)
	log.Println("Actor 2 stopped")
	return errors.New("error in actor 2")
}, func(error) {
	log.Println("Actor 2 interrupted")
})

// Run all actors and wait for their completion
if err := group.Run(); err != nil {
	log.Println("Error in group:", err)
}
```

### Benchmark

#### Pool

- [pond](github.com/alitto/pond)
- [gopool](github.com/devchat-ai/gopool)
- [ants/v2](github.com/panjf2000/ants/v2)
- [conc](github.com/sourcegraph/conc)

| Benchmark                | Iterations | Time (ns/op) | Memory (B/op) | Allocations (allocs/op) |
|--------------------------|------------|--------------|---------------|-------------------------|
| BenchmarkWr-8            | 17103136   | 365.3        | 33            | 3                       |
| BenchmarkConcErrorPool-8 | 17164700   | 397.1        | 66            | 4                       |
| BenchmarkGopool-8        | 18526196   | 345.1        | 33            | 3                       |
| BenchmarkAnts-8          | 17683563   | 337.2        | 40            | 3                       |
| BenchmarkPond-8          | 6861000    | 1232         | 32            | 3                       |


