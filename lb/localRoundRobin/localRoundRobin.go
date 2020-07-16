package localRoundRobin

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xiaomLee/go-plugin/registry"
)

// 轮训模式负载均衡实现
type RoundRobin struct {
	ServiceName string
	ReadyFlag   int32
	CloseFlag   int32
	ReloadFunc  func() error
	Nodes       []*registry.Node

	index int32

	sync.RWMutex
}

func (b *RoundRobin) SetServiceName(name string) {
	b.ServiceName = name
}

func (b *RoundRobin) SetEndPoints(nodes []*registry.Node) {
	b.Lock()
	defer b.Unlock()
	b.Nodes = nodes
}

func (b *RoundRobin) SetReloadFunc(f func() error) {
	b.ReloadFunc = f
}

func (b *RoundRobin) SetRegistry(reg registry.Registry) {
	panic("not implement")
}

func (b *RoundRobin) Ready() bool {
	return atomic.LoadInt32(&b.ReadyFlag) == 1
}

func (b *RoundRobin) GetNode(key string) *registry.Node {
	if atomic.LoadInt32(&b.CloseFlag) == 1 {
		return nil
	}

	i := atomic.AddInt32(&b.index, 1)

	b.RLock()
	defer b.RUnlock()

	i = i % int32(len(b.Nodes))
	return b.Nodes[int(i)]
}

func (b *RoundRobin) GetNodes() []*registry.Node {
	if atomic.LoadInt32(&b.CloseFlag) == 1 {
		return nil
	}

	b.RLock()
	defer b.RUnlock()

	return b.Nodes
}

func (b *RoundRobin) Start(TTL time.Duration) error {
	if b.ServiceName == "" {
		errors.New("serviceName empty")
	}

	if atomic.LoadInt32(&b.CloseFlag) == 1 {
		return nil
	}

	atomic.StoreInt32(&b.ReadyFlag, 1)

	return nil
}

func (b *RoundRobin) Close() {
	atomic.StoreInt32(&b.CloseFlag, 1)
}
