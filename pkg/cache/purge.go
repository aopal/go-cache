package cache

import (
	"log"
	"runtime"
	"time"
)

type PurgeMethod int

const (
	Purge PurgeMethod = iota
	Invalidate
)

func (c *Cache) PurgeURL(url string) {

}

func (c *Cache) PurgeByTag(tag string, method PurgeMethod) {
	l, loaded := c.tags.Load(tag)
	if !loaded {
		return
	}
	var tagList *tagList = l.(*tagList)

	tagList.Lock()
	defer tagList.Unlock()

	for e := tagList.tags.Front(); e != nil; e = e.Next() {
		var entry *cacheEntry = e.Value.(*cacheEntry)

		switch method {
		case Purge:
			// c.storage.Delete(entry.key)
			c.purgeNow(entry)
		case Invalidate:
			entry.purged = true
		}

		// e = e.Next()
		// prev := e.Prev()
		// if prev != nil {
		// 	tagList.tags.Remove(prev)
		// }

	}

	tagList.tags.Init() // clear taglist to remove references to *cacheEntries
	runtime.GC()
}

func (c *Cache) sweep() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			log.Println("Sweeping cache...")

			count := 0
			c.storage.Range(func(key, value any) bool {
				if c.purgeIfExpired(key, value) {
					count++
				}
				return true
			})

			runtime.GC()
			log.Printf("Purged %d entries\n", count)
		}
	}()
}

func (c *Cache) purgeIfExpired(key, value any) bool {
	var entry *cacheEntry = value.(*cacheEntry)
	if time.Now().After(entry.expiry) || entry.purged {
		// c.storage.Delete(entry.key)
		c.purgeNow(entry)
		return true
	}
	return false
}

func (c *Cache) purgeNow(entry *cacheEntry) {
	entry.val.Body = nil
	c.storage.Delete(entry.key)
}
