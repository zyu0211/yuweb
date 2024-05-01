package yu

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
    return func(c *Context) {
        t := time.Now()     // Start timer
        c.Next()            // Process request
        log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))  		// Calculate resolution time
	}
}