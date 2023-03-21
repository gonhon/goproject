package main

import (
	"fmt"
	"net/http"
)

type Engine struct {
}

func (engine *Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	switch req.URL.Path {
	case "/":
		fmt.Fprintf(resp, "URL Path: %v\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(resp, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(resp, "404 NOT FOUND: %s\n", req.URL)
	}
}

func main() {
	engine := new(Engine)
	fmt.Println("http starting 8080...")
	http.ListenAndServe(":8080", engine)
}
