package add_two

//https://leetcode.cn/problems/add-two-numbers/

type ListNode struct {
	Val  int
	Next *ListNode
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	head := &ListNode{
		Val: 0,
	}
	v1, v2, remainder, currentNode := 0, 0, 0, head

	for l1 != nil || l2 != nil || remainder != 0 {
		if l1 != nil {
			v1 = l1.Val
			l1 = l1.Next
		} else {
			v1 = 0
		}

		if l2 != nil {
			v2 = l2.Val
			l2 = l2.Next
		} else {
			v2 = 0
		}

		//计算出当前结点
		currentNode.Next = &ListNode{
			Val: (v1 + v2 + remainder) % 10,
		}
		//余数 下一步计算用到
		remainder = (v1 + v2 + remainder) / 10
		//当前结点指向新值
		currentNode = currentNode.Next
	}
	return head
}
