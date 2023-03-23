package web

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(*Context)

type Engine struct {
	router *router
}

func getKey(method, pattern string) string {
	return fmt.Sprintf("%s-%s", method, pattern)
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) AddRouter(method, pattern string, handler HandlerFunc) {
	engine.router.AddRouter(method, pattern, handler)
}

func (engine *Engine) Get(pattern string, handler HandlerFunc) {
	engine.AddRouter("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.AddRouter("POST", pattern, handler)
}

func (engine *Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	//"/favicon.ico"
	conext := newContext(resp, req)
	engine.router.handle(conext)
	// engine.router.handle(newContext(resp, req))
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}
