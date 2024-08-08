/*
 * @Author: gaoh
 * @Date: 2024-07-10 22:39:56
 * @LastEditTime: 2024-07-10 22:40:00
 */
package server

import (
	"bookstore/server/middleware"
	"bookstore/store"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type BookStoreServer struct {
	s   store.Store
	srv *http.Server
}

func NewBookStoreServer(addr string, s store.Store) *BookStoreServer {
	srv := &BookStoreServer{
		s: s,
		srv: &http.Server{
			Addr: addr,
		},
	}
	route := mux.NewRouter()
	route.HandleFunc("/book", srv.createBookHandler).Methods("POST")
	route.HandleFunc("/book/{id}", srv.getBookHandler).Methods("GET")
	route.HandleFunc("/book/{id}", srv.deleteBookHandler).Methods("DELETE")
	route.HandleFunc("/book/{id}", srv.updateBookHandler).Methods("PUT")
	route.HandleFunc("/book", srv.getAllBookHandler).Methods("GET")

	srv.srv.Handler = middleware.Logging(middleware.Validating(route))
	return srv

}

func (bs *BookStoreServer) ListenAndServe() (<-chan error, error) {
	var err error
	errChan := make(chan error)
	go func() {
		bs.srv.ListenAndServe()
		errChan <- err
	}()

	select {
	case err = <-errChan:
		return nil, err
	case <-time.After(time.Second):
		return errChan, nil
	}
}

func (bs *BookStoreServer) Shutdown(ctx context.Context) error {
	return bs.srv.Shutdown(ctx)
}

func getBodyBook(body io.ReadCloser) (*store.Book, error) {
	dec := json.NewDecoder(body)
	var book store.Book
	if err := dec.Decode(&book); err != nil {
		return nil, err
	}
	return &book, nil
}

func (bs *BookStoreServer) createBookHandler(w http.ResponseWriter, req *http.Request) {
	book, _ := getBodyBook(req.Body)

	if err := bs.s.Create(book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (bs *BookStoreServer) updateBookHandler(w http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(w, "id is nil", http.StatusBadRequest)
		return
	}
	book, _ := getBodyBook(req.Body)

	book.Id = id

	if err := bs.s.Update(book); err != nil {
		http.Error(w, "update errr", http.StatusBadRequest)
		return
	}

}

func (bs *BookStoreServer) getBookHandler(w http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(w, "id is nil", http.StatusBadRequest)
		return
	}
	book, err := bs.s.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response(w, book)

}

func (bs *BookStoreServer) getAllBookHandler(w http.ResponseWriter, req *http.Request) {

	books, err := bs.s.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response(w, books)
}

func (bs *BookStoreServer) deleteBookHandler(w http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(w, "id is nil", http.StatusBadRequest)
		return
	}
	if err := bs.s.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func response(w http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
