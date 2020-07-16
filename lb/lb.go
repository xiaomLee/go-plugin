package lb

var (
	supportLbs = []string{"roundRobin", "rand", "hash", "weightRobin"}
)

var (
	// LBS map loadBalance name and process function
	LBS = map[string]func() LoadBalance {
		"roundRobin":  NewRoundRobin,
		"weightRobin": NewWeightRobin,
		"hash":        NewHashBalance,
		"rand":        NewRandBalance,
	}
)

// LoadBalance loadBalance interface returns selected server's id
type LoadBalance interface {
	Select(servers []*Node, opts ...string) *Node
}

// GetSupportLBS return supported loadBalances
func GetSupportLBS() []string {
	return supportLbs
}

// NewLoadBalance create a LoadBalance,if LoadBalance function is not supported
// it will return NewRoundRobin
func NewLoadBalance(name string) LoadBalance {
	if l, ok := LBS[name]; ok {
		return l()
	}
	return NewRoundRobin()
}
