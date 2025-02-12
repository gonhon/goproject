package httpdemo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	AAPIURL  = "https://example.com/api/a" // A接口的URL
	BAPIURL  = "https://example.com/api/b" // B接口的URL
	PageSize = 10                          // 每页的大小
)

type Client struct {
	ClientCode string `json:"clientCode"`
}

type AResponse struct {
	Clients []Client `json:"clients"`
	Total   int      `json:"total"`
}

func fetchA(page int) (AResponse, error) {
	var response AResponse
	resp, err := http.Get(fmt.Sprintf("%s?page=%d&size=%d", AAPIURL, page, PageSize))
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, err
	}
	return response, nil
}

func callB(clientCode string) error {
	// 根据 clientCode 调用 B 接口的逻辑
	resp, err := http.Post(BAPIURL, "application/json", nil) // 根据实际情况构造请求体
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 处理B接口的响应
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("B API call failed with status: %s", resp.Status)
	}
	return nil
}
