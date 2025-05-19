package taskgroup

import (
	"errors"
	"os"
)

var (
	Interrupt error = signalError{os.Interrupt} //nolint:errname,gochecknoglobals // Satisfies the os.Interrupt.
	Kill      error = signalError{os.Kill}      //nolint:errname,gochecknoglobals // Satisfies the os.Kill.
)

var _ os.Signal = (*signalError)(nil)

type signalError struct {
	sig os.Signal
}

func (err signalError) String() string {
	return err.sig.String()
}

func (err signalError) Signal() {}

func (err signalError) Error() string {
	return err.sig.String()
}

func IsSignalError(err error) bool {
	return errors.As(err, new(signalError))
}
