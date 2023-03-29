package cache

import (
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
