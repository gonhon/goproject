package main

import "fmt"

type st1 struct {
	A int
	B st2
	D *st2
}

type st2 struct {
	C string
}

func main() {
	temp := st1{}
	fmt.Println(temp)
	temp2 := new(st1)
	fmt.Println(temp2)
}
