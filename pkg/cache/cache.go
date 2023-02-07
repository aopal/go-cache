package cache

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/aopal/go-cache/pkg/config"
	"github.com/aopal/go-cache/pkg/fetch"
)

type Cache struct {
	cfg *config.Config

	storage sync.Map // map[string]*CacheEntry
}

type cacheEntry struct {
	val         *fetch.Response
	varies      bool
	varyHeaders []string
	children    sync.Map // map[string]*fetch.Response
}

func New(cfg *config.Config) *Cache {
	return &Cache{
		cfg:       cfg,
	}
}

// Request is also passed for vary support
func (c *Cache) Lookup(cacheKey string, req *http.Request) (*fetch.Response, error) {
	e, inCache := c.storage.Load(cacheKey)
	if !inCache {
		return nil, nil
	}

	entry := e.(*cacheEntry)
	if !entry.varies {
		return entry.val, nil
	} else {
		varyKey := c.getVaryString(req, entry.varyHeaders)

		c, loaded := entry.children.Load(varyKey)
		if loaded {
			return c.(*fetch.Response), nil
		}
		return nil, nil
	}
}

// inserts without overwriting
func (c *Cache) Insert(cacheKey string, req *http.Request, resp *fetch.Response) {
	vary := resp.Header.Get("Vary")
	varyHeaders := c.getVaryHeaders(vary)
	varyKey := c.getVaryString(req, varyHeaders)

	e, _ := c.storage.LoadOrStore(cacheKey, &cacheEntry{
		varies:      vary != "",
		varyHeaders: varyHeaders,
		val:         resp,
	})

	entry := e.(*cacheEntry)
	if entry.varies {
		entry.children.LoadOrStore(varyKey, resp)
	}
}

func (c *Cache) getVaryHeaders(varyString string) []string {
	return strings.Split(strings.ReplaceAll(varyString, " ", ""), ",")
}

func (c *Cache) getVaryString(req *http.Request, varyHeaders []string) string {
	varyKey := ""

	for _, name := range varyHeaders {
		varyKey = fmt.Sprintf("%s=%s;%s", name, req.Header.Get(name), varyKey)
	}

	return varyKey
}
