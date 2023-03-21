package gee

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

func getKey(method, pattern string) string {
	return fmt.Sprintf("%s-%s", method, pattern)
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) AddRouter(method, pattern string, handler HandlerFunc) {
	engine.router[getKey(method, pattern)] = handler
}

func (engine *Engine) Get(pattern string, handler HandlerFunc) {
	engine.AddRouter("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.AddRouter("POST", pattern, handler)
}

func (engine *Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if handler, ok := engine.router[getKey(req.Method, req.URL.Path)]; ok {
		handler(resp, req)
	} else {
		fmt.Fprintf(resp, "404 NOT FOUND %v\n", req.URL.Path)
	}
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}
