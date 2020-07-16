package lb

import (
	"strconv"
)

// WeightRobin weight robin loadBalance impl
type WeightRobin struct {
	opts map[string]*weightRobin
}

// weightRobin used to save the weight info of server
type weightRobin struct {
	node *Node
	effectiveWeight int64
	currentWeight   int64
}

// NewWeightRobin create a WeightRobin
func NewWeightRobin() LoadBalance {
	return &WeightRobin{
		opts: make(map[string]*weightRobin, 1024),
	}
}

// Select select a server from servers using WeightRobin
func (w *WeightRobin) Select(servers []*Node, opts ...string) *Node {
	var total int64
	l := len(servers)
	if 0 >= l {
		return nil
	}

	best := ""
	for i := l - 1; i >= 0; i-- {
		svr := servers[i]
		weight, err := strconv.ParseInt(svr.Metadata["weight"], 10, 64)
		if err != nil {
			return nil
		}
		id := svr.Id

		if _, ok := w.opts[id]; !ok {
			w.opts[id] = &weightRobin{
				node: svr,
				effectiveWeight: weight,
			}
		}

		wt := w.opts[id]
		wt.currentWeight += wt.effectiveWeight
		total += wt.effectiveWeight

		if wt.effectiveWeight < weight {
			wt.effectiveWeight++
		}

		if best == "" || w.opts[best] == nil || wt.currentWeight > w.opts[best].currentWeight {
			best = id
		}
	}

	if best == "" {
		return nil
	}

	w.opts[best].currentWeight -= total

	return w.opts[best].node
}
