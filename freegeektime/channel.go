package freegeektime

import (
	"fmt"
	"sync"
	"time"
)

func produce(ch chan<- int) {
	for i := 0; i < 10; i++ {
		ch <- i + 1
		time.Sleep(time.Second)
	}
	close(ch)
}

func consume(ch <-chan int) {
	for n := range ch {
		println(n)
	}
}

type signal struct{}

func worker() {
	println("worker is working...")
	time.Sleep(1 * time.Second)
}
func workerI(i int) {
	fmt.Printf("worker %d: is working...\n", i)
	time.Sleep(1 * time.Second)
	fmt.Printf("worker %d: is works done\n", i)
}

func spawn(f func()) <-chan signal {
	c := make(chan signal)
	go func() {
		println("worker start to work...")
		f()
		c <- signal{}
	}()
	return c
}

func spawnGroup(f func(int), num int, groupSignal <-chan signal) <-chan signal {
	c := make(chan signal)
	var wg sync.WaitGroup

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			<-groupSignal
			fmt.Printf("worker %d: start to work...\n", i)
			f(i)
			wg.Done()
		}(i + 1)
	}

	go func() {
		wg.Wait()
		c <- signal{}
	}()
	return c
}

type counter struct {
	c chan int
	i int
}

func NewCounter() *counter {
	cter := &counter{
		c: make(chan int),
	}
	go func() {
		for {
			cter.i++
			cter.c <- cter.i
		}
	}()
	return cter
}

func (cter *counter) Incr() int {
	return <-cter.c
}
