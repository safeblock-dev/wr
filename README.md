# WR

`wr` is a Go library that provides convenient wrappers and utilities for concurrent programming, including a custom `WaitGroup` implementation with panic handling and a flexible worker pool.

## Features

- **wrgroup**: A wrapper around `sync.WaitGroup` with a custom panic handler.
- **wrpool**: A worker pool for managing and executing tasks concurrently with optional error and panic handling.

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

### Benchmark

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


