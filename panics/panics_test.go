package panics_test

import (
	"errors"
	"testing"

	"github.com/safeblock-dev/wr/panics"
	"github.com/stretchr/testify/require"
)

func TestNewRecovered(t *testing.T) {
	t.Parallel()

	t.Run("should create a Recovered instance", func(t *testing.T) {
		t.Parallel()

		recovered := panics.NewRecovered(0, "test panic")
		require.NotNil(t, recovered)
		require.Equal(t, "test panic", recovered.Value)
		require.NotEmpty(t, recovered.Stack)
	})
}

func TestRecovered_String(t *testing.T) {
	t.Parallel()

	t.Run("should return a formatted string", func(t *testing.T) {
		t.Parallel()

		recovered := panics.NewRecovered(0, "test panic")
		str := recovered.String()
		require.Contains(t, str, "panic: test panic")
		require.Contains(t, str, "stacktrace:")
	})
}

func TestRecovered_AsError(t *testing.T) {
	t.Parallel()

	t.Run("should return an error", func(t *testing.T) {
		t.Parallel()

		recovered := panics.NewRecovered(0, "test panic")
		err := recovered.AsError()
		require.Error(t, err)
		require.Contains(t, err.Error(), "panic: test panic")
	})

	t.Run("should unwrap the original error", func(t *testing.T) {
		t.Parallel()

		originalErr := errors.New("original error")
		recovered := panics.NewRecovered(0, originalErr)
		err := recovered.AsError()
		unwrappedErr := errors.Unwrap(err)
		require.ErrorIs(t, unwrappedErr, originalErr)
	})
}
