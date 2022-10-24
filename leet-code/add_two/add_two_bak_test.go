package add_two

import (
	"fmt"
	"testing"
)

func TestAddTwoNumbersBak(t *testing.T) {
	//a1 := []int{2, 4, 3}
	//a2 := []int{5, 6, 4}
	//a2 := []int{7, 0, 8}

	a1 := &ListNodeBak{
		Val: 2,
		Next: &ListNodeBak{
			Val: 4,
			Next: &ListNodeBak{
				Val: 3,
			},
		},
	}
	a2 := &ListNodeBak{
		Val: 5,
		Next: &ListNodeBak{
			Val: 6,
			Next: &ListNodeBak{
				Val: 4,
			},
		},
	}
	data := addTwoNumbersBak(a1, a2)
	for data != nil {
		fmt.Printf("%d ", data.Val)
		//t.Log(data.Val)
		data = data.Next
	}
	fmt.Println()

}
