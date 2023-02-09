package cache

import (
	"log"
	"time"

	"github.com/aopal/go-cache/pkg/fetch"
	"github.com/pquerna/cachecontrol/cacheobject"
)

var (
	unixEpoch = time.Time{}
)

func (c *Cache) parseCacheControl(resp *fetch.Response) (bool, time.Time) {
	directives, err := cacheobject.ParseResponseCacheControl(resp.Header.Get("Cache-Control"))
	if err != nil {
		log.Printf("error while parsing cache control header: %+v", err)
		return false, unixEpoch
	}

	shouldCache := false
	expiry := unixEpoch

	// don't cache
	if directives.NoCachePresent || directives.NoStore || directives.PrivatePresent {
		return shouldCache, expiry
	}

	if directives.SMaxAge != -1 {
		shouldCache = true
		expiry = time.Now().Add(time.Duration(directives.SMaxAge) * time.Second)
	} else if directives.MaxAge != -1 {
		shouldCache = true
		expiry = time.Now().Add(time.Duration(directives.MaxAge) * time.Second)
	}

	return shouldCache, expiry
}
