package gee

import (
	"fmt"
	"log"
	"testing"
)

/* func TestEngine(t *testing.T) {
	engine := New()
	engine.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL.Path: %v\n", r.URL.Path)
	})

	engine.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%v]: %v\n", k, v)
		}
	})

	engine.Run(":8080")
} */

func TestEngine(t *testing.T) {
	r := New()
	r.Get("/", func(c *Context) {
		fmt.Fprintf(c.Writer, "URL.Path: %v\n", c.Req.URL.Path)
	})

	r.Get("/hello", func(c *Context) {
		for k, v := range c.Req.Header {
			fmt.Fprintf(c.Writer, "Header[%v]: %v\n", k, v)
		}
	})

	r.POST("/login", func(ctx *Context) {
		ctx.Json(200, H{
			"username": ctx.PostFrom("username"),
			"password": ctx.PostFrom("password"),
		})
	})
	log.Fatalf("start... 8080")
	r.Run(":9090")
}
