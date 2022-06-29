package main

import (
	"errors"
	"fmt"
	_ "time"
)

func main() {
	// var scan=
	//go env -w GOPROXY = "https://proxy.golang.com.cn,direct"
	//go env -w GOPRIVATE=git.mycompany.com,github.com/my/private
	for i := 0; i < 10; i++ {
		fmt.Print(i, " ")
	}
	a := len("123")
	fmt.Printf("%v\n", a)

	test()

	e := test1(2)
	if e != nil {
		//  fmt.Printf("errMsg%v,type%T",e,e)
		panic(e)
		// fmt.Println("err====>",e)
	}
	fmt.Println("end...")
}
func test() {
	defer func() {

		if err := recover(); err != nil {
			fmt.Println("err:", err)
		}
	}()
	b := 0
	a := 2
	c := a / b
	fmt.Println("res:", c)
}

func test1(a int32) (err error) {
	if a == 2 {
		return errors.New("err....")
	}
	return nil
}
