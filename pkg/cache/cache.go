package cache

import (
	"container/list"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/aopal/go-cache/pkg/config"
	"github.com/aopal/go-cache/pkg/fetch"
)

type Cache struct {
	cfg *config.Config

	storage sync.Map // map[string]*cacheEntry
	tags    sync.Map // map[string]*tagList
}

type cacheEntry struct {
	key         string
	val         *fetch.Response
	varies      bool
	varyHeaders []string
	expiry      time.Time
	purged      bool
}

type tagList struct {
	tags *list.List // linked list of *cacheEntry
	sync.Mutex
}

func New(cfg *config.Config) *Cache {
	c := &Cache{
		cfg: cfg,
	}

	c.sweep()

	return c
}

// Request is also passed for vary support
func (c *Cache) Lookup(cacheKey string, req *http.Request) (*fetch.Response, error) {
	e, inCache := c.storage.Load(cacheKey)
	if !inCache {
		return nil, nil
	}

	var entry *cacheEntry = e.(*cacheEntry)
	if entry.varies {
		varyKey := c.getVaryString(req, entry.varyHeaders)

		c, loaded := c.storage.Load(varyKey + cacheKey)
		if loaded {
			entry = c.(*cacheEntry)
		} else {
			return nil, nil
		}
	}

	return entry.val, nil
}

// inserts without overwriting
func (c *Cache) Insert(cacheKey string, req *http.Request, resp *fetch.Response) {
	shouldCache, expiry := c.parseCacheControl(resp)
	if !shouldCache || expiry == unixEpoch {
		return
	}

	vary := resp.Header.Get("Vary")
	varyHeaders := split(vary)
	varyKey := c.getVaryString(req, varyHeaders)

	e, _ := c.storage.LoadOrStore(cacheKey, &cacheEntry{
		key:         cacheKey,
		varies:      vary != "",
		varyHeaders: varyHeaders,
		val:         resp,
		expiry:      expiry,
	})

	var entry *cacheEntry = e.(*cacheEntry)
	if entry.varies {
		e, _ = c.storage.LoadOrStore(varyKey+cacheKey, &cacheEntry{
			key:    varyKey + cacheKey,
			varies: false,
			val:    resp,
			expiry: expiry,
		})
		entry = e.(*cacheEntry)
	}

	c.tagEntry(entry, resp)
}

func (c *Cache) tagEntry(entry *cacheEntry, resp *fetch.Response) {
	tags := split(resp.Header.Get("Cache-Tag"))

	for _, tag := range tags {
		l, _ := c.tags.LoadOrStore(tag, &tagList{
			tags: list.New(),
		})
		var tagList *tagList = l.(*tagList)

		tagList.Lock()
		tagList.tags.PushBack(entry)
		tagList.Unlock()
	}
}

func split(str string) []string {
	return strings.Split(strings.ReplaceAll(str, " ", ""), ",")
}

func (c *Cache) getVaryString(req *http.Request, varyHeaders []string) string {
	varyKey := ""

	for _, name := range varyHeaders {
		varyKey = fmt.Sprintf("%s=%s;%s", name, req.Header.Get(name), varyKey)
	}

	return varyKey
}
