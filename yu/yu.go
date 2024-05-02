package yu

import (
	"log"
	"net/http"
	"path"
	"strings"
	"html/template"
)

// 定义 handerFunc 类型，以便复用
type HandlerFunc func(*Context)

/*
    定义结构体：Engine，实现 http.Handler 接口的 ServeHTTP 方法
*/
type Engine struct {
    *RouterGroup
    router *Router
    groups []*RouterGroup

    htmlTemplates *template.Template
    funcMap template.FuncMap
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

/*
    内置中间件
*/
func Default() *Engine {
    engine := New()
    engine.Use(Logger(), Recover())
    return engine
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
    engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(path string) {
    engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(path))
}

func (engine *Engine) addRoute(method string, path string, handler HandlerFunc) {
    engine.router.addRoute(method, path, handler)
}

func (engine *Engine) GET(path string, handler HandlerFunc) {
    engine.addRoute("GET", path, handler)
}

func (engine *Engine) POST(path string, handler HandlerFunc) {
    engine.addRoute("POST", path, handler)
}

func (engine *Engine) Run(addr string) (err error) {
    return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    var middleWares []HandlerFunc
    for _, group := range engine.groups {
        if strings.HasPrefix(req.URL.Path, group.prefix) {
            middleWares = append(middleWares, group.middleWares...)
        }
    }

    c := newContext(w, req)
    c.handlers = middleWares
    c.engine = engine
    engine.router.handle(c)
}

/*
    路由分组
*/
type RouterGroup struct {
    prefix string
    middleWares []HandlerFunc
    engine *Engine
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

/*
    向RouterGroup中注册中间件
*/
func (group *RouterGroup) Use(middleWares ...HandlerFunc) {
    group.middleWares = append(group.middleWares, middleWares...)
}

/*
    静态资源处理器
*/
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
    absolutePath := path.Join(group.prefix, relativePath)
    fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

    return func(c *Context) {
        file := c.Params["filepath"]
        if _, err := fs.Open(file); err != nil {
            c.Status(http.StatusNotFound)
            return
        }
        fileServer.ServeHTTP(c.Writer, c.Req)
    }
}

/*
    映射静态资源
*/
func (group *RouterGroup) Static(relativePath string, root string) {
    handler := group.createStaticHandler(relativePath, http.Dir(root))
    urlPath := path.Join(relativePath, "/*filepath")
    group.GET(urlPath, handler)
}