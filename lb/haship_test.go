package lb

import (
	"testing"
)

func Test_HashBalance(t *testing.T) {
	lb := NewHashBalance()

	var lastId string
	for i := 0; i < 10; i++ {
		node := lb.Select(Servers, "192.168.0.5")
		PrintNode(t, node)

		if lastId != "" && lastId != node.Id {
			t.Error("Test_HashBalance return node different from last")
		}
	}
}

func Benchmark_HashBalance(b *testing.B) {
	lb := NewHashBalance()
	for i := 0; i < b.N; i++ {
		node := lb.Select(Servers, "192.168.0.5")
		if node == nil {
			b.Error("Test_HashBalance select return nil")
		}
	}
}
