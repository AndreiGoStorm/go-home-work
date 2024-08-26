package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	sync.Mutex
}

type Item struct {
	key   Key
	value interface{}
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	defer lru.Mutex.Unlock()
	lru.Mutex.Lock()
	newItem := NewItem(key, value)
	item, ok := lru.items[key]
	if ok {
		item.Value = newItem
		lru.moveToFront(item, key)
	} else {
		lru.dequeueOldestItem()
		lru.items[key] = lru.queue.PushFront(newItem)
	}

	return ok
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	defer lru.Mutex.Unlock()
	lru.Mutex.Lock()
	if item, ok := lru.items[key]; ok {
		lru.moveToFront(item, key)
		currentItem := item.Value.(*Item)
		return currentItem.value, true
	}

	return nil, false
}

func (lru *lruCache) dequeueOldestItem() {
	if lru.queue.Len() >= lru.capacity {
		last := lru.queue.Back()
		item := last.Value.(*Item)
		delete(lru.items, item.key)
		lru.queue.Remove(last)
	}
}

func (lru *lruCache) moveToFront(item *ListItem, key Key) {
	lru.queue.MoveToFront(item)
	lru.items[key] = lru.queue.Front()
}

func (lru *lruCache) Clear() {
	defer lru.Mutex.Unlock()
	lru.Mutex.Lock()
	clear(lru.items)
	lru.queue = NewList()
}

func NewItem(key Key, value interface{}) *Item {
	return &Item{key, value}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
