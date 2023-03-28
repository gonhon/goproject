package cache

import "container/list"

//规定 Front(尾)为最新的元素  Back(头)为最久未使用
type Cache struct {
	//允许使用的最大内存
	maxBytes int64
	//当前已使用的内存
	nbytes int64
	//双线链表
	ll *list.List
	//value存双向链表的元素
	cache map[string]*list.Element
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

func (c *Cache) Add(key string, val Value) {
	//先检查有没有
	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element)
		kv := element.Value.(*entry)
		//修改容量 计算新的val容量
		c.nbytes += int64(val.Len()) - int64(kv.value.Len())
		kv.value = val
	} else {
		element = c.ll.PushFront(&entry{key: key, value: val})
		c.cache[key] = element
		//加入key val的容量
		c.nbytes += int64(val.Len()) + int64(len(key))
	}

	//检查是否容量是否达到maxBytes
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOld()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
