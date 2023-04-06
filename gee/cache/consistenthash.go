package cache

import (
	"hash/crc32"
	"log"
	"sort"
	"strconv"
)

type Hash func([]byte) uint32

// 一致性hash
type Map struct {
	//hash函数
	hash Hash
	//每个key虚拟几点个数
	replicas int
	//所以的虚拟节点hash值
	keys []int
	//key 虚拟节点hash值 val 对应真实节点的值
	hashMap map[int]string
}

func NewMap(replicas int, hash Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     hash,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			//生成hash
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			//将所有的虚拟节点的hash放入keys中
			m.keys = append(m.keys, hash)
			//虚拟节点的key(hash)映射到正式i二点
			m.hashMap[hash] = key
		}
	}
	//排序
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	//先根据key获取对应的hash
	hash := int(m.hash([]byte(key)))
	//在hash环上查到对应索引
	index := sort.Search(len(m.keys), func(i int) bool {
		//key对应的hash>=当前可以的hash
		return m.keys[i] >= hash
	})
	// return m.hashMap[m.keys[index&(len(m.keys)-1)]]
	val, ok := m.hashMap[m.keys[index%len(m.keys)]]
	if !ok {
		log.Printf("key :%s not found", key)
	}
	return val

}
