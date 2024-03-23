package wr

import (
	"github.com/safeblock-dev/wr/wrgroup"
	"github.com/safeblock-dev/wr/wrpool"
	"github.com/safeblock-dev/wr/wrstream"
	"github.com/safeblock-dev/wr/wrtask"
)

func NewWaitingGroup(options ...wrgroup.Option) *wrgroup.WaitGroup {
	return wrgroup.New(options...)
}

func NewPool(options ...wrpool.Option) *wrpool.Pool {
	return wrpool.New(options...)
}

func NewTaskGroup() wrtask.TaskGroup {
	return wrtask.New()
}

func NewStream(options ...wrstream.Option) *wrstream.Stream {
	return wrstream.New(options...)
}
