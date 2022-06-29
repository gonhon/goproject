package chans

import (
	"fmt"
	"sync"
	"time"
)

func test01() {

	flag := make(chan bool)

	go func(f chan bool) {
		fmt.Println("begin......")
		time.Sleep(3000 * time.Millisecond)
		f <- true
	}(flag)

	<-flag

	fmt.Printf("end....")
}

func test01Add(c chan int, a int) {
	time.Sleep(time.Millisecond * 2000)
	c <- a
}

func test02() {

	chan1 := make(chan int)
	chan2 := make(chan int)

	go test01Add(chan1, 1)
	go test01Add(chan2, 2)

	select {
	case v := <-chan1:
		fmt.Println("chan1===>", v)
	case v := <-chan2:
		fmt.Println("chan2===>", v)
	}
}

func Test3() {
	defer func() {
		if error := recover(); error != nil {
			fmt.Println("err:", error)
		}
	}()
	var x = 0
	chanLock := make(chan bool, 1)
	// chanLock := make(chan bool)
	var watiG sync.WaitGroup
	for i := 0; i < 10000; i++ {
		watiG.Add(1)
		go func() {
			chanLock <- true
			x += 1
			<-chanLock
			watiG.Done()
		}()
	}
	watiG.Wait()
	fmt.Println("x=", x)
}

var x = 0

func increment(wg *sync.WaitGroup, ch chan bool) {
	ch <- true
	x = x + 1
	<-ch
	wg.Done()
}
func test04() {
	var w sync.WaitGroup
	ch := make(chan bool, 1)
	for i := 0; i < 1000; i++ {
		w.Add(1)
		go increment(&w, ch)
	}
	w.Wait()
	fmt.Println("final value of x", x)
}

func test05() {
	ch := make(chan int)
	go func(c chan int) {
		for i := 0; i < 10; i++ {
			c <- i
		}
		close(c)
	}(ch)

	for {
		v, ok := <-ch
		if !ok {
			break
		}
		fmt.Println("res:", v, ok)
	}

}
