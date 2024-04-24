package yu

import (
	"fmt"
	"net/http"
)

// 定义 handerFunc 类型，以便复用
type HanderFunc func(http.ResponseWriter, *http.Request)

// 定义结构体：Engine，实现 http.Handler 接口的 ServeHTTP 方法
type Engine struct{
	router map[string]HanderFunc
}

// 构造方法
func New() *Engine {
    return &Engine{router: make(map[string]HanderFunc)}
}

func (engine *Engine) addRoute(method string, path string, handler HanderFunc) {
    key := method + "-" + path
	engine.router[key] = handler
}

func (engine *Engine) GET(path string, handler HanderFunc) {
    engine.addRoute("GET", path, handler)
}

func (engine *Engine) POST(path string, handler HanderFunc) {
    engine.addRoute("POST", path, handler)
}

func (engine *Engine) Run(addr string) (err error) {
    return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    key := req.Method + "-" + req.URL.Path
    if handler, ok := engine.router[key]; ok {
        handler(w, req)
    } else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}