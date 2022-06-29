package goroutines

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	maps = make(map[int]int, 10)
	lock sync.Mutex
)

//协程
func Testgoroutines1() {
	go test01()
	for i := 0; i < 10; i++ {
		fmt.Println("Testgoroutines1 ..." + strconv.Itoa(i))
		time.Sleep(time.Second)
	}
}
func test01() {
	for i := 0; i < 10; i++ {
		fmt.Println("test01 ..." + strconv.Itoa(i))
		time.Sleep(time.Second)
	}
}

//阶乘
func Testgoroutines2(count int) {
	for i := 1; i < count; i++ {
		go test02(i)
	}
	time.Sleep(time.Second * 10)

	lock.Lock()
	for k, v := range maps {
		fmt.Printf("%v-%v \n", k, v)
	}
	lock.Unlock()
}

func test02(count int) {
	res := 1
	for i := 1; i < count; i++ {
		res *= i
	}
	lock.Lock()
	maps[count] = res
	lock.Unlock()
}

//chan
func Testgoroutines3() {
	initChan := make(chan int, 50)
	exitChan := make(chan bool, 1)
	go writeData(initChan)
	go readData(initChan, exitChan)

	// time.Sleep(time.Second * 10)
	for {
		_, ok := <-exitChan
		if !ok {
			break
		}
	}
}

func writeData(initChan chan int) {
	for i := 0; i < 50; i++ {
		//放入数据
		initChan <- i
		time.Sleep(time.Millisecond)
	}
	close(initChan)
}

func readData(initChan chan int, exitChan chan bool) {
	for {
		val, ok := <-initChan
		if !ok {
			break
		}
		fmt.Println("读取到数据：", val)
	}
	exitChan <- true
	close(exitChan)
}

//select
func Testgoroutines4() {
	intChan := make(chan int, 10)
	for i := 0; i < 10; i++ {
		intChan <- i
	}

	strChan := make(chan string, 10)
	for i := 0; i < 10; i++ {
		strChan <- fmt.Sprintf("%d", i)
	}

	for {
		select {
		case v := <-strChan:
			fmt.Println("strChan-->", v)
		case v := <-intChan:
			fmt.Println("intChan-->", v)
		default:
			fmt.Println("暂未获取到")
			return
		}
	}
}
