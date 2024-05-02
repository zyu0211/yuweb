package main

import (
	"fmt"
	"net/http"
	"html/template"
	"time"

	"yu"
)

type student struct {
    Name string
    Age int8
}

func FormatAsDate(t time.Time) string {
    year, month, day := t.Date()
    return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
    engine := yu.New()
    engine.Use(yu.Logger()) // global midlleware
    engine.SetFuncMap(template.FuncMap{
        "FormatAsDate": FormatAsDate,
    })

    engine.LoadHTMLGlob("templates/*")
    engine.Static("/assets", "./static")

    stu1 := &student{Name: "zhangsan", Age: 20}
    stu2 := &student{Name: "lisi", Age: 22}

    engine.GET("/", func(c *yu.Context) {
        c.HTML(http.StatusOK, "css.tmpl", nil)
    })

    engine.GET("/students", func(c *yu.Context) {
        c.HTML(http.StatusOK, "arr.tmpl", yu.H{
            "title":  "yu",
            "stuArr": [2]*student{stu1, stu2},
        })
    })

    engine.GET("/date", func(c *yu.Context) {
        c.HTML(http.StatusOK, "custom_func.tmpl", yu.H{
            "title": "yu",
            "now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
        })
    })

    engine.Run(":9999")
}