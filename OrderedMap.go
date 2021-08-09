package utils

import (
	"reflect"
	"sync"
)

// OrderedItems tracks the items
type OrderedItems struct {
	KeySlice     []interface{}
	OrderedItems map[interface{}]interface{}
	mu           sync.Mutex
}

// OrderedItem represents a stored item
type OrderedItem struct {
	Index int
	Key   interface{}
	Value interface{}
}

// NewOrderedMap creates a new instance of OrderedItems
func NewOrderedMap() (items *OrderedItems) {
	items = new(OrderedItems)
	items.OrderedItems = make(map[interface{}]interface{})
	return
}

// Iter creates a channel to interates over OrderedItems
func (i *OrderedItems) Iter() <-chan OrderedItem {
	ch := make(chan OrderedItem)
	go func() {
		defer close(ch)
		for index, key := range i.KeySlice {
			val, ok := i.OrderedItems[key]
			if ok {
				ch <- OrderedItem{index, key, val}
			}
		}
	}()
	return ch
}

// Add adds an ordered item
func (i *OrderedItems) Add(key interface{}, value interface{}) interface{} {
	i.mu.Lock()
	defer i.mu.Unlock()
	_, ok := i.OrderedItems[key]
	i.OrderedItems[key] = value
	if !ok {
		i.KeySlice = append(i.KeySlice, key)
	}
	return i.OrderedItems[key]
}

// Update updes an existing ordered item, returning the value
func (i *OrderedItems) Update(key interface{}, value interface{}) interface{} {
	i.mu.Lock()
	defer i.mu.Unlock()
	if _, ok := i.OrderedItems[key]; ok {
		i.OrderedItems[key] = value
		return i.OrderedItems[key]
	}
	return struct{}{}
}

// Get returns an ordered item
func (i *OrderedItems) Get(key interface{}) (interface{}, bool) {
	value, ok := i.OrderedItems[key]
	return value, ok
}

// Has checks if an ordered item exists
func (i *OrderedItems) Has(key interface{}) bool {
	_, ok := i.OrderedItems[key]
	return ok
}

// Del removes an ordered item
func (i *OrderedItems) Del(key interface{}) bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	delete(i.OrderedItems, key)
	for id, val := range i.KeySlice {
		if val == key {
			i.KeySlice = append(i.KeySlice[:id], i.KeySlice[id+1:]...)
			return true
		}
	}
	return false
}

// Count gets the total number of ordered items
func (i *OrderedItems) Count() int {
	return len(i.OrderedItems)
}

// Keys returns all the ordered item keys
func (i *OrderedItems) Keys() []interface{} {
	keys := make([]interface{}, 0, len(i.OrderedItems))
	for k := range i.OrderedItems {
		keys = append(keys, k)
	}
	return keys
}

// StringKeys returns only ordered item keys that are strings
func (i *OrderedItems) StringKeys() []string {
	keys := make([]string, 0, len(i.OrderedItems))
	for k := range i.OrderedItems {
		if reflect.TypeOf(k).Kind() == reflect.String {
			keys = append(keys, k.(string))
		}
	}
	return keys
}
