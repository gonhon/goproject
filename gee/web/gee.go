package web

import (
	"fmt"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(*Context)

//路由组
type RouterGroup struct {
	prefix string
	//middleware
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

//维护路由关系
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

//RouterGroup和Engine使用了双向继承,双方都可以使用彼此的属性方法

func getKey(method, pattern string) string {
	return fmt.Sprintf("%s-%s", method, pattern)
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	//创建默认的组
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

//创建一个组  追加到engine后面
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		//把父类的前缀也加上
		prefix: engine.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) AddRouter(method, pattern string, handler HandlerFunc) {
	//路由添加组的前缀
	group.engine.router.AddRouter(method, group.prefix+pattern, handler)
}

func (group *RouterGroup) Get(pattern string, handler HandlerFunc) {
	group.AddRouter("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.AddRouter("POST", pattern, handler)
}

//添加middleware
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

//静态文件映射
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absulutePath := path.Join(group.prefix, relativePath)
	fileSystem := http.StripPrefix(absulutePath, http.FileServer(fs))
	return func(ctx *Context) {
		file := ctx.Params["filepath"]
		if _, err := fs.Open(file); err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}
		fileSystem.ServeHTTP(ctx.Writer, ctx.Req)
	}
}

//添加文件解析hander
func (group *RouterGroup) Static(relativePath, root string) {
	group.Get(path.Join(relativePath, "/*filepath"),
		group.createStaticHandler(relativePath, http.Dir(root)))
}

//Http路由都会走这里
func (engine *Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		//根据group前缀找到对应的middleware
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	//"/favicon.ico"
	//将http信息包装conext
	conext := newContext(resp, req)
	conext.middleware = middlewares
	//执行hander
	engine.router.handle(conext)
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}
