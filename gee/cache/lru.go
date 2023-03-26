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
