# Benchmark

### Pool

| Benchmark                | Iterations | Time (ns/op) | Memory (B/op) | Allocations (allocs/op) |
|--------------------------|------------|--------------|---------------|-------------------------|
| BenchmarkWrPool-8        | 17103136   | 365.3        | 33            | 3                       |
| BenchmarkConcErrorPool-8 | 17164700   | 397.1        | 66            | 4                       |
| BenchmarkGopool-8        | 18526196   | 345.1        | 33            | 3                       |
| BenchmarkAnts-8          | 17683563   | 337.2        | 40            | 3                       |
| BenchmarkPond-8          | 6861000    | 1232         | 32            | 3                       |

### Task

| Benchmark           | Iterations | Time (ns/op) | Memory (B/op) | Allocations (allocs/op) |
|---------------------|------------|--------------|---------------|-------------------------|
| BenchmarkWrTasks-8  | 1218404    | 4518         | 1680          | 60                      |
| BenchmarkRunTasks-8 | 1000000    | 5349         | 1560          | 58                      |

```sh
go test -bench=. -benchmem -benchtime=5s
```