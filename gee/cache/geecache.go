package cache

import (
	"fmt"
	"log"
	"sync"

	pb "github.com/limerence-code/goproject/gee/cache/cachepb"
	"github.com/limerence-code/goproject/gee/cache/singleflight"
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

	//lock Group
	loader *singleflight.Group
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
		loader:    &singleflight.Group{},
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
	return group.Load(key)
}

func (group *Group) RegisterPeers(peers PeerPicker) {
	if group.peers != nil {
		panic("RegisterPeers called more than once")
	}
	group.peers = peers
}

func (group *Group) getFromData(fun PeerGetter, key string) (ByteView, error) {

	req := &pb.Request{
		Group: group.name,
		Key:   key,
	}
	res := &pb.Response{}

	err := fun.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{
		b: res.Value,
	}, err

}

func (group *Group) Load(key string) (ByteView, error) {
	val, err := group.loader.Do(key, func() (interface{}, error) {
		if group.peers != nil {
			if peer, ok := group.peers.PickPeer(key); ok {
				//将值赋给value
				if value, err := group.getFromData(peer, key); err == nil {
					return value, err
				} else {
					log.Println("[Cache] Failed to get from peer", err)
				}
			}
		}
		return group.getLocally(key)
	})
	return val.(ByteView), err
}
