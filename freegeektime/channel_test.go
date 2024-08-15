package freegeektime

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestProduceConsume(t *testing.T) {
	ch := make(chan int, 5)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		produce(ch)
		wg.Done()
	}()

	go func() {
		consume(ch)
		wg.Done()
	}()
	wg.Wait()
}

func TestSpawn(t *testing.T) {
	println("start a working...")
	c := spawn(worker)
	<-c
	fmt.Println("worker work done!")
}

func TestSpawnGroup(t *testing.T) {
	groupSignal := make(chan signal)
	c := spawnGroup(workerI, 5, groupSignal)
	time.Sleep(5 * time.Second)
	fmt.Println("the group of workers start to work...")
	close(groupSignal)
	<-c
	fmt.Println("the group of workers work done!")

}

func TestIncr(t *testing.T) {
	cter := NewCounter()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			v := cter.Incr()
			fmt.Printf("groutine-%d:current counter value is %d\n", i, v)
			wg.Done()
		}(i)

	}
	wg.Wait()
}
