package gcache

import (
	"fmt"
	"time"
)

type BaseLRUIncrementer interface {
	Increment(k string, n int64) (interface{}, error)
}

func incrementValue(v interface{}, n int64) error {
	switch v.(type) {
	case int:
		v = v.(int) + int(n)
	case int8:
		v = v.(int8) + int8(n)
	case int16:
		v = v.(int16) + int16(n)
	case int32:
		v = v.(int32) + int32(n)
	case int64:
		v = v.(int64) + n
	case uint:
		v = v.(uint) + uint(n)
	case uintptr:
		v = v.(uintptr) + uintptr(n)
	case uint8:
		v = v.(uint8) + uint8(n)
	case uint16:
		v = v.(uint16) + uint16(n)
	case uint32:
		v = v.(uint32) + uint32(n)
	case uint64:
		v = v.(uint64) + uint64(n)
	case float32:
		v = v.(float32) + float32(n)
	case float64:
		v = v.(float64) + float64(n)
	default:
		return fmt.Errorf("the value %v is not an integer", v)
	}
	return nil
}

// make sure that LRUIncrementer implements BaseLRUIncrementer
var _ BaseLRUIncrementer = &LRUIncrementer{}

type LRUIncrementer struct {
	cache *LRUCache
}

func newLRUIncrementer(c *LRUCache) *LRUIncrementer {
	i := &LRUIncrementer{cache: c}
	return i
}

// Increment an item of type int, int8, int16, int32, int64, uintptr, uint,
// uint8, uint32, or uint64, float32 or float64 by n. Returns an error if the
// item's value is not an integer, if it was not found, or if it is not
// possible to incrementValue it by n.
func (i *LRUIncrementer) Increment(k string, n int64) (interface{}, error) {
	i.cache.mu.Lock()
	item, found := i.cache.items[k]
	if !found {
		i.cache.mu.Unlock()
		return nil, fmt.Errorf("item %s not found", k)
	}

	v := item.Value
	it := v.(*lruItem)
	now := time.Now()
	if it.IsExpired(&now) {
		i.cache.mu.Unlock()
		return nil, fmt.Errorf("item %s not found", k)
	}

	err := incrementValue(&v, n)
	if err != nil {
		i.cache.mu.Unlock()
		return nil, err
	}

	i.cache.items[k].Value = v
	i.cache.mu.Unlock()
	return v, nil
}
