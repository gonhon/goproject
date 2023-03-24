package web

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
	parts := parsePattern(pattern)
	if _, ok := r.roots[method]; !ok {
		//已请求类型作为key作为根节点
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[getKey(method, pattern)] = handler
}

func (r *router) handle(c *Context) {
	n, params := r.GetRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		r.handlers[getKey(c.Method, n.pattern)](c)
	} else {
		c.middleware = append(c.middleware, func(c *Context) {
			fmt.Fprintf(c.Writer, "404 NOT FOUND %v\n", c.Req.URL.Path)
		})
	}
	c.Next()
}

func (r *router) GetRoute(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	//根据方法获取根节点
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	//根据切割后的uri找到匹配的结点
	node := root.search(searchParts, 0)
	if node != nil {
		parts := parsePattern(node.pattern)
		for i, p := range parts {
			if p[0] == ':' {
				params[p[1:]] = searchParts[i]
			}
			if p[0] == '*' && len(p) > 1 {
				params[p[1:]] = strings.Join(searchParts[i:], "/")
				break
			}
		}
		return node, params
	}
	return nil, nil
}

// 将uri按/切割 转成数组 直到'*'截至
func parsePattern(pattern string) []string {
	parts := make([]string, 0)
	strs := strings.Split(pattern, "/")

	for _, item := range strs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}
