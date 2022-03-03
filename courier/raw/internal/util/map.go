package util

import (
	"sync"
)

// ConcurrentMap is an concurrent-safe wrapper over a key-value map.
type ConcurrentMap struct {
	store map[interface{}]interface{}
	mutex sync.Mutex
}

// NewConcurrentMap creates a new concurrent-safe map.
func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{store: make(map[interface{}]interface{})}
}

// Set sets a key-value pair.
func (c *ConcurrentMap) Set(key interface{}, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.store[key] = value
}

// Get gets the value associated with a key, if available.
func (c *ConcurrentMap) Get(key interface{}) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	value, ok := c.store[key]

	return value, ok
}

// Delete deletes the value associated with a key, if available.
func (c *ConcurrentMap) Delete(key interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.store, key)
}

// Size returns the current number of entries in the map.
func (c *ConcurrentMap) Size() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return len(c.store)
}

// MergeMaps takes one or more maps as input and merges them, in sequential order, into a new map.
// Note that later values take precedence in the merge, and that a cloned map is returned.
func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	for _, item := range maps {
		for k, v := range item {
			merged[k] = v
		}
	}

	return merged
}
