package wrgroup

import (
	"log"

	"github.com/safeblock-dev/wr/panics"
)

// Option represents an option that can be passed when instantiating a WaitGroup to customize it.
type Option func(wg *WaitGroup)

// defaultPanicHandler is the default panic handler that prints the panic information.
func defaultPanicHandler(recovered panics.Recovered) {
	log.Println(recovered.String())
}

// PanicHandler allows changing the panic handler function of a WaitGroup.
func PanicHandler(panicHandler func(recovered panics.Recovered)) Option {
	return func(wg *WaitGroup) {
		wg.panicHandler = panicHandler
	}
}
