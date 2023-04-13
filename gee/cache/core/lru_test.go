package core

import (
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestLruGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("a", String("a"))
	if val, ok := lru.Get("a"); ok {
		t.Logf("key: %s, value: %s", "a", val.(String))
	}
	if val, ok := lru.Get("b"); !ok || val == nil {
		t.Logf("key: %s is not found ...", "a")
	}

}

func TestLruRemove(t *testing.T) {
	k1, k2, k3 := "k1", "k2", "k3"
	val1, val2, val3 := "val1", "val2", "val3"
	//不包含key3的
	cap := len(k1 + val1 + k2 + val2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(val1))
	lru.Add(k2, String(val2))
	//加入第三对淘汰key1
	lru.Add(k3, String(val3))

	if val, ok := lru.Get(k1); !ok || val == nil {
		t.Logf("key: %s is not found remove...", k1)
	}
}

func TestLruOnEvicted(t *testing.T) {

	k1, k2, k3 := "k1", "k2", "k3"
	val1, val2, val3 := "val1", "val2", "val3"
	//不包含key3的
	cap := len(k1 + val1 + k2 + val2)
	lru := New(int64(cap), func(key string, val Value) {
		t.Logf("key: %s, value: %s OnEvicted ", key, val.(String))
	})

	lru.Add(k1, String(val1))
	lru.Add(k2, String(val2))
	//加入第三对淘汰key1
	lru.Add(k3, String(val3))

}
