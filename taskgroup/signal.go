package taskgroup

import (
	"errors"
	"os"
)

var (
	Interrupt error = errSignal{os.Interrupt} //nolint:errname,gochecknoglobals // Satisfies the os.Interrupt.
	Kill      error = errSignal{os.Kill}      //nolint:errname,gochecknoglobals // Satisfies the os.Kill.
)

//nolint:errname // local error.
type errSignal struct {
	sig os.Signal
}

func (err errSignal) Error() string {
	return err.sig.String()
}

func (err errSignal) String() string {
	return err.sig.String()
}

func IsSignalError(err error) bool {
	return errors.As(err, new(errSignal))
}
