/*
 * @Author: gaoh
 * @Date: 2024-07-10 22:40:31
 * @LastEditTime: 2024-07-10 22:59:57
 */
package factory

import (
	"bookstore/store"
	"fmt"
	"sync"
)

var (
	providerMutex sync.Mutex
	providers     = make(map[string]store.Store)
)

func Refister(name string, p store.Store) {
	providerMutex.Lock()
	defer providerMutex.Unlock()

	if p == nil {
		panic("store:refister provider is nuil")
	}

	if _, exist := providers[name]; exist {
		panic("store: provider exist")
	}
	providers[name] = p
}

func New(name string) (store.Store, error) {
	providerMutex.Lock()
	p, ok := providers[name]
	providerMutex.Unlock()
	if !ok {
		return nil, fmt.Errorf("store:unknown provider %s", name)
	}
	return p, nil
}
