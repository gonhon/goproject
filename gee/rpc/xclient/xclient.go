package xclient

import (
	"io"
	"sync"

	"github.com/limerence-code/goproject/gee/rpc"
)

type XClient struct {
	d       Discovery
	mode    SelectMode
	opt     *rpc.Option
	mutex   sync.Mutex
	clients map[string]*rpc.Client
}

var _ io.Closer = (*XClient)(nil)

func NewXClient(d Discovery, mode SelectMode, opt *rpc.Option) *XClient {
	return &XClient{d: d, mode: mode, opt: opt, clients: make(map[string]*rpc.Client)}
}

func (xc *XClient) Close() error {
	xc.mutex.Lock()
	defer xc.mutex.Unlock()

	for key, client := range xc.clients {
		client.Close()
		delete(xc.clients, key)
	}
	return nil
}
