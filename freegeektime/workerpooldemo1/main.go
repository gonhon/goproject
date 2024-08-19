package main

import (
	"fmt"
	"time"

	"github.com/gonhon/workerpoolplus"
)

func main() {
	// p := workerpoolplus.New(5)
	p := workerpoolplus.New(5, workerpoolplus.WithBlock(false), workerpoolplus.WithPreAllocWorkers(false))

	for i := 0; i < 10; i++ {
		err := p.Schedule(func() {
			time.Sleep(3 * time.Second)
		})
		if err != nil {
			fmt.Println("task:", i, "err:", err)
		}
	}
	p.Free()
}
