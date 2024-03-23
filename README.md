# WR

### GoPool

**gopool** is a Go library that provides a goroutine pool for efficient task execution. It allows you to manage a fixed number of goroutines to execute tasks concurrently, reducing the overhead of creating and destroying goroutines frequently.

### SyncGroup
**syncgroup** enhances Go's sync.WaitGroup with additional features such as panic recovery and custom panic handlers. It simplifies the management of goroutines, ensuring that any panics are handled gracefully and that all goroutines complete before proceeding.

### TaskGroup
**taskgroup** offers a way to group related tasks and manage their execution as a unit. It supports interruptible tasks and automatic error propagation, making it easier to handle errors and control the flow of concurrent tasks.

### GoStream
**gostream** provides a framework for executing tasks concurrently while ensuring that their callbacks are executed sequentially. This is useful for maintaining consistency when processing results from multiple concurrent operations.

## Installation

```sh
go get github.com/safeblock-dev/wr
```

### Examples:

- [GoPool](example/gopool/main.go)
- [GoPool (one fail)](example/gopool_one_fail/main.go)
- [GoStream](example/gostream/main.go)
- [SyncGroup](example/syncgroup/main.go)
- [TaskGroup](example/taskgroup/main.go)

## Benchmark Results

System Software Overview:

    Device: macbook air, m1
    System Version: macOS 14.1 (23B74)
    Kernel Version: Darwin 23.1.0


Run:

```sh
make bench
```

or 

```sh
cd benchmark && go test ./... -bench . -benchtime 5s -timeout 0 -run=XXX -cpu 1 -benchmem
```

- [pond](github.com/alitto/pond)
- [gopool](github.com/devchat-ai/gopool)
- [ants/v2](github.com/panjf2000/ants/v2)
- [conc](github.com/sourcegraph/conc)

## Pool

| Benchmark       | Iterations        | Time (ns/op)      | Memory (B/op)    | Allocations (allocs/op) |
|-----------------|-------------------|-------------------|------------------|-------------------------|
| Default         | 14,130,334        | 400.6             | 88               | 4                       |
| Wr              | 20,005,642        | 300.8             | 33               | 3                       |
| Conc            | 16,104,584        | 424.3             | 66               | 4                       |
| Gopool          | 13,473,462        | 437.8             | 33               | 3                       |
| Ants            | 16,403,790        | 366.7             | 40               | 3                       |
| Pond            | 20,494,136        | 299.3             | 32               | 3                       |

## Stream

| Benchmark       | Iterations        | Time (ns/op)      | Memory (B/op)    | Allocations (allocs/op) |
|-----------------|-------------------|-------------------|------------------|-------------------------|
| Wr              | 9,830,470         | 543.8             | 104              | 7                       |
| Conc            | 12,192,182        | 493.5             | 104              | 7                       |

## SyncGroup

| Benchmark       | Iterations        | Time (ns/op)      | Memory (B/op)    | Allocations (allocs/op) |
|-----------------|-------------------|-------------------|------------------|-------------------------|
| Default         | 2,299,267         | 2634              | 549              | 31                      |
| Wr              | 2,046,895         | 2958              | 712              | 41                      |
| Conc            | 2,028,130         | 2994              | 712              | 41                      |

## TaskGroup

| Benchmark       | Iterations        | Time (ns/op)      | Memory (B/op)    | Allocations (allocs/op) |
|-----------------|-------------------|-------------------|------------------|-------------------------|
| Wr              | 1,714,483         | 3492              | 1680             | 59                      |
| Run             | 1,512,882         | 3975              | 1560             | 57                      |