package common

//给定一个排序数组，你需要在 原地 删除重复出现的元素，使得每个元素只出现一次，返回移除后数组的新长度。
//不要使用额外的数组空间，你必须在 原地 修改输入数组 并在使用 O(1) 额外空间的条件下完成。
// nums = [1,1,2] ----> [1,2]
func removeDuplicates(nums []int) []int {
	i := 0
	for j := 1; j < len(nums); j++ {
		if nums[j] != nums[i] {
			i++
			nums[i] = nums[j]
		}
	}
	return nums[:i+1]
}

//给定一个数组，它的第 i 个元素是一支给定股票第 i 天的价格。
//
//设计一个算法来计算你所能获取的最大利润。你可以尽可能地完成更多的交易（多次买卖一支股票）。
//
//注意：你不能同时参与多笔交易（你必须在再次购买前出售掉之前的股票）。
//输入: [7,1,5,3,6,4]
//输出: 7

// 累计盈利
func maxProfitV1(prices []int) int {
	maxProfit := 0
	for i := 0; i < len(prices)-1; i++ {
		if prices[i] < prices[i+1] {
			maxProfit += prices[i+1] - prices[i]
		}
	}
	return maxProfit
}

// 波峰波谷
func maxProfitV2(prices []int) int {
	var i, valley, peak int
	profit := 0
	for i < len(prices)-1 {
		for i < len(prices)-1 {
			if prices[i] <= prices[i+1] {
				break
			}
			i++
		}
		valley = prices[i]
		for i < len(prices) {
			if prices[i] >= prices[i+1] {
				break
			}
			i++
		}
		peak = prices[i]
		profit += peak - valley
	}

	return profit
}

// 旋转数组
// 给定一个数组，将数组中的元素向右移动 k 个位置，其中 k 是非负数。
// 输入: [1,2,3,4,5,6,7] 和 k = 3
//输出: [5,6,7,1,2,3,4]
//解释:
//向右旋转 1 步: [7,1,2,3,4,5,6]
//向右旋转 2 步: [6,7,1,2,3,4,5]
//向右旋转 3 步: [5,6,7,1,2,3,4]
//说明:
//尽可能想出更多的解决方案，至少有三种不同的方法可以解决这个问题。
//要求使用空间复杂度为 O(1) 的 原地 算法。

// 暴力循环
func rotateV1(nums []int, k int) {
	for i := 0; i < k; i++ {
		pre := nums[len(nums)-1]
		for j := 0; j < len(nums); j++ {
			temp := nums[j]
			nums[j] = pre
			pre = temp
		}
	}
}

func reverse(nums []int, start, end int) {
	i := start
	j := end
	for i < j {
		temp := nums[i]
		nums[i] = nums[j]
		nums[j] = temp
		i++
		j--
	}
}

// 翻转法
// 输入: [1,2,3,4,5,6,7] 和 k = 3
// 第一步：[7,6,5,4,3,2,1]
// 第二步：[5,6,7,4,3,2,1]
// 第三步：[5,6,7,1,2,3,4]
func rotateV2(nums []int, k int) {
	reverse(nums, 0, len(nums)-1)
	reverse(nums, 0, k%len(nums)-1)
	reverse(nums, k%len(nums), len(nums)-1)
}
