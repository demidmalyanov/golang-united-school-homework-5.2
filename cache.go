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

	for i := range c.Storage {
		if c.Storage[i].key == key && (c.Storage[i].expirationTime.Before(time.Now()) || !c.Storage[i].willExpired) {
			return c.Storage[i].value, true
		}
	}

	return "", false
}

func (c *Cache) Put(key, value string) {

	needOverwrite := false
	for _, nonExpKey := range c.Keys() {
		if nonExpKey == key {
			needOverwrite = true
		}

	}

	if needOverwrite {
		for i := range c.Storage {

			if c.Storage[i].key == key {
				c.Storage[i].value = value
				c.Storage[i].willExpired = false

			}
		}

	} else {
		c.Storage = append(c.Storage, Key{key: key, value: value, willExpired: false})

	}

}

func (c *Cache) Keys() []string {
	var keys []string

	for i := range c.Storage {
		if c.Storage[i].expirationTime.Before(time.Now()) && !c.Storage[i].willExpired {
			keys = append(keys, c.Storage[i].key)
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
		for i := range c.Storage {
			if c.Storage[i].key == key {
				c.Storage[i].value = value
				c.Storage[i].willExpired = true
				c.Storage[i].expirationTime = deadline
			}
		}

	} else {
		c.Storage = append(c.Storage, Key{key: key, value: value, willExpired: true, expirationTime: deadline})

	}
}
