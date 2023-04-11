package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/limerence-code/goproject/gee/cache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *cache.Group {
	return cache.NewGroup("scores", 2<<10, cache.GetterFunc(func(key string) ([]byte, error) {
		log.Println("db search key", key)
		if val, ok := db[key]; ok {
			return []byte(val), nil
		} else {
			return nil, fmt.Errorf("%s is not found", key)
		}
	}))

}

func startCacheServer(addr string, addrs []string, group *cache.Group) {
	poll := cache.NewHttpPoll(addr)
	poll.Set(addrs...)
	group.RegisterPeers(poll)
	log.Println("cache is run ...", addr)
	log.Fatal(http.ListenAndServe(addr[7:], poll))
}

func startApiServer(apiAddr string, group *cache.Group) {
	http.Handle("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		view, err := group.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(view.ByteSlice())
	}))
	log.Println("api is run ...", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "cache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}
	group := createGroup()

	if api {
		go startApiServer(apiAddr, group)
	}
	startCacheServer(addrMap[port], addrs, group)
}
