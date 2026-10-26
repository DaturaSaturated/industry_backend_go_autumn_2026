package main

import (
	"container/list"
	"sync"
)

type LRU[K comparable, V any] interface {
	Get(key K) (value V, ok bool)
	Set(key K, value V)
}
type entry[K comparable, V any] struct {
	key   K
	value V
}
type LRUCache[K comparable, V any] struct {
	capacity int
	mu       sync.Mutex
	ll       list.List
	items    map[K]*list.Element
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	items := make(map[K]*list.Element)
	var ll list.List
	return &LRUCache[K, V]{items: items, capacity: capacity, ll: ll}
}

func (c *LRUCache[K, V]) Get(key K) (value V, ok bool) {
	if c.capacity <= 0 {
		var v V
		return v, false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	pointer, b := c.items[key]
	if b == false {
		var v V
		return v, false
	}
	value = (*pointer).Value.(entry[K, V]).value //библиотека list была написана до появления дженериков, поэтому надо явно указывать что за тип я хочу взять в поле Value (т.к. оно any, компилятор не знает что конкретно)
	c.ll.MoveToFront(pointer)
	return value, true
}

func (c *LRUCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	pointer, ok := c.items[key]
	if c.capacity <= 0 {
		return
	}
	if ok {
		(*pointer).Value = entry[K, V]{key: key, value: value}
		return
	}
	if c.ll.Len() < c.capacity {
		c.ll.PushFront(entry[K, V]{key: key, value: value})
		c.items[key] = c.ll.Front()
		return
	}
	if c.ll.Len() >= c.capacity {
		c.ll.PushFront(entry[K, V]{key: key, value: value})
		c.items[key] = c.ll.Front()
		pointer := c.ll.Back()
		k := (*pointer).Value.(entry[K, V]).key
		delete(c.items, k)
		c.ll.Remove(c.ll.Back())
		return
	}
}
