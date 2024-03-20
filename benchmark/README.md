### Benchmark

| Benchmark                | Iterations | Time (ns/op) | Memory (B/op) | Allocations (allocs/op) |
|--------------------------|------------|--------------|---------------|-------------------------|
| BenchmarkWr-8            | 17103136   | 365.3        | 33            | 3                       |
| BenchmarkConcErrorPool-8 | 17164700   | 397.1        | 66            | 4                       |
| BenchmarkGopool-8        | 18526196   | 345.1        | 33            | 3                       |
| BenchmarkAnts-8          | 17683563   | 337.2        | 40            | 3                       |
| BenchmarkPond-8          | 6861000    | 1232         | 32            | 3                       |

```go
go test -bench =.-benchmem -benchtime =5s
```