package xclient

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/limerence-code/goproject/gee/rpc/register"
)

type RpcRegisterDiscovery struct {
	*MultiServersDiscovery
	register string
	//超时时间
	timeout time.Duration
	//最后更新时间
	lastUpdate time.Time
}

const defaultUpdateTime = time.Second * 10

func NewRpcRegisterDiscovery(register string, timeout time.Duration) *RpcRegisterDiscovery {
	if timeout == 0 {
		timeout = defaultUpdateTime
	}
	return &RpcRegisterDiscovery{
		MultiServersDiscovery: NewMultiServersDiscovery(make([]string, 0)),
		register:              register,
		timeout:               timeout,
	}
}

func (d *RpcRegisterDiscovery) Update(servers []string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.servers = servers
	d.lastUpdate = time.Now()
	return nil
}

func (d *RpcRegisterDiscovery) Refresh() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.lastUpdate.Add(d.timeout).After(time.Now()) {
		return nil
	}
	log.Println("rpc registry: refresh servers from registry", d.register)
	resp, err := http.Get(d.register)
	if err != nil {
		log.Println("rpc registry refresh err:", err)
		return err
	}
	servers := strings.Split(resp.Header.Get(register.RpcHeader), register.RpcSplit)
	d.servers = make([]string, 0, len(servers))

	for _, serve := range servers {
		if strings.TrimSpace(serve) != "" {
			d.servers = append(d.servers, serve)
		}
	}
	d.lastUpdate = time.Now()

	return nil
}

func (d *RpcRegisterDiscovery) Get(mode SelectMode) (string, error) {
	if err := d.Refresh(); err != nil {
		return "", err
	}
	return d.MultiServersDiscovery.Get(mode)
}

func (d *RpcRegisterDiscovery) GetAll() ([]string, error) {
	if err := d.Refresh(); err != nil {
		return nil, err
	}
	return d.MultiServersDiscovery.GetAll()
}
