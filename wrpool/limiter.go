package wrpool

// limiter is a simple channel-based limiter for controlling the number of goroutines.
type limiter chan struct{}

func (l limiter) release() {
	if l != nil {
		<-l
	}
}

func (l limiter) limit() int {
	return cap(l)
}

func (l limiter) close() {
	if l != nil {
		close(l)
	}
}
