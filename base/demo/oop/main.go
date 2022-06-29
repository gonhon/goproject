package main

import (
	"fmt"
	"github.com/limerence-code/goproject/base/oop/student"
)

//引用其他模块
func main() {
	stu := student.New(1, "name")
	fmt.Println("stu:", stu)

}
