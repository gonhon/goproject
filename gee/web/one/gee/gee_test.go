package gee

import (
	"fmt"
	"net/http"
	"testing"
)

func TestEngine(t *testing.T) {
	engine := New()
	engine.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL.Path: %v\n", r.URL.Path)
	})

	engine.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%v]: %v\n", k, v)
		}
	})

	engine.Run(":8080")
}
