package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
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
