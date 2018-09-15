package filemap

import (
	"errors"
	"sync"
)

const (
	GET    = iota
	POST   = iota
	PUT    = iota
	DELETE = iota
)

type FileMap struct {
	sync.RWMutex
	m map[string]sync.RWMutex
}

func (fm *FileMap) insert(filename string) (*sync.RWMutex, error) {
	fm.Lock()
	defer fm.Unlock()
	if _, ok := fm.m[filename]; ok {
		return nil, errors.New("")
	}
	fm.m[filename] = sync.RWMutex{}
	l := fm.m[filename]

	return &l, nil
}

func (fm *FileMap) get(filename string) (*sync.RWMutex, error) {
	fm.RLock()
	defer fm.RUnlock()
	if l, ok := fm.m[filename]; ok {
		return &l, nil
	}
	return nil, errors.New("")
}

func (fm *FileMap) delete(filename string) (*sync.RWMutex, error) {
	fm.Lock()
	defer fm.Unlock()
	if l, present := fm.m[filename]; present {
		delete(fm.m, filename)
		return &l, nil
	}
	return nil, errors.New("Key not present")
}
