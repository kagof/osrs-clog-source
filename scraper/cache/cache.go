package cache

import "scraper/log"

// Cache a naive in-memory cache with a maximum size; if the size exceeded then old values roll off.
type Cache[K comparable, V any] struct {
	m      map[K]V
	keys   []K
	oldest int
	cap    int
}

func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity < 1 {
		capacity = 1
	}
	return &Cache[K, V]{
		m:      make(map[K]V, capacity),
		keys:   make([]K, 0, capacity),
		oldest: 0,
		cap:    capacity,
	}
}

func (c *Cache[K, V]) Memoized(f func(k K) (V, error)) func(k K) (V, error) {
	return func(k K) (V, error) {
		cached, ok := c.Get(k)
		if ok {
			log.Debugf("found cached value for key %v\n", k)
			return cached, nil
		}
		v, err := f(k)
		if err != nil {
			return v, err
		}
		c.Put(k, v)
		return v, nil
	}
}

func (c *Cache[K, V]) Get(k K) (V, bool) {
	v, ok := c.m[k]
	return v, ok
}

func (c *Cache[K, V]) Put(k K, v V) {
	// oldest key logic becomes wrong if we have multiple puts for the same key, but that's ok, we just accept it
	_, hasVal := c.m[k]
	if hasVal {
		c.m[k] = v
		return
	}

	if len(c.m) < c.cap {
		c.m[k] = v
		c.keys = append(c.keys, k)
		return
	}

	oldestKey := c.keys[c.oldest]
	delete(c.m, oldestKey)
	c.keys[c.oldest] = k
	c.oldest = (c.oldest + 1) % c.cap
	c.m[k] = v
}

func (c *Cache[K, V]) Len() int {
	return len(c.m)
}

func (c *Cache[K, V]) findNextOldest() {
	if len(c.m) == 0 {
		c.oldest = 0
		return
	}
	c.oldest = (c.oldest + 1) % c.cap
}
