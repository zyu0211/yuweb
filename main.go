package main

import (
	"net/http"
	"yu"
)



func main() {
    engine := yu.Default()
    engine.GET("/", func(c *yu.Context) {
		c.String(http.StatusOK, "Hello yuweb\n")
	})
	// index out of range for testing Recovery()
	engine.GET("/panic", func(c *yu.Context) {
		names := []string{"yuweb"}
		c.String(http.StatusOK, names[100])
	})

    engine.Run(":9999")
}