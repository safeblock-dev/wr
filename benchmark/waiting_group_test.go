package benchmark

import (
	"fmt"
	"testing"

	"github.com/safeblock-dev/wr/syncgroup"
	"github.com/sourcegraph/conc"
)

const (
	maxWGGoroutines = 10
)

func BenchmarkWrSyncGroup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := syncgroup.New()

		for j := 0; j < maxWGGoroutines; j++ {
			wg.Go(func() {
				_ = fmt.Sprintf("%d %d", i, j)
			})

		}

		wg.Wait()
	}
}

func BenchmarkConcGroup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := conc.NewWaitGroup()

		for j := 0; j < maxWGGoroutines; j++ {
			wg.Go(func() {
				_ = fmt.Sprintf("%d %d", i, j)
			})

		}

		wg.Wait()
	}
}
