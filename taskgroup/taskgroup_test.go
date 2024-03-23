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

		tasks := taskgroup.New()

		require.Panics(t, func() {
			tasks.Add(nil, taskgroup.SkipInterrupt())
		}, "expected panic when adding an actor with a nil execute function")
	})

	t.Run("add nil interrupt function", func(t *testing.T) {
		t.Parallel()

		tasks := taskgroup.New()

		require.Panics(t, func() {
			tasks.Add(func() error { return nil }, nil)
		}, "expected panic when adding an actor with a nil interrupt function")
	})
}

func TestTaskGroup_AddContext(t *testing.T) {
	t.Parallel()

	t.Run("should complete successfully with context", func(t *testing.T) {
		t.Parallel()

		group := taskgroup.New()
		var (
			interrupted = false
			syncer      = make(chan struct{})
		)

		group.Add(func() error {
			<-syncer

			return nil
		}, taskgroup.SkipInterrupt())
		group.AddContext(func(ctx context.Context) error {
			syncer <- struct{}{}
			<-ctx.Done()

			return nil
		}, func(_ context.Context, _ error) {
			interrupted = true
		})

		err := group.Run()
		require.NoError(t, err)
		require.True(t, interrupted)
	})

	t.Run("should interrupt task on error", func(t *testing.T) {
		t.Parallel()

		group := taskgroup.New()
		var interrupted bool

		group.AddContext(func(_ context.Context) error {
			return errors.New("task error")
		}, func(_ context.Context, _ error) {
			interrupted = true
		})

		err := group.Run()
		require.Error(t, err)
		require.True(t, interrupted)
	})

	t.Run("add nil execute function", func(t *testing.T) {
		t.Parallel()

		tasks := taskgroup.New()

		require.Panics(t, func() {
			tasks.AddContext(nil, func(_ context.Context, _ error) {})
		}, "expected panic when adding an actor with a nil execute function")
	})

	t.Run("add nil interrupt function", func(t *testing.T) {
		t.Parallel()

		tasks := taskgroup.New()

		require.Panics(t, func() {
			tasks.AddContext(func(_ context.Context) error { return nil }, nil)
		}, "expected panic when adding an actor with a nil interrupt function")
	})
}

func TestTaskGroup_Run(t *testing.T) {
	t.Parallel()

	t.Run("should complete successfully with no tasks", func(t *testing.T) {
		t.Parallel()

		group := taskgroup.New()
		err := group.Run()
		require.NoError(t, err)
	})

	t.Run("should return error from task", func(t *testing.T) {
		t.Parallel()

		group := taskgroup.New()
		expectedErr := errors.New("task error")

		group.Add(func() error {
			return expectedErr
		}, taskgroup.SkipInterrupt())

		err := group.Run()
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("should interrupt other tasks when one returns", func(t *testing.T) {
		t.Parallel()

		group := taskgroup.New()
		interrupted := false

		group.Add(func() error {
			return nil
		}, func(error) {
			interrupted = true
		})

		group.Add(func() error {
			return errors.New("returning early")
		}, taskgroup.SkipInterrupt())

		err := group.Run()
		require.Error(t, err)
		require.True(t, interrupted)
	})

	t.Run("should handle panic in task", func(t *testing.T) {
		t.Parallel()

		group := taskgroup.New()
		panicValue := "test panic"

		group.Add(func() error {
			panic(panicValue)
		}, taskgroup.SkipInterrupt())

		err := group.Run()
		require.Error(t, err)
		require.Contains(t, err.Error(), panicValue)
	})

	t.Run("should allow restarting group after Run", func(t *testing.T) {
		t.Parallel()

		var result [3]bool

		group := taskgroup.New()
		group.Add(func() error {
			result[0] = !result[0]

			return nil
		}, taskgroup.SkipInterrupt())
		group.Add(func() error {
			result[1] = !result[1]

			return nil
		}, taskgroup.SkipInterrupt())

		require.NoError(t, group.Run())
		require.Equal(t, [3]bool{true, true, false}, result)

		group.Add(func() error {
			result[2] = !result[2]

			return nil
		}, taskgroup.SkipInterrupt())

		require.NoError(t, group.Run())
		require.Equal(t, [3]bool{false, false, true}, result)
	})
}

func TestTaskGroup_Size(t *testing.T) {
	t.Parallel()

	tasks := taskgroup.New()
	tasks.Add(func() error { return nil }, taskgroup.SkipInterrupt())

	require.Equal(t, 1, tasks.Size())
}
