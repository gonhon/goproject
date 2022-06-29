package add_two

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	//a1 := []int{2, 4, 3}
	//a2 := []int{5, 6, 4}
	a1 := &ListNode{
		Val: 2,
		Next: &ListNode{
			Val: 4,
			Next: &ListNode{
				Val: 3,
			},
		},
	}
	a2 := &ListNode{
		Val: 5,
		Next: &ListNode{
			Val: 6,
			Next: &ListNode{
				Val: 4,
			},
		},
	}
	data := addTwoNumbers(a1, a2)
	for data != nil {
		fmt.Printf("%d ", data.Val)
		//t.Log(data.Val)
		data = data.Next
	}
	fmt.Println()

}
