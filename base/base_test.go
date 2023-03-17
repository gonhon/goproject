package base

import (
	"testing"
	"time"
)

func TestConvert(t *testing.T) {
	convert()
}

func TestChan(t *testing.T) {
	chanFunc()
}

func TestChanChange(t *testing.T) {
	chanChange()
}

func TestSelectFunc(t *testing.T) {
	selectFunc()
}

func TestSelect(t *testing.T) {
	select {
	case ret := <-asycService():
		t.Log(ret)
	case <-time.After(time.Microsecond * 100):
		t.Log("timeout...")
	}
}

func TestOpertion(t *testing.T) {
	val, _ := exec(1, 2, func(x, y int) (int, error) {
		return (x + y), nil
	})
	t.Log("val:", val)
}
func TestCalculate(t *testing.T) {
	genFunc := genCalculate(func(x, y int) (int, error) {
		return (x + y), nil
	})
	val, _ := genFunc(5, 8)

	t.Log("val:", val)
}

func TestListFunc(t *testing.T) {
	listFunc()
}
func TestSliceFunc(t *testing.T) {
	sliceFunc()
}
