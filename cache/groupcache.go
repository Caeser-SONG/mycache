package cache

import (
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (g GetterFunc) Get(key string) ([]byte, error) {
	return g(key)
}

// group是一个完整的cache单元 name+maincache+节点
type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	fmt.Printf("22222 ===  %p  \n", &g.mainCache)
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[cache] hit")
		return v, nil
	}
	return g.load(key)
}
func (g *Group) Set(key string, value string) ByteView {
	mu.Lock()
	defer mu.Unlock()
	b := ByteView{
		b: []byte(value),
	}
	fmt.Printf("22222 ===  %p  \n", &g.mainCache)
	g.mainCache.add(key, b)
	return b
}

// 远程下载key的内容
func (g *Group) load(key string) (value ByteView, err error) {
	if g.peers != nil {
		if peer, ok := g.peers.PickPeer(key); ok {
			if value, err := g.getFromPeer(peer, key); err == nil {
				return value, nil
			}
			log.Printf("[cache] failed to get from peer", err)
		}
	}
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	// use getter
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	// join data into cache
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	fmt.Println("add to cache")
	g.mainCache.add(key, value)
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	//  如果还没有注册peerpicker
	if g.peers != nil {
		panic("RegisterPeers call more than once")
	}
	g.peers = peers
}
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}
