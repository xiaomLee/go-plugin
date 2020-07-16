package lb

import "testing"

var (
	Servers = []*Node{
		&Node{
			Id:       "1",
			Address:  ":9001",
			Metadata: map[string]string{
				"weight": "10",
			},
		},
		&Node{
			Id:       "2",
			Address:  ":9002",
			Metadata: map[string]string{
				"weight": "20",
			},
		},
		&Node{
			Id:       "3",
			Address:  ":9003",
			Metadata: map[string]string{
				"weight": "20",
			},
		},
		&Node{
			Id:       "4",
			Address:  ":9004",
			Metadata: map[string]string{
				"weight": "20",
			},
		}, &Node{
			Id:       "5",
			Address:  ":9005",
			Metadata: map[string]string{
				"weight": "30",
			},
		},
	}
)

func PrintNode(t *testing.T, node *Node) {
	if node == nil {
		t.Error("Test_HashBalance select return nil")
	}
	t.Logf("node=%v", node)
}