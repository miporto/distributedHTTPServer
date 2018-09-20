package lockpool

import (
	"github.com/manuporto/distributedHTTPServer/pkg/util"
	"sync"
)

type LockPool struct {
	locks []*sync.RWMutex
	size  uint
}

func NewLockPool(size uint) LockPool {
	locks := make([]*sync.RWMutex, size)
	for i := uint(0); i < size; i++ {
		var l sync.RWMutex
		locks[i] = &l
	}
	return LockPool{locks, size}
}

func (lp LockPool) GetLock(path string) *sync.RWMutex {
	pathHash := util.CalculateHash(path)
	return lp.locks[pathHash%lp.size]
}
