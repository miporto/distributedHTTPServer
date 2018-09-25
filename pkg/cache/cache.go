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
	slots []*cacheEntry
	size  uint
}

func NewCache(size uint) Cache {
	slots := make([]*cacheEntry, size)
	if size == 0 {
		return Cache{slots, size}
	}
	for i := uint(0); i < size; i++ {
		var ce cacheEntry
		slots[i] = &ce
	}
	return Cache{slots, size}
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
