package add_two

//https://leetcode.cn/problems/add-two-numbers/

type ListNodeBak struct {
	Val  int
	Next *ListNodeBak
}

func addTwoNumbersBak(l1 *ListNodeBak, l2 *ListNodeBak) *ListNodeBak {
	head := &ListNodeBak{Val: 0}
	n1, n2, nextVal, current := 0, 0, 0, head

	for l1 != nil || l2 != nil || nextVal != 0 {
		if l1 == nil {
			n1 = 0
		} else {
			n1 = l1.Val
			l1 = l1.Next
		}

		if l2 == nil {
			n2 = 0
		} else {
			n2 = l2.Val
			l2 = l2.Next
		}
		//当前节点%10
		current.Next = &ListNodeBak{
			Val: (n1 + n2 + nextVal) % 10,
		}
		//当前结点变为尾节点
		current = current.Next
		//满10进1
		nextVal = (n1 + n2 + nextVal) / 10
	}
	return head.Next
}
