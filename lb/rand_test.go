package lb

import (
	"testing"
)

func Test_RandBalance(t *testing.T) {
	lb := NewRandBalance()
	for i := 0; i < 66; i++ {
		node := lb.Select(Servers)
		PrintNode(t, node)
	}
}

func Benchmark_RandBalance(b *testing.B) {
	lb := NewRandBalance()
	for i := 0; i < b.N; i++ {
		node := lb.Select(Servers)
		if node == nil {
			b.Error("select return nil")
		}
	}
}
