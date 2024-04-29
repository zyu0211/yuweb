package yu

import "net/http"


type Router struct {
    handlers map[string]HandlerFunc // 路由表
}

func newRouter() *Router {
    return &Router{
        handlers: make(map[string]HandlerFunc),
    }
}

func (r *Router) addRoute(method string, path string, handler HandlerFunc) {
    key := method + "-" + path
    r.handlers[key] = handler
}

func (r *Router) handle(c *Context) {
    key := c.Method + "-" + c.Path
    if handler, ok := r.handlers[key]; ok {
        handler(c)
    } else {
        c.String(http.StatusNotFound, "404 Not Found: %s", c.Path)
    }
}