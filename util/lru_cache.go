package common

import "runtime"

type LruCache struct {
	data map[string]*node
	len  int
	cap  int
	head *node
	tail *node
}

type node struct {
	pre   *node
	next  *node
	key   string
	value interface{}
}

func NewLruCache(size int) *LruCache {
	head := &node{}
	tail := &node{}
	head.next = tail
	tail.pre = head

	return &LruCache{
		data: make(map[string]*node),
		len:  0,
		cap:  size,
		head: head,
		tail: tail,
	}
}

func (c *LruCache) Get(key string) interface{} {
	v, ok := c.data[key]
	if !ok {
		return nil
	}
	c.moveToHead(v)
	return v.value
}

func (c *LruCache) GetAll() map[string]interface{} {
	list := make(map[string]interface{})
	for key, node := range c.data {
		list[key] = node.value
	}
	return list
}

func (c *LruCache) Set(key string, value interface{}) {
	v, ok := c.data[key]
	if ok {
		v.value = value
		c.moveToHead(v)
		return
	}
	v = &node{
		pre:   c.head,
		next:  c.head.next,
		key:   key,
		value: value,
	}

	c.data[key] = v
	c.head.next = v
	v.next.pre = v

	c.len++
	if c.len > c.cap {
		c.removeTail()
		c.len--
	}
}

func (c *LruCache) CleanData() {
	c.head.next = c.tail
	c.tail.pre = c.head
	count := c.len
	c.len = 0
	c.data = make(map[string]*node)
	if count > 10000 {
		runtime.GC()
	}
}

func (c *LruCache) Len() int {
	return c.len
}

func (c *LruCache) moveToHead(v *node) {
	v.pre.next = v.next
	v.next.pre = v.pre

	v.next = c.head.next
	c.head.next = v
	v.pre = c.head
}

func (c *LruCache) removeTail() {
	v := c.tail.pre
	delete(c.data, v.key)

	c.tail.pre = v.pre
	v.pre.next = c.tail
}
