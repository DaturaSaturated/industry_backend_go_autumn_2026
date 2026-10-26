package main

type Cache[K comparable, V any] struct {
	capacity int
	items    map[K]V
}

func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	a := make(map[K]V)
	return &Cache[K, V]{items: a, capacity: capacity}
}

func (c *Cache[K, V]) Get(k K) (v V, ok bool) {
	if c.capacity <= 0 {
		var value V
		return value, false
	}
	v, ok = c.items[k]
	return v, ok
}

func (c *Cache[K, V]) Set(k K, v V) bool {
	if c.capacity <= 0 {
		return false
	}

	if len(c.items) < c.capacity {
		c.items[k] = v
		return true
	}

	var ok bool
	_, ok = c.items[k]

	if len(c.items) == c.capacity && ok == true {
		c.items[k] = v
		return true
	}
	return false
}
