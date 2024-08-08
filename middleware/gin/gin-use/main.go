package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func test01() {
	var sl = [][]int{
		{1, 34, 26, 35, 78},
		{3, 45, 13, 24, 99},
		{101, 13, 38, 7, 127},
		{54, 27, 40, 83, 81},
	}
outerloop:
	for i := 0; i < len(sl); i++ {
		for j := 0; j < len(sl[i]); j++ {
			if sl[i][j] == 13 {
				fmt.Printf("found 13 at [%d, %d]\n", i, j)
				break outerloop
			}
		}
	}

	var m = map[string]int{
		"tony": 21,
		"tom":  22,
		"jim":  23,
	}
	counter := 0
	for k, v := range m {
		if counter == 0 {
			delete(m, "tony")
		}
		counter++
		fmt.Println(k, v)
	}
	fmt.Println("counter is ", counter)

	counter = 0
	for k, v := range m {
		if counter == 0 {
			m["lucy"] = 24
		}
		counter++
		fmt.Println(k, v)
	}
	fmt.Println("counter is ", counter)
}

type field struct {
	name string
}

func (f *field) printf() {
	fmt.Println(f.name)
}

func test02() {

	var f1 = []*field{{"one"}, {"two"}, {"three"}}

	for _, v := range f1 {
		go v.printf()
	}

	var f2 = []*field{{"four"}, {"five"}, {"six"}}

	for _, v := range f2 {
		go v.printf()
	}

	time.Sleep(3 * time.Second)

}
