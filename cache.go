package cache

import (
	"time"
)

type Key struct {
	key          string
	value        string
	expTime      time.Time
	shouldExpire bool
}

type Cache struct {
	storage []Key
}

func NewCache() Cache {
	return Cache{}
}

func (c *Cache) Get(key string) (string, bool) {

	for i := range c.storage {
		if c.storage[i].key == key && !c.storage[i].shouldExpire {
			return c.storage[i].value, true
		} else if c.storage[i].key == key && c.storage[i].expTime.After(time.Now()) {
			return c.storage[i].value, true
		}

	}

	return "", false
}

func (c *Cache) Put(key, value string) {

	validKeys := c.Keys()
	needToOverwrite := false

	for _, validKey := range validKeys {
		if validKey == key {
			needToOverwrite = true
			break
		}
	}

	if needToOverwrite {
		for i := range c.storage {
			if c.storage[i].key == key {
				c.storage[i].value = value
				c.storage[i].shouldExpire = false
				break
			}
		}

	} else {
		c.storage = append(c.storage, Key{key: key, value: value, shouldExpire: false})
	}
}

func (c *Cache) Keys() []string {

	var validKeys []string
	for i := range c.storage {
		if !c.storage[i].shouldExpire {

			validKeys = append(validKeys, c.storage[i].key)

		} else if c.storage[i].expTime.After(time.Now()) {
			validKeys = append(validKeys, c.storage[i].key)
		}

	}
	return validKeys
}

func (c *Cache) PutTill(key, value string, deadline time.Time) {

	validKeys := c.Keys()
	needToOverwrite := false

	for _, validKey := range validKeys {
		if validKey == key {
			needToOverwrite = true
			break
		}
	}

	if needToOverwrite {
		for i := range c.storage {
			if c.storage[i].key == key {
				c.storage[i].value = value
				c.storage[i].shouldExpire = true
				c.storage[i].expTime = deadline
				break
			}
		}

	} else {
		c.storage = append(c.storage, Key{key: key, value: value, shouldExpire: false, expTime: deadline})
	}
}
