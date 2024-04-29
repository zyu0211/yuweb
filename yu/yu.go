package yu

import (
	"log"
	"net/http"
)

// 定义 handerFunc 类型，以便复用
type HandlerFunc func(*Context)

type RouterGroup struct {
    prefix string
    middleWares []HandlerFunc
    engine *Engine
}

/*
    定义结构体：Engine，实现 http.Handler 接口的 ServeHTTP 方法
*/
type Engine struct {
    *RouterGroup
    router *Router
    groups []*RouterGroup
}

/*
    构造方法
*/
func New() *Engine {

    engine := &Engine{router: newRouter()}
    engine.RouterGroup = &RouterGroup{engine: engine}
    engine.groups = []*RouterGroup{engine.RouterGroup}
    return engine
}

func (engine *Engine) addRoute(method string, path string, handler HandlerFunc) {
    engine.router.addRoute(method, path, handler)
}

func (engine *Engine) GET(path string, handler HandlerFunc) {
    engine.router.addRoute("GET", path, handler)
}

func (engine *Engine) POST(path string, handler HandlerFunc) {
    engine.router.addRoute("POST", path, handler)
}

func (engine *Engine) Run(addr string) (err error) {
    return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    c := newContext(w, req)
    engine.router.handle(c)
}

/*
    构造一个RouterGroup
*/
func (group *RouterGroup) Group(prefix string) *RouterGroup {
    engine := group.engine
    newGroup := &RouterGroup {
        prefix: prefix,
        engine: engine,     // 新的Group和原Group共用一个 engine
    }
    engine.groups = append(engine.groups, newGroup)

    return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
    path := group.prefix + comp
    log.Printf("Route %4s - %s", method, path)
    group.engine.router.addRoute(method, path, handler)
}

func (group *RouterGroup) GET(path string, handler HandlerFunc) {
    group.addRoute("GET", path, handler)
}

func (group *RouterGroup) POST(path string, handler HandlerFunc) {
    group.addRoute("POST", path, handler)
}