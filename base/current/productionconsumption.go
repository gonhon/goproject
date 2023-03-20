package current

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	MAX int = 20
	NIL int = 0
)

var (
	lock           sync.Mutex
	productionCond *sync.Cond = sync.NewCond(&lock)
	consumeCond    *sync.Cond = sync.NewCond(&lock)
	group          sync.WaitGroup
	list           []int
	cache          = make(chan int, MAX)
)

type prodcons interface {
	//生产
	production()
	//消费
	consume()
}

// 使用锁
type ProdconsLock struct {
}

// 使用chan
type ProdconsChan struct {
}

func (ProdconsLock) production() {
	lock.Lock()
	defer lock.Unlock()

	//满了
	for len(list) == MAX {
		productionCond.Wait()
	}
	val := rand.Intn(100)
	fmt.Printf("生产者... %d\n", val)
	list = append(list, val)
	//唤醒消费
	consumeCond.Broadcast()

}

func (ProdconsLock) consume() {
	lock.Lock()
	defer lock.Unlock()

	//为空
	for len(list) == NIL {
		consumeCond.Wait()
	}
	var index int = 0
	//获取第一个元素
	val := list[index]
	//移除第一个元素
	list = append(list[:index], list[index+1:]...)
	fmt.Printf("消费者...获取的值:%d\n", val)
	//唤醒生产
	productionCond.Broadcast()
}

//使用chan

func (ProdconsChan) production() {
	val := rand.Intn(100)
	fmt.Printf("production... %d\n", val)
	cache <- val
}

func (ProdconsChan) consume() {
	val := <-cache
	fmt.Printf("consume... %d\n", val)
}

func Exec(parmas prodcons, productionSleep, consumeSleep time.Duration) {
	//生产者
	for i := 1; i <= MAX; i++ {
		go func() {
			group.Add(1)
			defer group.Done()
			time.Sleep(productionSleep)
			parmas.production()
		}()
	}

	//消费者
	for j := 1; j <= MAX; j++ {
		go func() {
			group.Add(1)
			defer group.Done()
			time.Sleep(consumeSleep)
			parmas.consume()
		}()
	}
	group.Wait()
}
