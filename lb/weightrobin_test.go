package lb

import (
	"testing"
)

func TestWeightRobin_Select(t *testing.T) {

	lb := NewWeightRobin()
	retCount := make(map[string]int)
	for i := 0; i < 100 ; i++ {
		node := lb.Select(Servers)
		PrintNode(t, node)
		if _, ok := retCount[node.Id]; !ok {
			retCount[node.Id] = 1
			continue
		}
		retCount[node.Id] += 1
	}

	t.Logf("retCount: %v\n", retCount)

}
