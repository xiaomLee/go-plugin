package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoveDuplicates(t *testing.T) {
	nums := []int{0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 3, 4, 4, 4, 4, 4, 4, 5, 6, 6, 6, 6}
	removeDuplicates(nums)
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6}, nums)
}

func TestReverse(t *testing.T) {
	nums := []int{1, 2, 3, 4}
	reverse(nums, 0, len(nums))
	assert.Equal(t, []int{4, 3, 2, 1}, nums)
	reverse(nums, 0, 3)
	assert.Equal(t, []int{2, 3, 4, 1}, nums)
	reverse(nums, 1, len(nums))
	assert.Equal(t, []int{2, 1, 4, 3}, nums)
}

func TestRotate(t *testing.T) {
	nums1 := []int{1, 2, 3, 4, 5, 6, 7}
	rotateV1(nums1, 3)
	assert.Equal(t, []int{5, 6, 7, 1, 2, 3, 4}, nums1)

	nums2 := []int{1, 2, 3, 4, 5, 6, 7}
	rotateV2(nums2, 3)
	assert.Equal(t, []int{5, 6, 7, 1, 2, 3, 4}, nums2)
}

func TestPermutation(t *testing.T) {
	s := "abcd"
	ret := Permutation(s)
	fmt.Println(len(ret), ret)
}
