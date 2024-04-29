package main

import (
    "net/http"

    "yu"
)

func main() {

    engine := yu.New()

    engine.GET("/", func(c *yu.Context) {
        c.HTML(http.StatusOK, "<h1>Hello yu!</h1>")
    })

    engine.GET("/hello", func(c *yu.Context) {
        c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
    })

    engine.POST("/login", func(c *yu.Context) {
        c.JSON(http.StatusOK, yu.H{
            "username": c.PostForm("username"),
            "password": c.PostForm("password"),
        })
    })

    engine.Run(":9999")
}