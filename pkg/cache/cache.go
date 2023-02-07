package cache

import (
	"net/http"

	"github.com/aopal/go-cache/pkg/config"
	"github.com/aopal/go-cache/pkg/fetch"
)

type Cache struct {
	cfg *config.Config

	storagev1 map[string]*fetch.Response
}

func New(cfg *config.Config) *Cache {
	return &Cache{
		cfg:       cfg,
		storagev1: make(map[string]*fetch.Response),
	}
}

// Request is also passed for vary support
func (c *Cache) Lookup(cacheKey string, _ *http.Request) (*fetch.Response, error) {
	if resp, ok := c.storagev1[cacheKey]; ok {
		return resp, nil
	}

	return nil, nil
}

func (c *Cache) Insert(cacheKey string, req *http.Request, resp *fetch.Response) {
	c.storagev1[cacheKey] = resp
}
