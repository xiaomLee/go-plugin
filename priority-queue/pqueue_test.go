package priority_queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPQueue_Put(t *testing.T) {
	pq := &PQueue{
		items: make([]int, 0),
	}
	pq.Put(3)
	pq.Put(4)
	pq.Put(2)
	pq.Put(5)
	pq.Put(1)
	t.Logf("pt.items:%v", pq.items)

	assert.EqualValues(t, 5, pq.Pop(), "")
	t.Logf("pt.items:%v", pq.items)
	assert.EqualValues(t, 4, pq.Pop(), "")
	t.Logf("pt.items:%v", pq.items)

	pq.Put(6)
	t.Logf("pt.items:%v", pq.items)

	assert.EqualValues(t, 6, pq.Pop(), "")

	pq.Put(7)
	pq.Put(8)
	pq.Put(9)
	t.Logf("pt.items:%v", pq.items)

	assert.EqualValues(t, 9, pq.Pop(), "")
	t.Logf("pt.items:%v", pq.items)
}
