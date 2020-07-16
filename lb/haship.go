package lb

import (
	"hash/fnv"
)

// HashBalance is hash IP loadBalance impl
type HashBalance struct {
}

// NewHashBalance create a HashBalance
func NewHashBalance() LoadBalance {
	lb := HashBalance{}
	return lb
}

// Select select a server from servers using HashBalance
func (b HashBalance) Select(servers []*Node, opts ...string) *Node {
	l := len(servers)
	if 0 >= l {
		return nil
	}

	if len(opts) == 0 {
		return nil
	}

	hash := fnv.New32a()
	// key is client ip
	key := opts[0]
	hash.Write([]byte(key))
	node := servers[hash.Sum32()%uint32(l)]
	return node
}
