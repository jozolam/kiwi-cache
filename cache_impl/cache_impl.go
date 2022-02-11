package cache_impl

import (
	"fmt"
	"sync"
	"time"
)

const numberOfShards = 10

type Cacher interface {
	Get(id int) (string, error)
}

type Fetcher interface {
	Fetch(id int) (string, error)
	FetchAll() (map[int]string, error)
}

type cache struct {
	cache   [numberOfShards]cacheShard
	fetcher Fetcher
	name    string
}

type cacheShard struct {
	lock  *sync.RWMutex
	cache *map[int]string
}

func (c *cache) Get(id int) (string, error) {
	cs := c.cache[calculateShardIndex(id)]

	cs.lock.RLock()

	v, f := (*(cs.cache))[id]
	if f == false {
		v1, e := c.fetcher.Fetch(id)
		if e != nil {
			return "", e
		}
		cs.lock.RUnlock()
		cs.lock.Lock()
		(*(cs.cache))[id] = v1
		cs.lock.Unlock()
		return v1, nil
	}
	cs.lock.RUnlock()
	return v, nil
}

func InitCache(updatePeriod int, f Fetcher, name string) Cacher {

	c := cache{fetcher: f, name: name}
	for k := range c.cache {
		c.cache[k] = cacheShard{lock: new(sync.RWMutex), cache: new(map[int]string)}
	}
	ticker := time.NewTicker(time.Duration(updatePeriod) * time.Millisecond)

	err := c.updateAll()
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for {
			t := <-ticker.C
			err := c.updateAll()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(fmt.Sprintf("Updating %v cache at time %v", c.name, t))
		}
	}()

	return &c
}

func (c *cache) updateAll() error {
	var mappedResults [numberOfShards]map[int]string
	results, err := c.fetcher.FetchAll()
	if err != nil {
		return err
	}

	for i := range mappedResults {
		mappedResults[i] = make(map[int]string)
	}

	for k, v := range results {
		mappedResults[calculateShardIndex(k)][k] = v
	}

	for k, v := range mappedResults {
		c.cache[k].lock.Lock()
		*(c.cache[k].cache) = make(map[int]string)
		for k1, v1 := range v {
			(*(c.cache[k].cache))[k1] = v1
		}
		c.cache[k].lock.Unlock()
	}
	return nil
}

func calculateShardIndex(id int) int {
	return id % numberOfShards
}
