package main

import (
	"cache"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"tom":  "123",
	"sam":  "2342",
	"jack": "777",
}

func createGroup() *cache.Group {
	return cache.NewGroup("scores", 2<<10, cache.GetterFunc(func(key string) ([]byte, error) {
		log.Println("search DB", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))
}

func startCacheServer(addr string, addrs []string, g *cache.Group) {
	peers := cache.NewHttpPool(addr)
	peers.Set(addrs...)
	g.RegisterPeers(peers)
	log.Println("cache is running at ", addr)
	log.Fatal(http.ListenAndServe(addr[7:], nil))
}

func startAPIserver(apiAddr string, g *cache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			fmt.Printf("111 == %p \n", g)
			view, err := g.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	http.Handle("/set", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			value := r.URL.Query().Get("value")
			fmt.Printf("11111 == %p \n", g)
			v := g.Set(key, value)
			w.Write(v.ByteSlice())
		}))
	log.Println("fonted server is running at ", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "cache server port")
	flag.BoolVar(&api, "api", false, "Start a api server")
	flag.Parse()
	fmt.Println(api)
	apiAddr := "http://0.0.0.0:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}
	var addrs []string

	for _, v := range addrMap {
		addrs = append(addrs, v)

	}

	g := createGroup()
	if api {
		go startAPIserver(apiAddr, g)
	}
	startCacheServer(addrMap[port], []string(addrs), g)
}
