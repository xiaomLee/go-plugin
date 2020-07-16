package lb

import (
	"sync/atomic"
)

// RoundRobin round robin loadBalance impl
type RoundRobin struct {
	ops *uint64
}

// NewRoundRobin create a RoundRobin
func NewRoundRobin() LoadBalance {
	var ops uint64
	ops = 0

	return RoundRobin{
		ops: &ops,
	}
}

// Select select a server from servers using RoundRobin
func (rr RoundRobin) Select(servers []*Node, opts ...string) *Node {
	l := uint64(len(servers))

	if 0 >= l {
		return nil
	}

	target := servers[int(atomic.AddUint64(rr.ops, 1)%l)]
	return target
}
