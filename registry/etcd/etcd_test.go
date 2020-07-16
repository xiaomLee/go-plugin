package etcd

import (
	"github.com/xiaomLee/go-plugin/registry"
	"testing"
	"time"
)

func TestRegistry(t *testing.T) {
	var register registry.Registry
	register = NewRegistry()
	if err := register.Init(registry.Timeout(3*time.Second), registry.Addrs("127.0.0.1:2379")); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	if err := register.Init(registry.Timeout(3*time.Second), registry.Addrs("127.0.0.1:2379")); err != nil {
		t.Fatal(err)
	}

	go func() {
		time.Sleep(time.Second * 3)
		ser := &registry.Service{
			Name:     "test",
			Version:  "v1.0",
			Metadata: map[string]string{"a": "b"},
			Nodes:    []*registry.Node{&registry.Node{Id: "1", Address: "127.0.0.1:8080", Metadata: map[string]string{"c": "d"}}},
		}

		if err := register.Register(ser, registry.RegisterTTL(time.Second*10)); err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second * 2)

		ser = &registry.Service{
			Name:     "test",
			Version:  "v1.0",
			Metadata: map[string]string{"a": "b"},
			Nodes:    []*registry.Node{&registry.Node{Id: "1", Address: "127.0.0.1:8080", Metadata: map[string]string{"c": "d"}}},
		}

		if err := register.Register(ser); err != nil {
			t.Fatal(err)
		}

		println("register success")
	}()

	// ss, err := register.GetService("test")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//
	// for _, s := range ss {
	// 	t.Logf(s.Name)
	// 	t.Logf(s.Nodes[0].Id)
	// }

	watcher, err := register.Watch(registry.WatchService("test"))
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()

	for {
		println("begin watch")
		res, err := watcher.Next()
		if err != nil {
			t.Fatal(err)
		}
		println(time.Since(now).String())
		println(res.Action)
		println(res.Service.Name)
		println(len(res.Service.Nodes))
		println(res.Service.Nodes[0].Id)
	}
}
