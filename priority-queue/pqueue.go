package priority_queue

type PriorityQueue interface {
	Put(item int)
	Pop() int
}

type PQueue struct {
	items []int
}

func (pq *PQueue) Put(item int) {
	// 1. append item to items
	// 2. swim items[length-1] and build max heap
	pq.items = append(pq.items, item)
	pq.swim(len(pq.items) - 1)
}

func pancakeSort(arr []int) []int {
	return sort(arr, len(arr))
}

func sort(arr []int, n int) []int {
	actions := make([]int, 0)
	if len(arr) <= 1 || n <= 1 {
		return actions
	}
	// find max
	maxIndex := 0
	max := arr[0]
	for i := 1; i < n; i++ {
		if arr[i] > max {
			max = arr[i]
			maxIndex = i
		}
	}
	// reverse max to first
	reverse(arr, 0, maxIndex)
	actions = append(actions, maxIndex+1)

	// reverse max to end
	reverse(arr, 0, n-1)
	actions = append(actions, n)

	actions = append(actions, sort(arr, n-1)...)

	return actions
}

func reverse(arr []int, i, j int) {
	for i < j {
		arr[i], arr[j] = arr[j], arr[i]
		i++
		j--
	}
}

func (pq *PQueue) swim(i int) {
	for {
		child := i
		root := i / 2
		if root >= 0 && pq.items[root] < pq.items[child] {
			// swap root child
			// make i == root
			pq.items[root], pq.items[child] = pq.items[child], pq.items[root]
			i = root
		}
		if i == child {
			break
		}
	}
}

func (pq *PQueue) Pop() int {
	// 1. swap items[0] items[length-1]
	// 2. pop items[length-1]
	// 3. sink item[0]
	n := len(pq.items)
	if n < 1 {
		return -1
	}

	pq.items[n-1], pq.items[0] = pq.items[0], pq.items[n-1]

	ans := pq.items[n-1]
	pq.items = pq.items[:n-1]

	pq.sink(0)
	return ans
}

func (pq *PQueue) sink(i int) {
	for {
		root := i
		left := 2*i + 1
		right := 2*i + 2
		if left < len(pq.items) && pq.items[root] < pq.items[left] {
			root = left
		}
		if right < len(pq.items) && pq.items[root] < pq.items[right] {
			root = right
		}
		if i == root {
			break
		}
		// root == max(left, right)
		// swap items[i] max(items[left], items[right]
		pq.items[i], pq.items[root] = pq.items[root], pq.items[i]
		// make i = max(left, right)
		i = root
	}
}
