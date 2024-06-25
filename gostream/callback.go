package gostream

import (
	"sync"
)

// callbackData represents data associated with a callback, including the callback function and any error.
type callbackData struct {
	fn  func() error // fn is the callback function to execute.
	err error        // err is any error that occurred during task execution.
}

// callbackChannel is a channel for sending callbackData.
type callbackChannel chan callbackData

// callbackChPool is a pool of callbackChannels to reduce allocations.
// callbackChPool is a global variable for callbackChannel pooling.
var callbackChPool = sync.Pool{ //nolint:gochecknoglobals
	New: func() interface{} {
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
