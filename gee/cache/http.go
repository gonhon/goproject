package cache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/cache/"

type HttpPoll struct {
	self     string
	basePath string
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
	// group :=
	// key :=

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
