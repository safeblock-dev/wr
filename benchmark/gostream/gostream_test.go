package gostream_test

import (
	"fmt"
	"testing"

	"github.com/safeblock-dev/wr/gostream"
	concStream "github.com/sourcegraph/conc/stream"
)

const (
	maxGoroutines = 30
)

func BenchmarkWr(b *testing.B) {
	stream := gostream.New(gostream.MaxGoroutines(maxGoroutines))
	defer stream.Wait()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stream.Go(func() (gostream.Callback, error) {
			msg := fmt.Sprintf("%d", i)

			return func() error {
				fmt.Sprintf("%s", msg)

				return nil
			}, nil
		})
	}
}

func BenchmarkConc(b *testing.B) {
	stream := concStream.New().WithMaxGoroutines(maxGoroutines)
	defer stream.Wait()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stream.Go(func() concStream.Callback {
			msg := fmt.Sprintf("%d", i)

			return func() {
				fmt.Sprintf("%s", msg)
			}
		})
	}
}
