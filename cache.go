package cache

import "time"

type Cache struct {
	Storage []Key
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

	for _, item := range c.Storage {
		if item.key == key && (item.expirationTime.Before(time.Now()) || !item.willExpired) {
			return item.value, true
		}
	}

	return "", false
}

func (c *Cache) Put(key, value string) {

	needOverwrite := false
	rangeM := c.Keys()
	for _, nonExpKey := range rangeM {
		if nonExpKey == key {
			needOverwrite = true
		}

	}

	for i := range c.Storage {

		if c.Storage[i].key == key {
			c.Storage[i].value = value
			c.Storage[i].willExpired = false

		}
	}

	if !needOverwrite {
		c.Storage = append(c.Storage, Key{key: key, value: value, willExpired: false})

	}

}

func (c *Cache) Keys() []string {
	var keys []string

	for _, item := range c.Storage {
		if item.expirationTime.Before(time.Now()) && !item.willExpired {
			keys = append(keys, item.key)
		}
	}

	return keys
}

func (c *Cache) PutTill(key, value string, deadline time.Time) {

	needOverwrite := false
	for _, nonExpKey := range c.Keys() {
		if nonExpKey == key {
			needOverwrite = true
		}

	}

	if needOverwrite {
		for _, item := range c.Storage {
			if item.key == key {
				item.value = value
				item.willExpired = true
				item.expirationTime = deadline
			}
		}

	} else {
		c.Storage = append(c.Storage, Key{key: key, value: value, willExpired: true, expirationTime: deadline})

	}
}
