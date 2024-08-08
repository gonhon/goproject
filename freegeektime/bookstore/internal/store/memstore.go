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

func (ms *MemStore) Update(book *mystore.Book) error {
	ms.Lock()
	defer ms.Unlock()

	oldBook, exist := ms.books[book.Id]
	if !exist {
		return mystore.ErrNotFound
	}
	nBook := *oldBook
	if book.Author != nil {
		nBook.Author = book.Author
	}
	if book.Name != "" {
		nBook.Name = book.Name
	}
	if book.Press != "" {
		nBook.Press = book.Press
	}
	ms.books[book.Id] = &nBook

	return nil
}
func (ms *MemStore) Get(id string) (mystore.Book, error) {
	ms.RLocker().Lock()
	defer ms.RLocker().Unlock()

	book, exist := ms.books[id]
	if !exist {
		return mystore.Book{}, mystore.ErrNotFound
	}

	return *book, nil
}
func (ms *MemStore) GetAll() ([]mystore.Book, error) {
	ms.RLocker().Lock()
	defer ms.RLocker().Unlock()

	allBooks := make([]mystore.Book, 0, len(ms.books))

	for _, bk := range ms.books {
		allBooks = append(allBooks, *bk)
	}
	return allBooks, nil
}
func (ms *MemStore) Delete(id string) error {
	ms.Lock()
	defer ms.Unlock()

	if _, exist := ms.books[id]; !exist {
		return mystore.ErrNotFound
	}
	delete(ms.books, id)
	return nil
}
