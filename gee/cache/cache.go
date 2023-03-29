package cache

import (
	"sync"
)

type cache struct {
	mutex      sync.Mutex
	lru        *LruCache
	cacheBytes int64
}

func (c *cache) Add(key string, val ByteView) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.lru == nil {
		//懒加载
		c.lru = New(c.cacheBytes, nil)
	}
	c.lru.Add(key, val)

}

func (c *cache) Get(key string) (val ByteView, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
