package wr

import (
	"github.com/safeblock-dev/wr/wrgroup"
	"github.com/safeblock-dev/wr/wrpool"
)

func NewWaitingGroup(options ...wrgroup.Option) *wrgroup.WaitGroup {
	return wrgroup.New(options...)
}

func NewPool(options ...wrpool.Option) *wrpool.Pool {
	return wrpool.New(options...)
}
