package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

// 下载对账文件
func tradeBillDownload(date string) {
	client := &http.Client{}
	url := fmt.Sprintf("http://222.188.64.226:18877/main/onlinehall/pay/billConfirm/test/tradeBillDownload?confirmDate=%s", url.QueryEscape(date+" 12:00:00"))
	log.Println("url:", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request for A:", err)
		return
	}
	req.Header.Set("Authorization", "Bearer 794767ad-c5fa-4f8c-88ab-2770e06bb48f") // 添加 Authorization 头部

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error calling A:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Called A with date:", date, "Status:", resp.Status)
}

// 执行对账
func tradeBillConfirm(date string) {
	client := &http.Client{}
	url := fmt.Sprintf("http://222.188.64.226:18877/main/onlinehall/pay/billConfirm/test/tradeBillConfirm?confirmDate=%s", url.QueryEscape(date+" 12:00:00"))
	log.Println("url:", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request for B:", err)
		return
	}
	req.Header.Set("Authorization", "Bearer 794767ad-c5fa-4f8c-88ab-2770e06bb48f") // 添加 Authorization 头部

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error calling B:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Called B with date:", date, "Status:", resp.Status)
}

func main() {
	startDate := time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 2, 20, 0, 0, 0, 0, time.UTC)

	for date := startDate; date.Before(endDate) || date.Equal(endDate); date = date.AddDate(0, 0, 1) {
		dateStr := date.Format("2006-01-02")

		tradeBillDownload(dateStr)
		time.Sleep(1 * time.Second)
		tradeBillConfirm(dateStr)
	}
}
