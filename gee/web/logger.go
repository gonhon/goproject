package web

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		// 先处理其他的
		c.Next()
		//最后打印耗时
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))

	}
}
