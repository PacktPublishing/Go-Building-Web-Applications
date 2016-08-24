package cache

import (
	"strings"
)

const (
	Location = "/var/cache/"
)

type CacheItem struct {
	TTL int
	Key string
}

func newCache(endpoint string, params ...[]string) CacheItem {
	cacheName := endponit + "_" + strings.Join(params, "_")
	c := CacheItem{}
	return c
}

func (c CacheItem) Get() (bool, string) {

	stats, err := os.Stat(c.Key)
	if err != nil {
		return false, ""
	}

	age := time.Nanoseconds() - stats.ModTime()
	if age <= c.TTL {
		cache, _ := ioutil.ReadFile(c.Key)
		return true, cache
	} else {
		return false, ""
	}
}

func (c CacheItem) Set(data []byte) bool {
	err := ioutil.WriteFile(c.Key, data, 0644)
}

func (c CacheItem) Clear() bool {

}
