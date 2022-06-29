package defers

import (
	"sync"
	"testing"
)

func TestArea(t *testing.T) {
	wait := sync.WaitGroup{}
	a := &rect{-10, 20}
	wait.Add(1)
	go a.area(&wait)

	b := &rect{10, -20}
	wait.Add(1)
	go b.area(&wait)

	c := &rect{10, 20}
	wait.Add(1)
	go c.area(&wait)

	wait.Wait()
}
