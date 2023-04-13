package base

import (
	"fmt"
	"net/http"
	"testing"
)

func TestCacheHttp(t *testing.T) {
	var base = map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	//存储个数
	loadCounts := make(map[string]int, len(base))
	NewGroup("key", 2<<5, GetterFunc(func(key string) ([]byte, error) {
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

	addr := "127.0.0.1:8900"
	httpPoll := NewHttpPoll(addr)
	t.Log("addr", addr)
	http.ListenAndServe(addr, httpPoll)

}
