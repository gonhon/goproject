package gee

import (
	"fmt"
	"log"
)

type router struct {
	handers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		handers: make(map[string]HandlerFunc),
	}
}

func (r *router) AddRouter(method, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	r.handers[getKey(method, pattern)] = handler
}

func (r *router) handle(c *Context) {
	key := getKey(c.Method, c.Path)
	if handler, ok := r.handers[key]; ok {
		handler(c)
	} else {
		fmt.Fprintf(c.Writer, "404 NOT FOUND %v\n", c.Req.URL.Path)
	}
}
