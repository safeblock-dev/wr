# WR

### GoPool

**gopool** is a Go library that provides a goroutine pool for efficient task execution. It allows you to manage a fixed number of goroutines to execute tasks concurrently, reducing the overhead of creating and destroying goroutines frequently.

### SyncGroup

**syncgroup** enhances Go's sync.WaitGroup with additional features such as panic recovery and custom panic handlers. It simplifies the management of goroutines, ensuring that any panics are handled gracefully and that all goroutines complete before proceeding.

### TaskGroup

**taskgroup** offers a way to group related tasks and manage their execution as a unit. It supports interruptible tasks and automatic error propagation, making it easier to handle errors and control the flow of concurrent tasks.

### GoStream

**gostream** provides a framework for executing tasks concurrently while ensuring that their callbacks are executed sequentially. This is useful for maintaining consistency when processing results from multiple concurrent operations.

### GoStreamCh

**gostreamch** extends gostream by providing error handling capabilities. It allows you to submit tasks that may return errors and provides an error channel to retrieve these errors asynchronously.

### GoPoolCh

**gopoolch** is an extension of gopool that includes custom panic and error handlers. It allows you to manage goroutines efficiently with built-in panic recovery and error handling mechanisms.

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
- [GoStreamCh](example/gostreamch/main.go)
- [GoPoolCh](example/gopoolch/main.go)

## Benchmark Results

System Software Overview:

    Device: MacBook Air, M1
    System Version: macOS 14.1 (23B74)
    Kernel Version: Darwin 23.1.0

### Pool

| Benchmark       | Iterations        | Time (ns/op)      | Memory (B/op)    | Allocations (allocs/op) |
|-----------------|-------------------|-------------------|------------------|-------------------------|
| Default         | 20,630,466        | 285.5             | 88               | 4                       |
| Wr              | 16,181,502        | 353.8             | 33               | 3                       |
| Conc            | 14,231,066        | 370.3             | 66               | 4                       |
| Gopool          | 15,818,833        | 348.5             | 33               | 3                       |
| Ants            | 16,583,290        | 351.8             | 40               | 3                       |
| Pond            | 6,783,289         | 1233              | 32               | 3                       |

### Stream

| Benchmark       | Iterations        | Time (ns/op)      | Memory (B/op)    | Allocations (allocs/op) |
|-----------------|-------------------|-------------------|------------------|-------------------------|
| Wr              | 11,183,976        | 505.8             | 112              | 7                       |
| Conc            | 14,056,011        | 435.8             | 104              | 7                       |

### SyncGroup

| Benchmark       | Iterations        | Time (ns/op)      | Memory (B/op)    | Allocations (allocs/op) |
|-----------------|-------------------|-------------------|------------------|-------------------------|
| Default         | 2,299,267         | 2634              | 549              | 31                      |
| Wr              | 2,046,895         | 2958              | 712              | 41                      |
| Conc            | 2,028,130         | 2994              | 712              | 41                      |

### TaskGroup

| Benchmark       | Iterations        | Time (ns/op)      | Memory (B/op)    | Allocations (allocs/op) |
|-----------------|-------------------|-------------------|------------------|-------------------------|
| Wr              | 1,714,483         | 3492              | 1680             | 59                      |
| Run             | 1,512,882         | 3975              | 1560             | 57                      |
