package gee

import (
	"fmt"
	"log"
	"strings"
)

type router struct {
	handlers map[string]HandlerFunc
	roots    map[string]*node
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
		roots:    make(map[string]*node),
	}
}

func (r *router) AddRouter(method, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	parts := parsePattent(pattern)
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[getKey(method, pattern)] = handler
}

func (r *router) handle(c *Context) {
	key := getKey(c.Method, c.Path)
	n, params := r.GetRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		r.handlers[key](c)
	} else {
		fmt.Fprintf(c.Writer, "404 NOT FOUND %v\n", c.Req.URL.Path)
	}
}

func (r *router) GetRoute(method, path string) (*node, map[string]string) {
	searchParts := parsePattent(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	node := root.search(searchParts, 0)
	if node != nil {
		parts := parsePattent(node.pattern)
		for i, p := range parts {
			if p[0] == ':' {
				params[p[1:]] = searchParts[i]
			}
			if p[0] == '*' || len(p) > 1 {
				params[p[1:]] = strings.Join(searchParts[i:], "/")
				break
			}
		}
		return node, params
	}
	return nil, nil
}

func parsePattent(pattern string) []string {
	paths := strings.Split(pattern, "/")

	parts := make([]string, len(paths))

	for _, item := range paths {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}
