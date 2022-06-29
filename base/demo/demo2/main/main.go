package main

import (
	"fmt"
)

func main() {
	test1()
}

func test1() {
	f := func(a *int, b string) bool {
		fmt.Printf("请输入%v：", b)
		fmt.Scanf("%d", a)
		if *a <= 0 {
			fmt.Println("输入错误")
			return false
		}
		return true
	}

	var y, m, d int
	for {
		if a := f(&y, "年"); a == false {
			continue
		}
		if b := f(&m, "月"); b == false {
			continue
		}
		if c := f(&d, "日"); c == false {
			continue
		}
		fmt.Printf("%d年%d月%d日\n", y, m, d)
	}

}
