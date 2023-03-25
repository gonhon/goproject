/*
 * @Author: ljm
 * @Date: 2023-03-25 18:02:32
 * @LastEditTime: 2023-03-25 18:15:48
 */
package web

import (
	"fmt"
	"html/template"
	"net/http"
	"testing"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func TestTmpl(t *testing.T) {
	r := New()
	r.Use(Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHtmlGlob("templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.Get("/", func(c *Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.Get("/students", func(c *Context) {
		c.HTML(http.StatusOK, "arr.tmpl", H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.Get("/date", func(c *Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", H{
			"title": "gee",
			"now":   time.Date(2023, 3, 25, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":9999")
}
