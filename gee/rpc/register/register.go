package register

import (
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type RpcRegister struct {
	timeout time.Duration
	mutex   sync.Mutex
	servers map[string]*ServerItem
}

type ServerItem struct {
	Addr  string
	start time.Time
}

const (
	defaultPath = "/rpc/register"
	//默认超时时间
	defaultTimeout = time.Minute * 5
	RpcHeader      = "X-Rpc-Server"
	RpcSplit       = ","
)

func New(timeout time.Duration) *RpcRegister {
	return &RpcRegister{timeout: timeout, servers: make(map[string]*ServerItem)}
}

var DefaultRpcRegister = New(defaultTimeout)

// 添加实例到注册中心
func (r *RpcRegister) putServer(addr string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	s, exist := r.servers[addr]
	if !exist {
		r.servers[addr] = &ServerItem{Addr: addr, start: time.Now()}
	} else {
		s.start = time.Now()
	}
}

// 获取可用的服务
func (r *RpcRegister) aliveServers() (alive []string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for addr, s := range r.servers {
		if r.timeout == 0 || s.start.Add(r.timeout).After(time.Now()) {
			alive = append(alive, addr)
		} else {
			//超时移除服务
			delete(r.servers, addr)
		}
	}
	sort.Strings(alive)
	return
}

func (r *RpcRegister) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		w.Header().Set(RpcHeader, strings.Join(r.aliveServers(), RpcSplit))
	case "POST":
		addr := req.Header.Get(RpcHeader)
		if addr == "" {
			w.WriteHeader(http.StatusInternalServerError)
		}
		r.putServer(addr)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (r *RpcRegister) HandleHTTP(registerPath string) {
	http.Handle(registerPath, r)
	log.Println("rpc register path:", registerPath)
}

func HandleHTTP() {
	DefaultRpcRegister.HandleHTTP(defaultPath)
}

// 心跳检查
func Heartbeat(register, addr string, duration time.Duration) {
	if duration == 0 {
		duration = defaultTimeout - time.Duration(1)*time.Minute
	}
	var err error
	err = sendHeartbeat(register, addr)
	go func() {
		t := time.NewTicker(duration)
		for err == nil {
			<-t.C
			err = sendHeartbeat(register, addr)
		}
	}()
}

func sendHeartbeat(register, addr string) error {
	log.Println(addr, " send heart beat to register ", register)
	httpClient := &http.Client{}
	req, _ := http.NewRequest("POST", register, nil)
	req.Header.Set(RpcHeader, addr)
	if _, err := httpClient.Do(req); err != nil {
		log.Panicln("rpc server heart beat err:", err)
		return err
	}
	return nil
}
