package gee

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.AddRouter("GET", "/", nil)
	r.AddRouter("GET", "/hello/:name", nil)
	r.AddRouter("GET", "/hello/b/c", nil)
	r.AddRouter("GET", "/hi/:name", nil)
	r.AddRouter("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", "name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestRouter(t *testing.T) {
	r := newTestRouter()
	log.Printf("%v", r)
	// n, ps := r.GetRoute("GET", "/hello/geektutu")
	n, ps := r.GetRoute("GET", "/hello/:name")

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if ps["name"] != "geektutu" {
		t.Fatal("name should be equal to 'geektutu'")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])

}

func TestRun(t *testing.T) {
	r := New()
	r.Get("/", func(c *Context) {
		c.Html(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.Get("/hello", func(c *Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.Get("/hello/:name", func(c *Context) {
		// expect /hello/geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.Get("/assets/*filepath", func(c *Context) {
		c.Json(http.StatusOK, H{"filepath": c.Param("filepath")})
	})
	r.Run(":9999")
}
