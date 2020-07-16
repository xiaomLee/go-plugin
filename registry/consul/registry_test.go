package consul

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	consul "github.com/hashicorp/consul/api"
	"github.com/xiaomLee/go-plugin/registry"
)

type mockRegistry struct {
	body   []byte
	status int
	err    error
	url    string
}

func encodeData(obj interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func newMockServer(rg *mockRegistry, l net.Listener) error {
	mux := http.NewServeMux()
	mux.HandleFunc(rg.url, func(w http.ResponseWriter, r *http.Request) {
		if rg.err != nil {
			http.Error(w, rg.err.Error(), 500)
			return
		}
		w.WriteHeader(rg.status)
		w.Write(rg.body)
	})
	return http.Serve(l, mux)
}

func newConsulTestRegistry(r *mockRegistry) (*consulRegistry, func()) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		// blurgh?!!
		panic(err.Error())
	}
	cfg := consul.DefaultConfig()
	println(l.Addr().String())
	cfg.Address = l.Addr().String()

	go newMockServer(r, l)

	var cr = &consulRegistry{
		config:      cfg,
		Address:     []string{cfg.Address},
		options:     registry.Options{},
		register:    make(map[string]uint64),
		lastChecked: make(map[string]time.Time),
		queryOptions: &consul.QueryOptions{
			AllowStale: true,
		},
	}
	cr.Client()

	return cr, func() {
		l.Close()
	}
}

func newServiceList(svc []*consul.ServiceEntry) []byte {
	bts, _ := encodeData(svc)
	return bts
}

func TestConsul_GetService_WithError(t *testing.T) {
	cr, cl := newConsulTestRegistry(&mockRegistry{
		err: errors.New("client-error"),
		url: "/v1/health/service/service-name",
	})
	defer cl()

	if _, err := cr.GetService("test-service"); err == nil {
		t.Fatalf("Expected error not to be `nil`")
	}
}

func TestConsul_GetService_WithHealthyServiceNodes(t *testing.T) {
	// warning is still seen as healthy, critical is not
	svcs := []*consul.ServiceEntry{
		newServiceEntry(
			"node-name-1", "node-address-1", "service-name", "v1.0.0",
			[]*consul.HealthCheck{
				newHealthCheck("node-name-1", "service-name", "passing"),
				newHealthCheck("node-name-1", "service-name", "warning"),
			},
		),
		newServiceEntry(
			"node-name-2", "node-address-2", "service-name", "v1.0.0",
			[]*consul.HealthCheck{
				newHealthCheck("node-name-2", "service-name", "passing"),
				newHealthCheck("node-name-2", "service-name", "warning"),
			},
		),
	}

	cr, cl := newConsulTestRegistry(&mockRegistry{
		status: 200,
		body:   newServiceList(svcs),
		url:    "/v1/health/service/service-name",
	})
	defer cl()

	svc, err := cr.GetService("service-name")

	for _, s := range svc {
		fmt.Printf("service:%+v \n", s)
	}

	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	if exp, act := 1, len(svc); exp != act {
		t.Fatalf("Expected len of svc to be `%d`, got `%d`.", exp, act)
	}

	if exp, act := 2, len(svc[0].Nodes); exp != act {
		t.Fatalf("Expected len of nodes to be `%d`, got `%d`.", exp, act)
	}
}

func TestConsul_GetService_WithUnhealthyServiceNode(t *testing.T) {
	// warning is still seen as healthy, critical is not
	svcs := []*consul.ServiceEntry{
		newServiceEntry(
			"node-name-1", "node-address-1", "service-name", "v1.0.0",
			[]*consul.HealthCheck{
				newHealthCheck("node-name-1", "service-name", "passing"),
				newHealthCheck("node-name-1", "service-name", "warning"),
			},
		),
		newServiceEntry(
			"node-name-2", "node-address-2", "service-name", "v1.0.0",
			[]*consul.HealthCheck{
				newHealthCheck("node-name-2", "service-name", "passing"),
				newHealthCheck("node-name-2", "service-name", "critical"),
			},
		),
	}

	cr, cl := newConsulTestRegistry(&mockRegistry{
		status: 200,
		body:   newServiceList(svcs),
		url:    "/v1/health/service/service-name",
	})
	defer cl()

	svc, err := cr.GetService("service-name")
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	if exp, act := 1, len(svc); exp != act {
		t.Fatalf("Expected len of svc to be `%d`, got `%d`.", exp, act)
	}

	if exp, act := 1, len(svc[0].Nodes); exp != act {
		t.Fatalf("Expected len of nodes to be `%d`, got `%d`.", exp, act)
	}
}

func TestConsul_GetService_WithUnhealthyServiceNodes(t *testing.T) {
	// warning is still seen as healthy, critical is not
	svcs := []*consul.ServiceEntry{
		newServiceEntry(
			"node-name-1", "node-address-1", "service-name", "v1.0.0",
			[]*consul.HealthCheck{
				newHealthCheck("node-name-1", "service-name", "passing"),
				newHealthCheck("node-name-1", "service-name", "critical"),
			},
		),
		newServiceEntry(
			"node-name-2", "node-address-2", "service-name", "v1.0.0",
			[]*consul.HealthCheck{
				newHealthCheck("node-name-2", "service-name", "passing"),
				newHealthCheck("node-name-2", "service-name", "critical"),
			},
		),
	}

	cr, cl := newConsulTestRegistry(&mockRegistry{
		status: 200,
		body:   newServiceList(svcs),
		url:    "/v1/health/service/service-name",
	})
	defer cl()

	svc, err := cr.GetService("service-name")
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	if exp, act := 1, len(svc); exp != act {
		t.Fatalf("Expected len of svc to be `%d`, got `%d`.", exp, act)
	}

	if exp, act := 0, len(svc[0].Nodes); exp != act {
		t.Fatalf("Expected len of nodes to be `%d`, got `%d`.", exp, act)
	}
}

func TestRegister(t *testing.T) {
	var register registry.Registry
	register = NewRegistry()

	if err := register.Init(registry.Addrs("localhost:8500"), TCPCheck(5*time.Second)); err != nil {
		t.Fatal(err)
	}

	doreg := func() {
		ser := &registry.Service{
			Name:     "test",
			Version:  "v1.0",
			Metadata: map[string]string{"a": "b"},
			Nodes:    []*registry.Node{&registry.Node{Id: "1", Address: "127.0.0.1:3306", Metadata: map[string]string{"c": "d"}}},
		}

		if err := register.Register(ser); err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second * 2)

		ser = &registry.Service{
			Name:     "test",
			Version:  "v1.0",
			Metadata: map[string]string{"a": "b"},
			Nodes:    []*registry.Node{&registry.Node{Id: "2", Address: "mysql:3306", Metadata: map[string]string{"c": "d"}}},
		}

		if err := register.Register(ser); err != nil {
			t.Fatal(err)
		}

		println("register success")
	}

	//doreg()
	go doreg()

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

	for {
		println("begin watch")
		now := time.Now()
		res, err := watcher.Next()
		if err != nil {
			t.Fatal(err)
		}
		println(time.Since(now).String())
		println(res.Action)
		fmt.Printf("service:%+v \n", res.Service)
	}
}
