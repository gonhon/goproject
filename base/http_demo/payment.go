package httpdemo

import (
	"fmt"
	"net/http"
	"time"
)

func Payment() {
	startDate := "20240601"

	// 解析开始日期
	start, err := time.Parse("20060102", startDate)
	if err != nil {
		fmt.Println("Error parsing start date:", err)
		return
	}

	// 创建一个日期切片
	var currentDate time.Time
	for currentDate = start; currentDate.Before(time.Now()) || currentDate.Equal(time.Now()); currentDate = currentDate.AddDate(0, 0, 1) {
		dateStr := currentDate.Format("20060102")
		url := fmt.Sprintf("http://218.202.93.230:8877/main/payment/api/weixin/downloadWeixinBill?date=%s", dateStr)

		fmt.Printf("url:%s\n", url)

		// 发起 HTTP 请求
		resp, err := http.Post(url, "application/json;charset=UTF-8", nil)
		if err != nil {
			fmt.Printf("Error fetching data for date %s: %v\n", dateStr, err)
			continue
		}
		defer resp.Body.Close()

		// 检查响应状态
		if resp.StatusCode == http.StatusOK {
			fmt.Printf("Successfully fetched data for date: %s\n", dateStr)
		} else {
			fmt.Printf("Failed to fetch data for date %s: Status Code %d\n", dateStr, resp.StatusCode)
		}
	}
}
