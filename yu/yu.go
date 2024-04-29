package yu

import (
    "net/http"
)

// 定义 handerFunc 类型，以便复用
type HandlerFunc func(*Context)

// 定义结构体：Engine，实现 http.Handler 接口的 ServeHTTP 方法
type Engine struct{
    router *Router
}

// 构造方法
func New() *Engine {
    return &Engine{router: newRouter()}
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