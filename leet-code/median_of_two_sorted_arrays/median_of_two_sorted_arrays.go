package median_of_two_sorted_arrays

import (
	leet_code "github.com/limerence-code/goproject/leet-code"
)

func findMedianSortedArrays(num1 []int, num2 []int) float64 {
	// 假设 num1 的长度小
	if len(num1) > len(num2) {
		num1, num2 = num2, num1
	}

	low, high, k, n1len, n2len := 0, len(num1), (len(num1)+len(num2)+1)>>1, 0, 0
	for low <= high {
		// num1:  ……………… num1[n1len-1] | num1[n1len] ……………………
		// num2:  ……………… num2[n2len-1] | num2[n2len] ……………………
		n1len = low + (high-low)>>1 // 分界限右侧是 mid，分界线左侧是 mid - 1
		n2len = k - n1len
		if n1len > 0 && num1[n1len-1] > num2[n2len] { // num1 中的分界线划多了，要向左边移动
			high = n1len - 1
		} else if n1len != len(num1) && num1[n1len] < num2[n2len-1] { // num1 中的分界线划少了，要向右边移动
			low = n1len + 1
		} else {
			// 找到合适的划分了，需要输出最终结果了
			// 分为奇数偶数 2 种情况
			break
		}
	}
	midLeft, midRight := 0, 0
	if n1len == 0 {
		midLeft = num2[n2len-1]
	} else if n2len == 0 {
		midLeft = num1[n1len-1]
	} else {
		midLeft = leet_code.Max(num1[n1len-1], num2[n2len-1])
	}
	if (len(num1)+len(num2))&1 == 1 {
		return float64(midLeft)
	}
	if n1len == len(num1) {
		midRight = num2[n2len]
	} else if n2len == len(num2) {
		midRight = num1[n1len]
	} else {
		midRight = leet_code.Min(num1[n1len], num2[n2len])
	}
	return float64(midLeft+midRight) / 2
}
func findMedianSortedArraysBak(num1 []int, num2 []int) float64 {
	if len(num1) > len(num2) {
		num1, num2 = num2, num1
	}
	n1len, n2len := len(num1), len(num2)

	left, right, totalLeft := 0, n1len, (n1len+n2len)>>1

	for left < right {
		l := left + (right-left+1)>>1
		r := totalLeft - l
		if num1[l-1] > num2[r] {
			right = l - 1
		} else {
			left = l
		}
	}
	right = totalLeft - left
	num1LeftMax, num1RightMin, num2LeftMax, num2RightMin := 0, 0, 0, 0

	if left != 0 {
		num1LeftMax = num1[left-1]
	}
	if left != n1len {
		num1RightMin = num1[left]
	}
	if right != 0 {
		num2LeftMax = num2[right-1]
	}
	if right != n1len {
		num2RightMin = num2[right]
	}

	if (n1len+n2len)%2 == 0 {
		return float64(leet_code.Max(num1LeftMax, num2LeftMax))
	} else {
		return float64((leet_code.Max(num1LeftMax, num2LeftMax) + leet_code.Min(num1RightMin, num2RightMin)) >> 1)
	}
}

func findMedianSortedArraysTest(num1 []int, num2 []int) float64 {
	// 假设 num1 的长度小
	if len(num1) > len(num2) {
		return findMedianSortedArraysTest(num2, num1)
	}
	low, high, k, n1len, n2len := 0, len(num1), (len(num1)+len(num2)+1)>>1, 0, 0
	for low <= high {
		// num1:  ……………… num1[n1len-1] | num1[n1len] ……………………
		// num2:  ……………… num2[n2len-1] | num2[n2len] ……………………
		n1len = low + (high-low)>>1 // 分界限右侧是 mid，分界线左侧是 mid - 1
		n2len = k - n1len
		if n1len > 0 && num1[n1len-1] > num2[n2len] { // num1 中的分界线划多了，要向左边移动
			high = n1len - 1
		} else if n1len != len(num1) && num1[n1len] < num2[n2len-1] { // num1 中的分界线划少了，要向右边移动
			low = n1len + 1
		} else {
			// 找到合适的划分了，需要输出最终结果了
			// 分为奇数偶数 2 种情况
			break
		}
	}
	midLeft, midRight := 0, 0
	if n1len == 0 {
		midLeft = num2[n2len-1]
	} else if n2len == 0 {
		midLeft = num1[n1len-1]
	} else {
		midLeft = max(num1[n1len-1], num2[n2len-1])
	}
	if (len(num1)+len(num2))&1 == 1 {
		return float64(midLeft)
	}
	if n1len == len(num1) {
		midRight = num2[n2len]
	} else if n2len == len(num2) {
		midRight = num1[n1len]
	} else {
		midRight = min(num1[n1len], num2[n2len])
	}
	return float64(midLeft+midRight) / 2
}

func finArrays(nums1 []int, nums2 []int) float64 {
	if len(nums1) > len(nums2) {
		nums1, nums2 = nums2, nums1
	}
	low, high, count, index1, index2 := 0, len(nums1), (len(nums1)+len(nums2)+1)>>1, 0, 0
	for low <= high {
		//二分查找
		index1 = low + (high-low)>>1
		index2 = count - index1
		//左边的数大了 num1需要向 左移动
		if index1 > 0 && nums1[index1-1] > nums2[index2] {
			high = index1 - 1
		} else if index1 != len(nums1) && nums1[index1] < nums2[index2-1] { //右边偏大 num2需要向 右边移
			low = index1 + 1
		} else {
			break
		}
	}
	left, rigth := 0, 0

	//获取左边最大值
	if index1 == 0 { //num1 到左边了
		left = nums2[index2-1]
	} else if index2 == 0 { //num2到右边了
		left = nums1[index1-1]
	} else {
		//左边最大值
		left = max(nums1[index1-1], nums2[index2-1])
	}

	//最后一位是0 表示奇数 0 表示偶数
	if (len(nums1)+len(nums2))&1 != 0 {
		//左边本来就多一位
		return float64(left)
	}

	//获取右边最大值
	if index1 == len(nums1) { //num2到右边了
		rigth = nums2[index2]
	} else if index2 == len(nums2) {
		rigth = nums1[index1]
	} else {
		//右边最小值
		rigth = min(nums1[index1], nums2[index2])
	}
	//return float64((left + rigth) >> 1)
	return float64(left+rigth) / 2
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}
