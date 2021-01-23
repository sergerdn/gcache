package gcache

import (
	"fmt"
	"sync"
)

type BaseIncrementer interface {
	Increment(k string, n int64) (interface{}, error)
}

func incrementValue(v interface{}, n int64) (interface{}, error) {
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
		return nil, fmt.Errorf("the value %v is not an integer", v)
	}
	return v, nil
}

// make sure that LRUIncrementer implements BaseIncrementer
var _ BaseIncrementer = &LRUIncrementer{}

type LRUIncrementer struct {
	cache *LRUCache
	lock  sync.RWMutex
}

func newLRUIncrementer(c *LRUCache) *LRUIncrementer {
	i := &LRUIncrementer{cache: c}
	return i
}

// TODO: change to core implementation and remove custom lock
func (i *LRUIncrementer) Increment(k string, n int64) (interface{}, error) {
	i.lock.Lock()
	defer i.lock.Unlock()
	v, err := i.cache.Get(k)
	if err != nil {
		i.lock.Unlock()
		return nil, fmt.Errorf("item %s not found", k)
	}
	vNew, err := incrementValue(v, n)

	if err != nil {
		return nil, err
	}
	return vNew, i.cache.Set(k, vNew)
}
