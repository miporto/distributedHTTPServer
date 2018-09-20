package cache

import (
	"github.com/manuporto/distributedHTTPServer/pkg/util"
	"sync"
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
	for i := uint(0); i < size; i++ {
		var ce cacheEntry
		slots[i] = &ce
	}
	return Cache{slots, size}
}

func (c Cache) Get(path string) ([]byte, bool) {
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
	pathHash := util.CalculateHash(path)
	slot := c.slots[pathHash%c.size]
	slot.Lock()
	defer slot.Unlock()
	slot.path = path
	slot.body = body
}
