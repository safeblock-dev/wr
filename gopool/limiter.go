package gopool

// limiter is a simple channel-based limiter for controlling the number of goroutines.
type limiter chan struct{}

// release releases a permit from the limiter.
func (l limiter) release() {
	if l != nil {
		<-l
	}
}

// limit returns the capacity of the limiter.
func (l limiter) limit() int {
	return cap(l)
}

// close closes the limiter channel.
func (l limiter) close() {
	if l != nil {
		close(l)
	}
}
