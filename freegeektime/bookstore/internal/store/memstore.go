/*
 * @Author: gaoh
 * @Date: 2024-07-10 22:38:47
 * @LastEditTime: 2024-07-10 22:38:52
 */
package store

import (
	mystore "bookstore/store"
	factory "bookstore/store/factory"
	"sync"
)

func init() {
	factory.Refister("mem", &MemStore{
		books: make(map[string]*mystore.Book),
	})
}

type MemStore struct {
	sync.RWMutex
	books map[string]*mystore.Book
}

func (ms *MemStore) Create(book *mystore.Book) error {
	ms.Lock()
	defer ms.Unlock()

	if _, ok := ms.books[book.Id]; ok {
		return mystore.ErrExist
	}
	newBook := *book
	ms.books[book.Id] = &newBook

	return nil
}

func (ms *MemStore) Update(*mystore.Book) error {
	return nil
}
func (ms *MemStore) Get(string) (mystore.Book, error) {
	return mystore.Book{}, nil
}
func (ms *MemStore) GetAll() ([]mystore.Book, error) {
	ms.Lock()
	defer ms.Unlock()

	allBooks := make([]mystore.Book, 0, len(ms.books))

	for _, bk := range ms.books {
		allBooks = append(allBooks, *bk)
	}
	return allBooks, nil
}
func (ms *MemStore) Delete(string) error {
	return nil
}
