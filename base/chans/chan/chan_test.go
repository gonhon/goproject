package chans

import (
	"fmt"
	"testing"
	"time"
)

func TestChan1(t *testing.T) {
	test01()
}
func TestChan2(t *testing.T) {
	test02()
}
func TestChan3(t *testing.T) {
	Test3()
	a := make(chan bool)
	fmt.Println("len:", len(a))
	fmt.Println("cap:", cap(a))

}
func TestChan4(t *testing.T) {
	// test04()
	b := make(chan bool)

	go func() {
		time.Sleep(time.Millisecond * 3000)
		b <- true
	}()

	<-b
	fmt.Println("end...")

}

func TestChan5(t *testing.T) {
	test05()
}
