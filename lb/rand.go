package lb

import (
	"math/rand"
	"time"
)

// RandBalance is rand loadBalance impl
type RandBalance struct {
}

// NewRandBalance create a RandBalance
func NewRandBalance() LoadBalance {
	rand.Seed(time.Now().UnixNano())
	lb := RandBalance{}
	return lb
}

// Select select a server from servers using rand
func (rb RandBalance) Select(servers []*Node, opts ...string) *Node {
	l := len(servers)
	if 0 >= l {
		return nil
	}
	node := servers[rand.Intn(l)]
	return node
}
