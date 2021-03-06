package gocache

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

// constants for default cache life
const (
	Eternal = -1
)

// ***** Models ****

// cached file content
type Item struct {
	Content interface{}
	Life    int64
}

// garbage collector
type garbageCollector struct {
	interval time.Duration
	stop     chan bool
}

// Cache properties
type Cache struct {
	DefaultLife      time.Duration
	pool             *workerPool
	mutex            sync.RWMutex
	items            map[string]Item
	garbageCollector *garbageCollector
}

//***** Garbage Collector Functions *****

// Garbage Collector
func (j *garbageCollector) Run(c *Cache) {
	ticker := time.NewTicker(j.interval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

// Stop garbage collector
func stopGC(c *Cache) {
	c.garbageCollector.stop <- true
}

func runGC(c *Cache, ci time.Duration) {
	gc := &garbageCollector{
		interval: ci,
		stop:     make(chan bool),
	}
	c.garbageCollector = gc

	go gc.Run(c)
}

func (c *Cache) DeleteBackground(key string) {
	if c.pool == nil {
		log.Println("Pool is empty")
		return
	}
	c.pool.Submit(&deleteJob{job: job{cache: c, key: key}})
}

func (c *Cache) DeleteCachedItem(key string) {
	c.mutex.Lock()
	delete(c.items, key)
	c.mutex.Unlock()
}

func (c *Cache) ClearCache() {
	c.mutex.Lock()

	for key := range c.items {
		delete(c.items, key)
	}

	c.mutex.Unlock()
}

func (c *Cache) DeleteExpired() {
	now := time.Now().UnixNano()

	c.mutex.Lock()

	for key, item := range c.items {
		// Delete items excepts eternal one(v.Life=-1)
		if item.Life > 0 && now > item.Life {
			delete(c.items, key)
		}
	}

	c.mutex.Unlock()
}

//***** Cache Functions *****
func (c *Cache) set(key string, content interface{}, life time.Duration) {
	var lifeTime int64

	if life > 0 {
		lifeTime = time.Now().Add(life).UnixNano()
	} else {
		lifeTime = Eternal
	}

	c.items[key] = Item{
		Content: content,
		Life:    lifeTime,
	}
}

func (c *Cache) Set(key string, content interface{}, life time.Duration) {
	var lifeTime int64

	if life > 0 {
		lifeTime = time.Now().Add(life).UnixNano()
	} else {
		lifeTime = Eternal
	}

	c.mutex.Lock()

	c.items[key] = Item{
		Content: content,
		Life:    lifeTime,
	}

	c.mutex.Unlock()
}

func (c *Cache) get(key string) (interface{}, bool) {
	item, found := c.items[key]

	if !found {
		return nil, false
	}

	return item.Content, true
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()

	item, found := c.items[key]

	if !found {
		c.mutex.RUnlock()
		return nil, false
	}

	c.mutex.RUnlock()

	return item.Content, true
}

func (c *Cache) Add(key string, content interface{}, life time.Duration) (bool, string) {
	c.mutex.Lock()

	_, found := c.get(key)

	if found {
		c.mutex.Unlock()
		return false, fmt.Sprintf("Item with key: %s already exists", key)
	}

	c.set(key, content, life)

	c.mutex.Unlock()

	return true, "OK"
}

func (c *Cache) Print() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	fmt.Printf("Length: %d\n", len(c.items))
	for k, value := range c.items {
		fmt.Printf("key: %s value: %v\n", k, value)
	}
}

func (c *Cache) AddBackground(key string, content interface{}, life time.Duration) {
	if c.pool == nil {
		log.Println("Pool is empty")
		return
	}
	c.pool.Submit(&addingJob{job: job{cache: c, key: key, content: content, life: life}})
}

func (c *Cache) StartWorkerPoolWith(workerCount int) {
	if c.pool != nil {
		return
	}
	c.pool = newWorkerPool(workerCount)
	c.pool.Start()
}

func (c *Cache) StopWorkerPool() {
	c.pool.Stop()
}

func NewCache(defaultLife, gbcInterval time.Duration) *Cache {
	items := make(map[string]Item)
	c := &Cache{
		DefaultLife: defaultLife,
		items:       items,
	}

	if gbcInterval > 0 {
		runGC(c, gbcInterval)
		runtime.SetFinalizer(c, stopGC)
	}

	return c
}
