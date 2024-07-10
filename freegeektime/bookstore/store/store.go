/*
 * @Author: gaoh
 * @Date: 2024-07-10 22:40:44
 * @LastEditTime: 2024-07-10 22:47:47
 */
package store

type Book struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Author []string `json:"author"`
	Press  string   `json:"press"`
}

type Store interface {
	Create(*Book) error
	Update(*Book) error
	Get(string) error
	GetAll() ([]Book, error)
	Delete(string) error
}
