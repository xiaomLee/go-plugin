package priority_queue

func HeapSort(arr []int) {
	// 1. build max heap: start from last parent node
	// 2. swap max with last
	// 3. repeat sink for arr[:length-2]

	n := len(arr)
	for i := len(arr)/2 - 1; i >= 0; i-- {
		sink(arr, i, len(arr))
	}

	for j := n - 1; j > 0; j-- {
		// swap arr[0] arr[j]
		arr[0], arr[j] = arr[j], arr[0]
		sink(arr, 0, j)
	}
}

// sink i
func sink(arr []int, i, n int) {
	for {
		root := i
		left := 2*i + 1
		right := 2*i + 2
		if left < n && arr[root] < arr[left] {
			root = left
		}
		if right < n && arr[root] < arr[right] {
			root = right
		}
		// i > left && i > right
		if i == root {
			break
		}
		// swap  arr[i] arr[root]
		arr[i], arr[root] = arr[root], arr[i]
		i = root
	}
}
