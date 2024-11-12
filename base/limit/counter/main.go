package main

import (
	"fmt"
	"sync"
	"time"
)

type Limiter struct {
	mu        sync.Mutex    // 保证并发安全的锁
	rate      int           // 限制请求数
	interval  time.Duration // 时间窗口
	timestamp time.Time     // 上一次重置的时间戳
	requests  int           // 当前窗口内的请求数
}

// NewLimiter 创建一个限流器实例
func NewLimiter(rate int, interval time.Duration) *Limiter {
	return &Limiter{
		rate:      rate,
		interval:  interval,
		timestamp: time.Now(),
		requests:  0,
	}
}

// Allow 方法判断是否允许通过请求
func (l *Limiter) Allow() bool {
	l.mu.Lock()         // 加锁，确保并发安全
	defer l.mu.Unlock() // 解锁

	// 如果当前时间超过了时间窗口，则重置计数器
	if time.Since(l.timestamp) > l.interval {
		l.requests = 0
		l.timestamp = time.Now()
	}

	// 检查请求数是否超过阈值
	if l.requests < l.rate {
		l.requests++ // 增加请求计数
		return true  // 允许通过
	}

	return false // 拒绝请求
}

func main() {
	// 创建一个限流器实例，每分钟最多允许1000次请求
	limiter := NewLimiter(1000, time.Minute)

	// 模拟请求
	for i := 0; i < 1100; i++ { // 超过1000次请求，测试限流效果
		if limiter.Allow() {
			fmt.Println("请求通过", i+1)
		} else {
			fmt.Println("请求被拒绝", i+1)
		}
	}
}
