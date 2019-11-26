package imgvips

import "C"

import "sync"

var cStringsCache = &stringsCache{
	cache: make(map[string]*C.char),
}

type stringsCache struct {
	cache map[string]*C.char
	mu    sync.RWMutex
}

func (c *stringsCache) get(key string) *C.char {
	c.mu.RLock()
	if str, ok := c.cache[key]; ok {
		c.mu.RUnlock()

		return str
	}
	c.mu.RUnlock()

	c.mu.Lock()
	str := C.CString(key)
	c.cache[key] = str
	c.mu.Unlock()

	return str
}
