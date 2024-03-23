package syncgroup

import (
	"fmt"
	"sync"
	"testing"

	"github.com/safeblock-dev/wr/syncgroup"
	"github.com/sourcegraph/conc"
)

const (
	maxGoroutines = 10
)

func BenchmarkDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for j := 0; j < maxGoroutines; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				_ = fmt.Sprintf("%d %d", i, j)
			}()
		}

		wg.Wait()
	}
}

func BenchmarkWr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := syncgroup.New()

		for j := 0; j < maxGoroutines; j++ {
			wg.Go(func() {
				_ = fmt.Sprintf("%d %d", i, j)
			})

		}

		wg.Wait()
	}
}

func BenchmarkConc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := conc.NewWaitGroup()

		for j := 0; j < maxGoroutines; j++ {
			wg.Go(func() {
				_ = fmt.Sprintf("%d %d", i, j)
			})

		}

		wg.Wait()
	}
}
