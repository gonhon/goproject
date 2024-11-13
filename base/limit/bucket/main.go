package main

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	tokens     chan struct{}
	refillRate time.Duration
	bucketSize int
	mu         sync.Mutex
}

func NewTokenBucket(bucketSize int, refillRate time.Duration) *TokenBucket {
	tb := &TokenBucket{
		tokens:     make(chan struct{}, bucketSize),
		refillRate: refillRate,
		bucketSize: bucketSize,
	}

	// 启动一个协程来定期填充令牌
	go tb.refill()
	return tb
}

func (tb *TokenBucket) refill() {
	//创建一个定时器
	ticker := time.NewTicker(tb.refillRate)
	defer ticker.Stop()

	for range ticker.C {
		tb.mu.Lock()
		// 如果桶没有满，添加一个令牌
		if len(tb.tokens) < tb.bucketSize {
			tb.tokens <- struct{}{}
		}
		tb.mu.Unlock()
	}
}

func (tb *TokenBucket) GetToken() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	select {
	case <-tb.tokens:
		return true // 成功获取令牌
	default:
		return false // 无可用令牌
	}
}

func main() {
	// 创建一个大小为5，速率为1秒1个令牌的令牌桶
	tb := NewTokenBucket(5, 1*time.Second)

	// 模拟请求
	for i := 0; i < 10; i++ {
		if tb.GetToken() {
			fmt.Printf("请求 %d 被允许\n", i+1)
		} else {
			fmt.Printf("请求 %d 被拒绝\n", i+1)
		}
		time.Sleep(300 * time.Millisecond) // 模拟请求间隔
	}
}
