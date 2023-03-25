package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"
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
	port := 8080
	log.Printf("start... %d", port)
	r.Run(":" + strconv.Itoa(port))
}

func TestGroupRun(t *testing.T) {
	r := New()
	r.Get("/index", func(c *Context) {
		c.Html(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := r.Group("/v1")
	{
		v1.Get("/", func(c *Context) {
			c.Html(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		v1.Get("/hello", func(c *Context) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.Get("/hello/:name", func(c *Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *Context) {
			c.Json(http.StatusOK, H{
				"username": c.PostFrom("username"),
				"password": c.PostFrom("password"),
			})
		})

	}

	r.Run(":9999")

}

func onlyForV2() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func TestMiddlewaresRun(t *testing.T) {
	r := New()
	r.Use(Logger()) // global midlleware
	r.Get("/", func(c *Context) {
		c.Html(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.Get("/hello/:name", func(c *Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	r.Run(":9999")
}

//静态文件解析
func TestStaticRun(t *testing.T) {
	r := New()
	//解析当前文件夹
	r.Static("/assets", "./")
	//http://127.0.0.1:9999/assets/gee.go
	//访问当前下的 gee.go
	r.Run(":9999")
}
