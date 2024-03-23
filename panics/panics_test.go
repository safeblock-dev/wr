package panics_test

import (
	"errors"
	"fmt"
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

func TestRecoveredError_Error(t *testing.T) {
	t.Parallel()

	// Simulate a panic with a string value.
	recovered := panics.NewRecovered(0, "test panic")
	recoveredErr := recovered.AsError().(*panics.RecoveredError)

	// Check that the Error method returns the expected string.
	expectedErrMsg := fmt.Sprintf("panic: %v\nstacktrace:\n%s\n", recovered.Value, recovered.Stack)
	require.Equal(t, expectedErrMsg, recoveredErr.Error(), "Error message should match the recovered panic information")
}

func TestRecoveredError_Unwrap(t *testing.T) {
	t.Parallel()

	// Simulate a panic with an error value.
	originalErr := errors.New("original error")
	recovered := panics.NewRecovered(0, originalErr)
	recoveredErr := recovered.AsError().(*panics.RecoveredError)

	// Check that the Unwrap method returns the original error.
	unwrappedErr := errors.Unwrap(recoveredErr)
	require.Equal(t, originalErr, unwrappedErr, "Unwrap should return the original error")
}

func TestRecoveredError_Unwrap_NonError(t *testing.T) {
	t.Parallel()

	// Simulate a panic with a non-error value.
	recovered := panics.NewRecovered(0, "test panic")
	recoveredErr := recovered.AsError().(*panics.RecoveredError)

	// Check that the Unwrap method returns nil for non-error values.
	unwrappedErr := errors.Unwrap(recoveredErr)
	require.Nil(t, unwrappedErr, "Unwrap should return nil for non-error panic values")
}
