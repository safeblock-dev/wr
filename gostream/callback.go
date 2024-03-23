package gostream

import (
	"sync"

	"github.com/safeblock-dev/wr/panics"
)

// callbackData represents data associated with a callback, including the callback function and any error.
type callbackData struct {
	fn    func() error // The callback function to execute.
	err   error        // Any error that occurred during task execution.
	panic *panics.Recovered
}

// callbackChannel is a channel for sending callbackData.
type callbackChannel chan callbackData

// callbackChPool is a pool of callbackChannels to reduce allocations.
var callbackChPool = sync.Pool{ //nolint:gochecknoglobals // optimization
	New: func() any {
		return make(callbackChannel, 1) // Buffer size of 1 to prevent blocking.
	},
}

// getCallbackChannel retrieves a callbackChannel from the pool or creates a new one if necessary.
func getCallbackChannel() callbackChannel {
	return callbackChPool.Get().(callbackChannel)
}

// putCallbackChannel returns a callbackChannel to the pool for reuse.
func putCallbackChannel(ch callbackChannel) {
	callbackChPool.Put(ch)
}
