package cache

import (
	"sync"
)

type cacheEntry struct {
	path string
	body []byte
}

type Cache struct {
	l       *sync.RWMutex
	slots   []*cacheEntry
	size    uint
	toEvict uint
}

func NewCache(size uint) Cache {
	slots := make([]*cacheEntry, size)
	var l sync.RWMutex
	if size == 0 {
		return Cache{&l, slots, size, 0}
	}
	for i := uint(0); i < size; i++ {
		var ce cacheEntry
		slots[i] = &ce
	}
	return Cache{&l, slots, size, 0}
}

func (c *Cache) Get(path string) ([]byte, bool) {
	if c.size == 0 {
		return nil, false
	}
	c.l.RLock()
	defer c.l.RUnlock()
	for _, element := range c.slots {
		if element.path == "" {
			break
		} else if element.path == path {
			return element.body, true
		}
	}
	return nil, false
}

func (c *Cache) Insert(path string, body []byte) {
	if c.size == 0 {
		return
	}
	c.l.Lock()
	defer c.l.Unlock()
	for _, element := range c.slots {
		if element.path == "" {
			element.path = path
			element.body = body
			return
		}
	}
	c.slots[c.toEvict%c.size].path = path
	c.slots[c.toEvict%c.size].body = body
	c.toEvict++
	if c.toEvict == c.size {
		c.toEvict = 0
	}
}

func (c *Cache) Update(path string, body []byte) {
	if c.size == 0 {
		return
	}
	c.l.Lock()
	defer c.l.Unlock()
	for _, element := range c.slots {
		if element.path == path {
			element.body = body
			return
		}
	}
}

func (c *Cache) Delete(path string) {
	if c.size == 0 {
		return
	}
	c.l.Lock()
	defer c.l.Unlock()
	for _, element := range c.slots {
		if element.path == path {
			element.path = ""
			element.body = nil
			return
		}
	}
}
