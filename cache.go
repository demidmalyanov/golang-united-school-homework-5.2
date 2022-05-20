package cache

import "time"

type Cache struct {
	storage []Key
}

type Key struct {
	key            string
	value          string
	expirationTime time.Time
	willExpired    bool
}

func NewCache() Cache {
	return Cache{}
}

func (c *Cache) Get(key string) (string, bool) {

	for _, item := range c.storage {
		if item.key == key && item.expirationTime.Before(time.Now()) {
			return item.value, true
		}
	}

	return "", false
}

func (c *Cache) Put(key, value string) {
	c.storage = append(c.storage, Key{key: key, value: value, willExpired: false})

}

func (c *Cache) Keys() []string {
	var keys []string

	for _, item := range c.storage {
		if item.expirationTime.Before(time.Now()) && !item.willExpired {
			keys = append(keys, item.key)
		}
	}

	return keys
}

func (c *Cache) PutTill(key, value string, deadline time.Time) {
	c.storage = append(c.storage, Key{key: key, value: value, willExpired: true, expirationTime: deadline})
}
