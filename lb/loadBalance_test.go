package lb

import (
	"strconv"
	"testing"
	"time"

	"github.com/xiaomLee/go-plugin/lb/localRoundRobin"
	"github.com/xiaomLee/go-plugin/lb/roundRobin"
	"github.com/xiaomLee/go-plugin/registry"
	"github.com/xiaomLee/go-plugin/registry/etcd"
)

func TestRemainder(t *testing.T) {
	var reg registry.Registry
	reg = etcd.NewRegistry(registry.Timeout(time.Second*3), registry.Addrs("127.0.0.1:2379"))
	go func() {
		for i := 1; i <= 9; i++ {
			ser := &registry.Service{
				Name:  "test",
				Nodes: []*registry.Node{&registry.Node{Id: strconv.Itoa(i), Address: "127.0.0.1:8" + strconv.Itoa(i)}},
			}
			if err := reg.Register(ser, registry.RegisterTTL(time.Second*10)); err != nil {
				t.Logf(err.Error())
			}
			time.Sleep(time.Second * 10)
		}
	}()

	// lb := &remainder.RemainderLoadBalance{}
	lb := &roundRobin.RoundRobinLoadBalance{}
	lb.SetServiceName("test")
	lb.SetRegistry(reg)
	lb.SetReloadFunc(func() error {
		println("触发reload")
		return nil
	})

	time.Sleep(time.Second)

	if err := lb.Start(time.Second * 10); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		ser := lb.GetNode(strconv.Itoa(i))
		if ser != nil {
			println("get server:", i, ser.Id, ser.Address)
		}
		time.Sleep(time.Second)
	}
}

func TestLocalRoundRobin(t *testing.T) {
	lb := &localRoundRobin.RoundRobin{}
	lb.SetServiceName("test")
	nodes := make([]*registry.Node, 0)
	nodes = append(nodes, &registry.Node{
		Id:       "1",
		Address:  "8080",
		Metadata: nil,
	})
	nodes = append(nodes, &registry.Node{
		Id:       "2",
		Address:  "8090",
		Metadata: nil,
	})
	lb.SetEndPoints(nodes)
	lb.Start(0)

	for i := 0; i < 100; i++ {
		ser := lb.GetNode(strconv.Itoa(i))
		if ser != nil {
			println("get server:", i, ser.Id, ser.Address)
		}
		time.Sleep(time.Second)
	}
}
