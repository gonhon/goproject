package towSum

import (
	"fmt"
	"testing"
)

func TestTowSum(t *testing.T) {
	arrays := []int{1, 5, 6, 8}
	fmt.Println(arrays)
	arrays = towSum(arrays, 11)
	fmt.Println(arrays)
}
