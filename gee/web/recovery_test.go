/*
 * @Author: gaoh
 * @Date: 2023-03-25 18:32:57
 * @LastEditTime: 2023-03-25 18:35:15
 */
package web

import (
	"net/http"
	"testing"
)

func TestRecovery(t *testing.T) {
	r := Default()
	r.Get("/", func(c *Context) {
		c.String(http.StatusOK, "Hello World\n")
	})
	// index out of range for testing Recovery()
	r.Get("/panic", func(c *Context) {
		names := []string{"Hello"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}
