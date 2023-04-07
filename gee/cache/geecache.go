package cache

import (
	"fmt"
	"sync"
)

// 接口形函数
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker
}

var (
	//使用读写锁
	mutex  sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("getter is nil")
	}
	mutex.Lock()
	defer mutex.Unlock()

	group := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	mutex.RLock()
	defer mutex.RUnlock()
	return groups[name]
}

// 加入元素
func (group *Group) populateCache(key string, val ByteView) {
	group.mainCache.Add(key, val)
}

// 回调Get方法初始化缓存
func (group *Group) getLocally(key string) (ByteView, error) {
	bytes, err := group.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	val := ByteView{b: cloneBytes(bytes)}
	group.populateCache(key, val)
	return val, nil
}

func (group *Group) load(key string) (val ByteView, err error) {
	return group.getLocally(key)
}

// Get 根据key获取数据
func (group *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is empty")
	}
	//存在及从缓存取
	if v, ok := group.mainCache.Get(key); ok {
		return v, nil
	}
	//不存在调用Getter
	return group.load(key)
}

func (group *Group) RegisterPeers(peers PeerPicker) {
	if group.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	group.peers = peers
}
