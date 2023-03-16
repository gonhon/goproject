package base

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func convert() {
	intface := map[int]string{
		0: "0",
		1: "1",
		2: "2",
	}
	Convert(intface)
}

func Convert(params interface{}) {
	maps := params.(map[int]string)
	val, ok := maps[0]
	if !ok {
		panic("值不存在")
	}
	fmt.Printf("val:%s\n", val)
}

//--------------------chan--------------------

func chanFunc() {
	var c chan int = make(chan int, 3)
	c <- 1
	c <- 1
	c <- 1
	val := <-c
	fmt.Printf("chan val:%d", val)
}

func chanChange() {
	var chans chan int = make(chan int, 2)
	go func() {
		for i := 0; i < 10; i++ {
			chans <- i
			fmt.Printf("prod send chan params:%d\n", i)
		}
		fmt.Printf("close chan")
		close(chans)
	}()

	for {
		val, ok := <-chans
		if !ok {
			fmt.Printf("chan cloase")
			break
		}
		fmt.Printf("custemer val:%d\n", val)
	}

	fmt.Println("End.")
}

//--------------------select--------------------
func selectFunc() {
	intChan := [3]chan int{
		make(chan int, 1),
		make(chan int, 2),
		make(chan int, 3),
	}
	index := rand.Intn(3)
	intChan[index] <- index

	select {
	case <-intChan[0]:
		fmt.Printf("chan 0 \n")
	case <-intChan[1]:
		fmt.Printf("chan 1 \n")
	case <-intChan[2]:
		fmt.Printf("chan 2 \n")
	default:
		fmt.Printf("not val ....")
	}

}

func asycService() chan string {
	retCh := make(chan string, 1)
	go func() {
		time.Sleep(time.Duration(1) * time.Minute)
		fmt.Println("return result.")
		retCh <- "val"
		fmt.Println("service exited.")
	}()
	return retCh
}

//--------------------闭包--------------------
type opertion func(x, y int) (int, error)

func exec(x, y int, op opertion) (int, error) {
	if op == nil {
		return 0, errors.New("opertion is nil")
	}
	return op(x, y)
}

type calculateFunc func(x, y int) (int, error)

func genCalculate(op opertion) calculateFunc {
	return func(x, y int) (int, error) {
		if op == nil {
			return 0, errors.New("opertion is nil")
		}
		return op(x, y)
	}
}
