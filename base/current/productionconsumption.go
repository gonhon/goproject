package current

import (
	"fmt"
	"math/rand"
	"sync"
)

const (
	MAX int = 20
	NIL int = 0
)

var (
	lock           sync.Mutex
	productionCond *sync.Cond = sync.NewCond(&lock)
	consumeCond    *sync.Cond = sync.NewCond(&lock)
	list           []int
)

func productionLocal() {
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

func consumeLock() {
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
