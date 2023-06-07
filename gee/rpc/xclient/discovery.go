package xclient

import (
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"
)

type SelectMode int

const (
	RandomSelect SelectMode = iota
	RoundRobinSelect
)

// 服务发现注册接口
type Discovery interface {
	Refresh() error
	Update(servers []string) error
	Get(mode SelectMode) (string, error)
	GetAll() ([]string, error)
}

type MultiServersDiscovery struct {
	r       *rand.Rand
	mutex   sync.Mutex
	servers []string
	index   int
}

func NewMultiServersDiscovery(servers []string) *MultiServersDiscovery {
	msd := &MultiServersDiscovery{
		servers: servers,
		//设置随机种子
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	//index随机设置
	msd.index = msd.r.Intn(math.MaxInt32 - 1)
	return msd
}

var _ Discovery = (*MultiServersDiscovery)(nil)

func (d *MultiServersDiscovery) Refresh() error {
	return nil
}

func (d *MultiServersDiscovery) Update(servers []string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.servers = servers
	return nil
}

func (d *MultiServersDiscovery) Get(mode SelectMode) (string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	lens := len(d.servers)
	if lens == 0 {
		return "", errors.New("rpc discovery no available servers")
	}

	switch mode {
	case RandomSelect:
		return d.servers[d.r.Intn(lens)], nil
	case RoundRobinSelect:
		s := d.servers[d.index%lens]
		d.index = (d.index + 1) % lens
		return s, nil
	default:
		return "", errors.New("rpc discovery no supported select mode")
	}
}

func (d *MultiServersDiscovery) GetAll() ([]string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	servers := make([]string, len(d.servers), len(d.servers))
	copy(servers, d.servers)
	return servers, nil
}
