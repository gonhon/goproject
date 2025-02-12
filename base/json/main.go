package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	inputName  = flag.String("input", "data.json", "源 JSON 文件名")
	outputName = flag.String("output", fmt.Sprintf("%s.json", time.Now().Format("2006-01-02-150405")), "目标 JSON 文件名")
)

/**
{"clientCode":"1020071","impactClientCode":"1020071","batchId":"1339001861449060352","cciId":"4ab6475d5961621f7cf77f289d0839e0"}
转换
"{\"clientCode\":\"1020071\",\"impactClientCode\":\"1020071\",\"batchId\":\"1339001861449060352\",\"cciId\":\"4ab6475d5961621f7cf77f289d0839e0\"}"
**/

func main() {
	flag.Parse() // 解析命令行参数
	log.Printf("inputName:%s,outputName:%s\n", *inputName, *outputName)
	output(*inputName, *outputName)
}

func print(fileName string) {
	// 读取 JSON 文件
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// 定义一个空的切片来存储解码后的 JSON 数据
	var jsonData []interface{}

	// 解码 JSON 数据
	if err := json.Unmarshal(data, &jsonData); err != nil {
		log.Fatal(err)
	}

	// 创建一个新的切片用于存储转义后的字符串
	var escapedJSON []string

	// 遍历解码后的数据并将每个对象转义
	for _, item := range jsonData {
		// 将每个对象编码为字符串
		jsonStr, err := json.Marshal(item)
		if err != nil {
			log.Fatal(err)
		}
		// 转义字符串并添加到新的切片中
		escapedJSON = append(escapedJSON, string(jsonStr))
	}

	// 输出最终的转义字符串数组
	output, err := json.Marshal(escapedJSON)
	if err != nil {
		log.Fatal(err)
	}

	// 输出结果
	fmt.Println(string(output))
}

func output(inputName, outputName string) {
	// 读取源 JSON 文件
	data, err := ioutil.ReadFile(inputName)
	if err != nil {
		log.Fatal(err)
	}

	// 定义一个空的切片来存储解码后的 JSON 数据
	var jsonData []interface{}

	// 解码 JSON 数据
	if err := json.Unmarshal(data, &jsonData); err != nil {
		log.Fatal(err)
	}

	// 创建一个新的切片用于存储转义后的字符串
	var escapedJSON []string

	// 遍历解码后的数据并将每个对象转义
	for _, item := range jsonData {
		// 将每个对象编码为字符串
		jsonStr, err := json.Marshal(item)
		if err != nil {
			log.Fatal(err)
		}
		// 转义字符串并添加到新的切片中
		escapedJSON = append(escapedJSON, string(jsonStr))
	}

	// 输出最终的转义字符串数组
	output, err := json.Marshal(escapedJSON)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))

	// 写入到目标 JSON 文件
	err = ioutil.WriteFile(outputName, output, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// 打印成功消息
	fmt.Printf("转义后的 JSON 已写入到 %s 文件中。\n", outputName)
}
