package cache

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFunc(t *testing.T) {
	var f GetterFunc = func(s string) ([]byte, error) {
		return []byte(s), nil
	}

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Log("callback failed")
	} else {
		t.Logf("val:%s", string(v))
	}
}

func TestGroupRun(t *testing.T) {

	var base = map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	//存储个数
	loadCounts := make(map[string]int, len(base))
	group := NewGroup("key", 2<<5, GetterFunc(func(key string) ([]byte, error) {
		//从base中取出 模拟数据库
		if val, ok := base[key]; ok {
			t.Logf("init %s data ...", key)
			if _, ok := loadCounts[key]; !ok {
				loadCounts[key] = 0
			}
			loadCounts[key] += 1
			return []byte(val), nil
		}
		t.Logf("%s not  data ...", key)
		return nil, fmt.Errorf("%s not found", key)
	}))

	for k, v := range base {
		if view, err := group.Get(k); err != nil || view.String() != v {
			t.Logf("%s val not match", k)
		}

		if _, err := group.Get(k); err != nil || loadCounts[k] > 1 {
			t.Logf("   %s not found", k)
		}
	}

	if v, err := group.Get("key1"); err == nil {
		t.Log(v.String())
	}
	group.Get("key6")

}
