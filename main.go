package main

import (
    "net/http"

    "yu"
)

func main() {

    engine := yu.New()
    engine.GET("/index", func(c *yu.Context) {
        c.HTML(http.StatusOK, "<h1>Index Page</h1>")
    })

    v1 := engine.Group("/v1")
    {
        v1.GET("/", func(c *yu.Context) {
            c.HTML(http.StatusOK, "<h1>Hello yu</h1>")
        })

        v1.GET("/hello", func(c *yu.Context) {
            c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
        })
    }
    v2 := engine.Group("/v2")
    {
        v2.GET("/hello/:name", func(c *yu.Context) {
            c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
        })
        v2.POST("/login", func(c *yu.Context) {
            c.JSON(http.StatusOK, yu.H{
                "username": c.PostForm("username"),
                "password": c.PostForm("password"),
            })
        })

    }

    engine.Run(":9999")
}