package median_of_two_sorted_arrays

import "testing"

func Test01(t *testing.T) {
	a, b := []int{1, 3, 4, 9}, []int{2, 8}
	temp := findMedianSortedArraysTest(a, b)
	t.Log(temp)
	temp = findMedianSortedArrays(a, b)
	t.Log(temp)
	temp = finArrays(a, b)
	t.Log(temp)
	temp = findMedianSortedArraysBak(a, b)
	t.Log(temp)
	//2 8
	//1 3 4 9
	//1 2 3 4 8 9
}
