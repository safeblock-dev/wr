package panics

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

// NewRecovered creates a panics.Recovered from a panic value and a collected
// stacktrace. The skip parameter allows the caller to skip stack frames when
// collecting the stacktrace. Calling with a skip of 0 means include the call to
// NewRecovered in the stacktrace.
func NewRecovered(skip int, value any) Recovered {
	// 64 frames should be plenty
	var callers [64]uintptr
	n := runtime.Callers(skip+1, callers[:])

	return Recovered{
		Value:   value,
		Callers: callers[:n],
		Stack:   debug.Stack(),
	}
}

// Recovered is a panic that was caught with recover().
type Recovered struct {
	// The original value of the panic.
	Value any
	// The caller list as returned by runtime.Callers when the panic was
	// recovered. Can be used to produce a more detailed stack information with
	// runtime.CallersFrames.
	Callers []uintptr
	// The formatted stacktrace from the goroutine where the panic was recovered.
	// Easier to use than Callers.
	Stack []byte
}

// String renders a human-readable formatting of the panic.
func (p Recovered) String() string {
	return fmt.Sprintf("panic: %v\nstacktrace:\n%s\n", p.Value, p.Stack)
}

// AsError casts the panic into an error implementation. The implementation
// is unwrappable with the cause of the panic, if the panic was provided one.
func (p Recovered) AsError() error {
	return &RecoveredError{p}
}

// RecoveredError wraps a panics.Recovered in an error implementation.
type RecoveredError struct{ Recovered }

var _ error = (*RecoveredError)(nil)

func (p *RecoveredError) Error() string { return p.String() }

func (p *RecoveredError) Unwrap() error {
	if err, ok := p.Value.(error); ok {
		return err
	}

	return nil
}
