package main

import (
    "log"
    "net/http"
    "time"

    "yu"
)

func onlyForV2() yu.HandlerFunc {
    return func(c *yu.Context) {
        // Start timer
        t := time.Now()
        // if a server error occurred
        c.HTML(500, "Internal Server Error")
        // Calculate resolution time
        log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
    }
}

func main() {
    engine := yu.New()
    engine.Use(yu.Logger()) // global midlleware
    engine.GET("/", func(c *yu.Context) {
        c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
    })

    v2 := engine.Group("/v2")
    v2.Use(onlyForV2()) // v2 group middleware
    {
        v2.GET("/hello/:name", func(c *yu.Context) {
            // expect /hello/geektutu
            c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
        })
    }

    engine.Run(":9999")
}