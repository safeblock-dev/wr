package benchmark_test

import (
	"fmt"
	"testing"

	"github.com/oklog/run"
	"github.com/safeblock-dev/wr/wrtask"
)

const maxTasksGoroutines = 10

func BenchmarkWrTasks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		task := wrtask.New()

		for j := 0; j < maxTasksGoroutines; j++ {
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

func BenchmarkRunTasks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		task := run.Group{}

		for j := 0; j < maxTasksGoroutines; j++ {
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
