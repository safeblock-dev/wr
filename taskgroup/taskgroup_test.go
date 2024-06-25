package taskgroup_test

import (
	"context"
	"errors"
	"testing"

	"github.com/safeblock-dev/wr/taskgroup"
	"github.com/stretchr/testify/require"
)

func TestTaskGroup_Add(t *testing.T) {
	t.Parallel()

	t.Run("add nil execute function", func(t *testing.T) {
		t.Parallel()

		tg := taskgroup.New()

		require.Panics(t, func() {
			tg.Add(nil, taskgroup.SkipInterrupt())
		}, "expected panic when adding an actor with a nil execute function")
	})

	t.Run("add nil interrupt function", func(t *testing.T) {
		t.Parallel()

		tg := taskgroup.New()

		require.Panics(t, func() {
			tg.Add(func() error { return nil }, nil)
		}, "expected panic when adding an actor with a nil interrupt function")
	})
}

func TestTaskGroup_AddContext(t *testing.T) {
	t.Parallel()

	t.Run("should complete successfully with context", func(t *testing.T) {
		t.Parallel()

		tg := taskgroup.New()

		var (
			interrupted bool
			syncCh      = make(chan struct{})
		)

		tg.Add(func() error {
			<-syncCh
			return nil
		}, taskgroup.SkipInterrupt())

		tg.AddContext(func(ctx context.Context) error {
			syncCh <- struct{}{}
			<-ctx.Done()
			return nil
		}, func(_ context.Context, _ error) {
			interrupted = true
		})

		err := tg.Run()
		require.NoError(t, err)
		require.True(t, interrupted)
	})

	t.Run("should interrupt task on error", func(t *testing.T) {
		t.Parallel()

		tg := taskgroup.New()

		var interrupted bool

		tg.AddContext(func(_ context.Context) error {
			return errors.New("task error")
		}, func(_ context.Context, _ error) {
			interrupted = true
		})

		err := tg.Run()
		require.Error(t, err)
		require.True(t, interrupted)
	})

	t.Run("add nil execute function", func(t *testing.T) {
		t.Parallel()

		tg := taskgroup.New()

		require.Panics(t, func() {
			tg.AddContext(nil, taskgroup.SkipInterruptCtx())
		}, "expected panic when adding an actor with a nil execute function")
	})

	t.Run("add nil interrupt function", func(t *testing.T) {
		t.Parallel()

		tg := taskgroup.New()

		require.Panics(t, func() {
			tg.AddContext(func(_ context.Context) error { return nil }, nil)
		}, "expected panic when adding an actor with a nil interrupt function")
	})
}

func TestTaskGroup_Run(t *testing.T) {
	t.Parallel()

	t.Run("should complete successfully with no tasks", func(t *testing.T) {
		t.Parallel()

		tg := taskgroup.New()
		err := tg.Run()
		require.NoError(t, err)
	})

	t.Run("should return error from task", func(t *testing.T) {
		t.Parallel()

		tg := taskgroup.New()
		expectedErr := errors.New("task error")

		tg.Add(func() error {
			return expectedErr
		}, taskgroup.SkipInterrupt())

		err := tg.Run()
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("should interrupt other tasks when one returns", func(t *testing.T) {
		t.Parallel()

		tg := taskgroup.New()
		var interrupted bool

		tg.Add(func() error {
			return nil
		}, func(error) {
			interrupted = true
		})

		tg.Add(func() error {
			return errors.New("returning early")
		}, taskgroup.SkipInterrupt())

		err := tg.Run()
		require.Error(t, err)
		require.True(t, interrupted)
	})

	t.Run("should handle panic in task", func(t *testing.T) {
		t.Parallel()

		tg := taskgroup.New()
		panicValue := "test panic"

		tg.Add(func() error {
			panic(panicValue)
		}, taskgroup.SkipInterrupt())

		err := tg.Run()
		require.Error(t, err)
		require.Contains(t, err.Error(), panicValue)
	})

	t.Run("should allow restarting group after Run", func(t *testing.T) {
		t.Parallel()

		var result [3]bool

		tg := taskgroup.New()
		tg.Add(func() error {
			result[0] = !result[0]
			return nil
		}, taskgroup.SkipInterrupt())

		tg.Add(func() error {
			result[1] = !result[1]
			return nil
		}, taskgroup.SkipInterrupt())

		require.NoError(t, tg.Run())
		require.Equal(t, [3]bool{true, true, false}, result)

		tg.Add(func() error {
			result[2] = !result[2]
			return nil
		}, taskgroup.SkipInterrupt())

		require.NoError(t, tg.Run())
		require.Equal(t, [3]bool{false, false, true}, result)
	})
}

func TestTaskGroup_Size(t *testing.T) {
	t.Parallel()

	tg := taskgroup.New()
	tg.Add(func() error { return nil }, taskgroup.SkipInterrupt())

	require.Equal(t, 1, tg.Size())
}
