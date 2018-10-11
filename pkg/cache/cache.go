package cache

import "container/list"

type cacheEntry struct {
	path string
	body []byte
}

type Cache struct {
	size  uint
	refs  map[string]*list.Element // *cacheEntry
	nodes *list.List
}

func NewCache(size uint) Cache {
	refs := make(map[string]*list.Element)
	nodes := list.New()
	return Cache{size: size, refs: refs, nodes: nodes}
}

func (cch *Cache) Get(path string) ([]byte, bool) {
	element, ok := cch.refs[path]
	if ok {
		b := element.Value.(*cacheEntry)
		cch.nodes.MoveToFront(element)
		return b.body, ok
	}
	return nil, ok
}

func (cch *Cache) Set(path string, body []byte) {
	if cch.size == 0 {
		return
	}
	element, ok := cch.refs[path]
	if ok {
		element.Value = &cacheEntry{path: path, body: body}
		cch.nodes.MoveToFront(element)
	} else {
		if uint(cch.nodes.Len()) == cch.size {
			elemToEvict := cch.nodes.Back()
			entry := elemToEvict.Value.(*cacheEntry)
			delete(cch.refs, entry.path)
			cch.nodes.Remove(elemToEvict)
		}
		newElement := cch.nodes.PushFront(&cacheEntry{path: path, body: body})
		cch.refs[path] = newElement
	}
}

func (cch *Cache) Delete(path string) {
	element, ok := cch.refs[path]
	if ok {
		delete(cch.refs, path)
		cch.nodes.Remove(element)
	}
}
