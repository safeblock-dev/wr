package wrpool

type limiter chan struct{}

func (l limiter) release() {
	if l != nil {
		<-l
	}
}

func (l limiter) close() {
	if l != nil {
		close(l)
	}
}
