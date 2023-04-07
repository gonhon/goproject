package cache

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/cache/"
	defaultReplicas = 50
)

type HttpPoll struct {
	self        string
	basePath    string
	mutex       sync.Mutex
	peers       *Map
	httpGetters map[string]*httpGetter
}

func NewHttpPoll(self string) *HttpPoll {
	return &HttpPoll{self: self, basePath: defaultBasePath}
}

func (p *HttpPoll) Log(format string, args ...interface{}) {
	log.Printf("server %s: %s", p.self, fmt.Sprintf(format, args...))
}

func (p *HttpPoll) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HttpPoll serving unexpected path: " + r.URL.Path)
	}

	p.Log("%s %s", r.Method, r.URL.Path)

	paths := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)

	if len(paths) != 2 {
		http.Error(w, "request failed", http.StatusBadRequest)
		return
	}

	group := GetGroup(paths[0])
	if group == nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}

	view, err := group.Get(paths[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//写出文件
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}

func (p *HttpPoll) Set(peers ...string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.peers = NewMap(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))

	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseUrl: peer + p.basePath}
	}
}

func (p *HttpPoll) PickPeer(key string) (PeerGetter, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %v", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

// 确保这个类型实现了这个接口 如果没有实现会报错的

var _ PeerPicker = (*HttpPoll)(nil)

type httpGetter struct {
	baseUrl string
}

func (h *httpGetter) Get(group, key string) ([]byte, error) {
	u := fmt.Sprintf("%v%v/%v", h.baseUrl, url.QueryEscape(group), url.QueryEscape(key))
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if bytes, err := ioutil.ReadAll(res.Body); err != nil {
		return bytes, nil
	} else {
		return nil, err
	}
}

// 确保这个类型实现了这个接口 如果没有实现会报错的
var _ PeerGetter = (*httpGetter)(nil)
