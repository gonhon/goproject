package web

import (
	"fmt"
	"log"
	"net/http"
)

func Run() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL Path: %v\n", r.URL.Path)
	})

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q]=%q\n", k, v)
		}
	})
	fmt.Println("http starting 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

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

func RunServeHTTP() {
	engine := new(Engine)
	fmt.Println("http starting 8080...")
	http.ListenAndServe(":8080", engine)
}
