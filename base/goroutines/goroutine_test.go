package goroutines

import (
	"testing"
)

func TestMain(t *testing.T) {
	Testgoroutines1()
}

func TestMain1(t *testing.T) {
	Testgoroutines2(200)
	t.Logf("456")
}

func TestMain2(t *testing.T) {
	Testgoroutines3()
}

func TestMain4(t *testing.T) {
	Testgoroutines4()
}

// go test -v  -test.run TestMain1 指定方法
