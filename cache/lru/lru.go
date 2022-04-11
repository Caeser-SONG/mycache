package lru

import (
	"container/list" // 双向链表
	"fmt"
)

type Cache struct {
	maxBytes  int64                    // 最大内存
	nbytes    int64                    // 当前内存
	ll        *list.List               // 链表
	cache     map[string]*list.Element // 链表的元素类型的指针
	OnEvicted func(string, Value)      // 某条记录被移除时的回调函数，可以为 nil。
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int //占用内存的大小
}

// 创建一个cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry) // 类型转换？？
		return kv.value, true
	}
	fmt.Printf("now: %v", c.nbytes)
	return
}
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	// 如果key存在于缓存里,更新并移动到队尾
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	fmt.Printf("max: %v", c.maxBytes)

	fmt.Printf("now: %v", c.nbytes)
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}
func (c *Cache) Len() int {
	return c.ll.Len()
}
