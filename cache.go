package cache

import (
	"sync"
	"time"
)

type Item struct {
	value        string
	expTime      time.Time
	shouldExpire bool
}

// Cache with mutex to protect it from race conditions
type Cache struct {
	Storage map[string]Item
	sync.Mutex
}

//Cache constructor
func NewCache() *Cache {
	return &Cache{Storage: map[string]Item{}}
}

func (c *Cache) Get(key string) (string, bool) {
	c.Lock()
	defer c.Unlock()

	item, found := c.Storage[key]
	if !found {
		return "", false
	}

	//Check if it expired
	if item.shouldExpire && !item.expTime.After(time.Now()){
		return "",false
	}

	return item.value, true
}

func (c *Cache) Put(key string, value string) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Storage[key]; ok {
		c.Storage[key] = Item{shouldExpire: false, value: value}
	} else {
		c.Storage[key] = Item{
			value:        value,
			shouldExpire: false,
		}
	}

}

func (c *Cache) Keys() []string {
	c.Lock()
	defer c.Unlock()

	keys := []string{}
	for key := range c.Storage {
		if !c.Storage[key].shouldExpire {

			keys = append(keys, key)

		} else if c.Storage[key].expTime.After(time.Now()) {
			keys = append(keys, key)
		}
	}
	return keys
}

func (c *Cache) PutTill(key, value string, deadline time.Time) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Storage[key]; ok {
		c.Storage[key] = Item{value: value, shouldExpire: true, expTime: deadline}
	} else {
		c.Storage[key] = Item{
			value:        value,
			shouldExpire: true,
			expTime:      deadline,
		}
	}

}
