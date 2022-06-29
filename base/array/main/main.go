package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	test2()
}

func test2() {
	rand.Seed(time.Now().UnixNano())
	var arrays [5]int
	for i, _ := range arrays {
		arrays[i] = rand.Intn(100)
	}

	fmt.Println(arrays)
	var temp int
	lens := len(arrays)
	for i := 0; i < lens/2; i++ {
		temp = arrays[lens-1-i]
		arrays[lens-1-i] = arrays[i]
		arrays[i] = temp
	}
	fmt.Println(arrays)

}

func test1() {
	// var nums [3]int=[3]int {1,23,4}
	//   nums:= [3]int {1,23,4}
	nums := [...]int{1, 23, 4}
	/* for i := 0; i < len(nums); i++ {
		fmt.Printf("%v \t", nums[i])
	} */
	for i, v := range nums {
		fmt.Printf("%v %v \t", i, v)
	}
	fmt.Println()

	var c [26]byte
	for i, _ := range c {
		c[i] = 'A' + byte(i)
	}
	fmt.Println()
	for _, v := range c {
		fmt.Printf("%c \t", v)
	}
}
