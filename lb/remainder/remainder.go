package remainder

import (
	"github.com/becent/golang-common/loadBalance/roundRobin"
	"github.com/xiaomLee/go-plugin/registry"
	"github.com/mitchellh/hashstructure"
)

// 取余模式负载均衡实现
type RemainderLoadBalance struct {
	roundRobin.RoundRobinLoadBalance
}

func (b *RemainderLoadBalance) GetService(key string) *registry.Node {
	id, err := hashstructure.Hash(key, nil)
	if err != nil {
		return nil
	}

	b.RLock()
	defer b.RUnlock()

	if len(b.Nodes) == 0 {
		return nil
	}

	return b.Nodes[int(id%uint64(len(b.Nodes)))]
}
