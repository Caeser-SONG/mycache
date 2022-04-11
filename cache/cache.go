package cache

import (
	"cache/lru"
	"fmt"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Printf("44444 ====  %p", &c)
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
		fmt.Printf("33333 ==== %p", c.lru)
	}
	c.lru.Add(key, value)
}
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Printf("333333333 === %p", c.lru)
	if c.lru == nil {
		return
	}
	fmt.Println("zzzzzzzzzzzzzzaaaaaaaaaaaaa")
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return

}
