package syncgroup

import (
	"log"
)

// Option represents an option that can be passed when instantiating a WaitGroup to customize it.
type Option func(wg *WaitGroup)

// PanicHandler allows changing the panic handler function of a WaitGroup.
// It accepts a function that handles a panic and assigns it to WaitGroup's panicHandler field.
func PanicHandler(panicHandler func(pc any)) Option {
	return func(wg *WaitGroup) {
		wg.panicHandler = panicHandler
	}
}

// defaultPanicHandler is the default panic handler that prints the panic information to the log.
// It uses ANSI escape sequences to colorize the error message in red.
func defaultPanicHandler(pc any) {
	const red = "\u001B[31m"
	log.Printf("[%[1]sERROR%[1]s] %v", red, pc)
}
