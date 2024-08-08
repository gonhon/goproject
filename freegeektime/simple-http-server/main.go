package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	test02()
}

func test01() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world!"))
	})

	http.ListenAndServe(":8080", nil)
}

func validateAuth(s string) error {
	if s != "123456" {
		return fmt.Errorf("%s", "bad auth token")
	}
	return nil
}
func greetings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}
func logHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t)
		h.ServeHTTP(w, r)
	})
}
func authHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := validateAuth(r.Header.Get("auth"))
		if err != nil {
			http.Error(w, "bad auth param", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func test02() {
	http.ListenAndServe(":8080", logHandler(authHandler(http.HandlerFunc(greetings))))
}

type a func(a string, b int)

func (f a) AA(a string, b int) {
	f(a, b)
}
