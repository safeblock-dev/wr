package gopool_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/alitto/pond"
	devchatGopool "github.com/devchat-ai/gopool"
	"github.com/panjf2000/ants/v2"
	"github.com/safeblock-dev/wr/gopool"
	concPool "github.com/sourcegraph/conc/pool"
)

const (
	maxGoroutines = 30
)

func BenchmarkDefault(b *testing.B) {
	var wg sync.WaitGroup
	defer wg.Wait()

	goroutineLimit := make(chan struct{}, maxGoroutines)
	errorHandler := func(_ error) {}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		goroutineLimit <- struct{}{}
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			defer func() {
				if pc := recover(); pc != nil {
					fmt.Println(pc)
				}
			}()

			msg := fmt.Sprintf("%d", i)
			if i%10 == 0 {
				errorHandler(errors.New(msg))
			}

			<-goroutineLimit
		}(i)
	}
}

func BenchmarkWr(b *testing.B) {
	pool := gopool.New(gopool.MaxGoroutines(maxGoroutines))
	defer pool.Wait()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.Go(func() error {
			msg := fmt.Sprintf("%d", i)
			if i%10 == 0 {
				return errors.New(msg)
			} else {
				return nil
			}
		})
	}
}

func BenchmarkConc(b *testing.B) {
	pool := concPool.New().WithMaxGoroutines(maxGoroutines).WithErrors()
	defer pool.Wait()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.Go(func() error {
			msg := fmt.Sprintf("%d", i)
			if i%10 == 0 {
				return errors.New(msg)
			} else {
				return nil
			}
		})
	}
}

func BenchmarkGopool(b *testing.B) {
	pool := devchatGopool.NewGoPool(maxGoroutines)
	defer pool.Release()
	defer pool.Wait()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.AddTask(func() (interface{}, error) {
			msg := fmt.Sprintf("%d", i)
			if i%10 == 0 {
				return nil, errors.New(msg)
			} else {
				return nil, nil
			}
		})
	}
}

func BenchmarkAnts(b *testing.B) {
	var wg sync.WaitGroup
	p, _ := ants.NewPool(maxGoroutines)
	defer p.Release()
	defer wg.Wait()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		_ = p.Submit(func() {
			defer wg.Done()

			msg := fmt.Sprintf("%d", i)
			if i%10 == 0 {
				errors.New(msg)
			}
		})
	}
}

func BenchmarkPond(b *testing.B) {
	pool := pond.New(maxGoroutines, 0, pond.MinWorkers(maxGoroutines))
	defer pool.StopAndWait()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.Submit(func() {
			msg := fmt.Sprintf("%d", i)
			if i%10 == 0 {
				errors.New(msg)
			}
		})
	}
}
