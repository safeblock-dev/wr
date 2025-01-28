package taskgroup_test

import (
	"errors"
	"testing"

	"github.com/safeblock-dev/werr"
	"github.com/safeblock-dev/wr/taskgroup"
	"github.com/stretchr/testify/require"
)

func TestIsSignalError(t *testing.T) {
	t.Parallel()

	t.Run("when signal error", func(t *testing.T) {
		t.Parallel()

		require.True(t, taskgroup.IsSignalError(taskgroup.Interrupt))
	})

	t.Run("when wrapped signal error", func(t *testing.T) {
		t.Parallel()

		require.True(t, taskgroup.IsSignalError(werr.Wrap(taskgroup.Interrupt)))
	})

	t.Run("when there is no signal error", func(t *testing.T) {
		t.Parallel()

		require.False(t, taskgroup.IsSignalError(errors.New("not a signal error")))
	})
}
