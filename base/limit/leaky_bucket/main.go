package main

import (
	"fmt"
	"sync"
	"time"
)

type LeakyBucket struct {
	capacity     int
	rate         time.Duration
	current      int
	lastLeakTime time.Time
	mu           sync.Mutex
}

func NewLeakyBucket(capacity int, rate time.Duration) *LeakyBucket {
	return &LeakyBucket{
		capacity:     capacity,
		rate:         rate,
		current:      0,
		lastLeakTime: time.Now(),
	}
}

func (lb *LeakyBucket) leak() {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(lb.lastLeakTime)

	// 计算漏出的请求数量
	leaked := int(elapsed / lb.rate)
	if leaked > 0 {
		lb.current -= leaked
		if lb.current < 0 {
			lb.current = 0
		}
		lb.lastLeakTime = now
	}
}

func (lb *LeakyBucket) Allow() bool {
	lb.leak() // 先执行漏水操作

	lb.mu.Lock()
	defer lb.mu.Unlock()

	if lb.current < lb.capacity {
		lb.current++
		return true // 允许请求
	}
	return false // 拒绝请求
}

func main() {
	// 创建一个容量为5，漏水速率为每秒1个请求的漏斗
	lb := NewLeakyBucket(5, 1*time.Second)

	// 模拟请求
	for i := 0; i < 10; i++ {
		if lb.Allow() {
			fmt.Printf("请求 %d 被允许\n", i+1)
		} else {
			fmt.Printf("请求 %d 被拒绝\n", i+1)
		}
		time.Sleep(300 * time.Millisecond) // 模拟请求间隔
	}
}
