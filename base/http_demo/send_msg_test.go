package httpdemo

import (
	"fmt"
	"testing"
)

func TestMain(m *testing.T) {
	page := 1
	for {
		aResponse, err := fetchA(page)
		if err != nil {
			fmt.Printf("Error fetching A API: %v\n", err)
			break
		}

		if len(aResponse.Clients) == 0 {
			break // 如果没有更多数据，退出循环
		}

		for _, client := range aResponse.Clients {
			if err := callB(client.ClientCode); err != nil {
				fmt.Printf("Error calling B API for client %s: %v\n", client.ClientCode, err)
			}
		}

		// 检查是否还有更多页面
		if len(aResponse.Clients) < PageSize {
			break // 若当前页面的数量小于 PageSize，说明已到最后一页
		}
		page++
	}
}
