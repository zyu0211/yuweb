package yu

import (
	"net/http"
	"strings"
)


type Router struct {
    // roots存储每种请求方式的Trie树根节点，key为 "GET", "POST" 等
    roots map[string]*Node
    // handlers存储每个请求的HandlerFunc，key为 "POST-/p/book" 等
    handlers map[string]HandlerFunc // 路由表
}

func newRouter() *Router {
    return &Router{
        roots: make(map[string]*Node),
        handlers: make(map[string]HandlerFunc),
    }
}

func parsePath(path string) []string {
    path_split := strings.Split(path, "/")

    parts := make([]string, 0)
    for _, item := range path_split {
        if item != "" {
            parts = append(parts, item)
            if item[0] == '*' {
                break
            }
        }
    }
    return parts
}

func (r *Router) addRoute(method string, path string, handler HandlerFunc) {

    parts := parsePath(path)
    _, ok := r.roots[method]
    if !ok {
        r.roots[method] = &Node{}
    }
    r.roots[method].insert(path, parts, 0)

    key := method + "-" + path
    r.handlers[key] = handler
}

func (r *Router) getRoute(method string, path string) (*Node, map[string]string) {

    searchParts := parsePath(path)
    params := make(map[string]string)

    root, ok := r.roots[method]
    if !ok {
        return nil, nil
    }

    node := root.search(searchParts, 0)
    if node != nil {
        parts := parsePath(node.path)
        for i, part := range parts {
            if part[0] == ':' {
                params[part[1:]] = searchParts[i]
            }
            if part[0] == '*' {
                params[part[1:]] = strings.Join(searchParts[i:], "/")
            }
        }
        return node, params
    }

    return nil, nil
}

func (r *Router) handle(c *Context) {
    node, params := r.getRoute(c.Method, c.Path)

    if node != nil {
        c.Params = params
        key := c.Method + "-" + node.path
        c.handlers = append(c.handlers, r.handlers[key])
    } else {
        c.handlers = append(c.handlers, func(c *Context){
            c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
        })
    }

    c.Next()
}