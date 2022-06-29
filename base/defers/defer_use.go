package defers

import (
	"fmt"
	"sync"
)

type rect struct {
	length int
	width  int
}

func (r *rect) area(wg *sync.WaitGroup) int {
	defer wg.Done()
	if r.length <= 0 {
		fmt.Printf("rect %v's length should be greater than zero\n", r)
		return 0
	}
	if r.width <= 0 {
		fmt.Printf("rect %v's width should be greater than zero\n", r)
		return 0
	}
	areas := r.length * r.width
	fmt.Printf("rect-->length:%d width:%d areas:%d\n", r.length, r.width, areas)
	return areas
}
