package web

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix     string
	middleware HandlerFunc
	parent     *RouterGroup
	engine     *Engine
}
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

func getKey(method, pattern string) string {
	return fmt.Sprintf("%s-%s", method, pattern)
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: engine.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) AddRouter(method, pattern string, handler HandlerFunc) {
	group.engine.router.AddRouter(method, group.prefix+pattern, handler)
}

func (group *RouterGroup) Get(pattern string, handler HandlerFunc) {
	group.AddRouter("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.AddRouter("POST", pattern, handler)
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
