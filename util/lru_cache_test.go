package common

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestLruCache_Set(t *testing.T) {
	cache := NewLruCache(4)

	cache.Set("1", 1)
	cache.Set("2", 2)
	assert.Equal(t, "1", cache.data["2"].next.key)

	cache.Set("3", 3)
	assert.Equal(t, "2", cache.data["3"].next.key)

	cache.Set("4", 4)
	assert.Equal(t, "3", cache.data["4"].next.key)

	cache.Set("5", 5)
	assert.Equal(t, 4, cache.Len())
	assert.Equal(t, "4", cache.data["5"].next.key)

	assert.Equal(t, []string{"5", "4", "3", "2"}, getKeyOrder(cache))

	cache.Set("2", 20)
	assert.Equal(t, []string{"2", "5", "4", "3"}, getKeyOrder(cache))

	cache.Set("1", 10)
	assert.Equal(t, []string{"1", "2", "5", "4"}, getKeyOrder(cache))
}

func TestLruCache_Get(t *testing.T) {
	cache := NewLruCache(4)

	cache.Set("1", 1)
	cache.Set("2", 2)
	cache.Set("3", 3)
	cache.Set("4", 4)
	cache.Set("5", 5)
	assert.Equal(t, []string{"5", "4", "3", "2"}, getKeyOrder(cache))

	cache.Get("3")
	assert.Equal(t, []string{"3", "5", "4", "2"}, getKeyOrder(cache))

	cache.Get("2")
	assert.Equal(t, []string{"2", "3", "5", "4"}, getKeyOrder(cache))
}

func getKeyOrder(cache *LruCache) []string {
	list := make([]string, 0)
	v := cache.head.next
	for v.next != nil {
		list = append(list, v.key)
		v = v.next
	}
	return list
}
