package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
)

func main() {
	//test01()
	test02()
}
func test01() {
	// maps := make(map[string]string, 2)
	maps := map[string]string{
		"a": "666",
	}
	maps["a"] = "123"
	fmt.Println(maps)
	a := new(int)
	fmt.Println(a)

	k, v := maps["b"]
	fmt.Printf("Key:%v,Val:%v", k, v)
}

var intMap map[int]int
var cnt = 8192

func test02() {
	mutex := sync.Mutex{}
	mutex.Lock()
	mutex.Unlock()

	printMemStats() //打印出memory情况

	initMap()       // 创建map
	runtime.GC()    // 强制执行GC
	printMemStats() // 在强制GC之后 再打印memory情况

	log.Println(len(intMap)) // 查看map的元素个数
	for i := 0; i < cnt; i++ {
		delete(intMap, i) // 执行delete的操作
	}
	log.Println(len(intMap)) // 验证执行delete操作对实际map的元素个数影响

	runtime.GC() // 强制执行GC
	printMemStats()

	intMap = nil // 将map置为nil 释放其占用的内存空间
	runtime.GC()
	printMemStats()
}

func initMap() {
	intMap = make(map[int]int, cnt)

	for i := 0; i < cnt; i++ {
		intMap[i] = i
	}
}

func printMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("Alloc = %v TotalAlloc = %v Sys = %v NumGC = %v\n", m.Alloc/1024, m.TotalAlloc/1024, m.Sys/1024, m.NumGC)
}
