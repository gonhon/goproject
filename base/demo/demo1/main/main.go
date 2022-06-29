package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var a int64 = 123
	var b int64 = 456
	fmt.Printf("a %v %T size %d \t", a, a, unsafe.Sizeof(a))
	fmt.Printf("b %v %T size %d \n", b, b, unsafe.Sizeof(b))
	a, b = b, a
	fmt.Printf("a %v %T size %d \t", a, a, unsafe.Sizeof(a))
	fmt.Printf("b %v %T size %d \n", b, b, unsafe.Sizeof(b))
}
