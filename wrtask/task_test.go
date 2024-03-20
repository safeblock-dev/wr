package wrtask_test

import (
	"context"
	"errors"
	"testing"

	"github.com/safeblock-dev/wr/wrtask"
	"github.com/stretchr/testify/require"
)

func TestTaskGroup_AddContext(t *testing.T) {
	t.Parallel()

	t.Run("should complete successfully with context", func(t *testing.T) {
		t.Parallel()

		group := wrtask.New()
		var (
			interrupted = false
			syncer      = make(chan struct{})
		)

		group.Add(func() error {
			<-syncer

			return nil
		}, wrtask.SkipInterrupt())
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

		group := wrtask.New()
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
}

func TestTaskGroup_Run(t *testing.T) {
	t.Parallel()

	t.Run("should complete successfully with no tasks", func(t *testing.T) {
		t.Parallel()

		group := wrtask.New()
		err := group.Run()
		require.NoError(t, err)
	})

	t.Run("should return error from task", func(t *testing.T) {
		t.Parallel()

		group := wrtask.New()
		expectedErr := errors.New("task error")

		group.Add(func() error {
			return expectedErr
		}, wrtask.SkipInterrupt())

		err := group.Run()
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("should interrupt other tasks when one returns", func(t *testing.T) {
		t.Parallel()

		group := wrtask.New()
		interrupted := false

		group.Add(func() error {
			return nil
		}, func(error) {
			interrupted = true
		})

		group.Add(func() error {
			return errors.New("returning early")
		}, wrtask.SkipInterrupt())

		err := group.Run()
		require.Error(t, err)
		require.True(t, interrupted)
	})

	t.Run("should handle panic in task", func(t *testing.T) {
		t.Parallel()

		group := wrtask.New()
		panicValue := "test panic"

		group.Add(func() error {
			panic(panicValue)
		}, wrtask.SkipInterrupt())

		err := group.Run()
		require.Error(t, err)
		require.Contains(t, err.Error(), panicValue)
	})

	t.Run("should allow restarting group after Run", func(t *testing.T) {
		t.Parallel()

		var (
			result [3]bool
			syncer = make(chan struct{})
		)

		group := wrtask.New()
		group.Add(func() error {
			result[0] = !result[0]
			syncer <- struct{}{}

			return errors.New("test error")
		}, wrtask.SkipInterrupt())
		group.Add(func() error {
			result[1] = !result[1]
			<-syncer

			return nil
		}, wrtask.SkipInterrupt())

		require.Error(t, group.Run())
		require.Equal(t, [3]bool{true, true, false}, result)

		group.Add(func() error {
			result[2] = !result[2]

			return nil
		}, wrtask.SkipInterrupt())

		require.NoError(t, group.Run())
		require.Equal(t, [3]bool{false, false, true}, result)
	})
}
