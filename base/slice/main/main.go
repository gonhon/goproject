package main

import "fmt"

func main() {
	/* arrays := [...]int{1, 2, 345, 6, 75}
	slice := arrays[1:2]
	fmt.Println(slice)
	fmt.Printf("len:%v cap:%v\n", len(slice), cap(slice))

	se := make([]int, 10)
	fmt.Println(se)
	se = append(se, 20, 30)
	// se[10] = 14
	fmt.Println(se[10]) */
	/* 	fb := fbl(10)
	   	fmt.Println(fb)
	   	fmt.Println(fbs(3)) */

	//定义切片
	array := []int{1, 2, 3, 4, 5}

	array = append(array, 9)
	fmt.Println(array)
}

func fbl(n int) []int {
	fbls := make([]int, n)
	fbls[0] = 1
	fbls[1] = 1
	for i := 2; i < n; i++ {
		fbls[i] = fbls[i-1] + fbls[i-2]
	}
	return fbls
}

func fbs(n int) int {
	if n == 1 || n == 2 {
		return 1
	}
	return fbs(n-1) + fbs(n-2)
}
