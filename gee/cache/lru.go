package cache

import "container/list"

type Cache struct {
	//允许使用的最大内存
	maxBytes int64
	//当前已使用的内存
	nbytes int64
	ll     *list.List
	cache  map[string]*list.Element
	//某条记录被移除时的回调函数
	OnEvicted func(string, Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}
func (c *Cache) Get(key string) (val Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		//将val移到尾
		c.ll.MoveToFront(element)
		v := element.Value.(*entry)
		return v.value, true
	}
	return nil, false
}

//移除淘汰结点
func (c *Cache) RemoveOld() {
	//从头部移除
	element := c.ll.Back()
	if element != nil {
		c.ll.Remove(element)
		kv := element.Value.(*entry)
		//删除map
		delete(c.cache, kv.key)
		//减去移除的大小
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}

}
