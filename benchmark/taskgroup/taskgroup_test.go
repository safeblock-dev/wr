package taskgroup_test

import (
	"fmt"
	"testing"

	"github.com/oklog/run"
	"github.com/safeblock-dev/wr/taskgroup"
)

const maxGoroutines = 10

func BenchmarkWr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		task := taskgroup.New()

		for j := 0; j < maxGoroutines; j++ {
			task.Add(func() error {
				_ = fmt.Sprintf("%d", i)

				return nil
			}, func(error) {
				return
			})
		}

		_ = task.Run()
	}
}

func BenchmarkRun(b *testing.B) {
	for i := 0; i < b.N; i++ {
		task := run.Group{}

		for j := 0; j < maxGoroutines; j++ {
			task.Add(func() error {
				_ = fmt.Sprintf("%d", i)

				return nil
			}, func(error) {
				return
			})
		}

		_ = task.Run()
	}
}
