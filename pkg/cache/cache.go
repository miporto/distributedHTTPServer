package cache

import (
	"sync"

	"github.com/manuporto/distributedHTTPServer/pkg/util"
)

type cacheEntry struct {
	sync.RWMutex
	path string
	body []byte
}

type Cache struct {
	l     *sync.RWMutex
	slots []*cacheEntry
	size  uint
}

func NewCache(size uint) Cache {
	slots := make([]*cacheEntry, size)
	var l sync.RWMutex
	if size == 0 {
		return Cache{&l, slots, size}
	}
	for i := uint(0); i < size; i++ {
		var ce cacheEntry
		slots[i] = &ce
	}
	return Cache{&l, slots, size}
}

func (c *Cache) Get2(path string) ([]byte, bool) {
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

func (c Cache) Get(path string) ([]byte, bool) {
	if c.size == 0 {
		return nil, false
	}
	pathHash := util.CalculateHash(path)
	slot := c.slots[pathHash%c.size]
	slot.RLock()
	defer slot.RUnlock()
	if path != slot.path {
		return nil, false
	}
	return slot.body, true
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
}

func (c *Cache) Update2(path string, body []byte) {
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

func (c Cache) Update(path string, body []byte) {
	if c.size == 0 {
		return
	}
	pathHash := util.CalculateHash(path)
	slot := c.slots[pathHash%c.size]
	slot.Lock()
	defer slot.Unlock()
	slot.path = path
	slot.body = body
}

func (c *Cache) Delete2(path string) {
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

func (c Cache) Delete(path string) {
	if c.size == 0 {
		return
	}
	pathHash := util.CalculateHash(path)
	slot := c.slots[pathHash%c.size]
	slot.Lock()
	defer slot.Unlock()
	slot.path = ""
	slot.body = nil
}
