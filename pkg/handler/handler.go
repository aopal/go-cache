package handler

import (
	"github.com/aopal/go-cache/pkg/cache"
	"github.com/aopal/go-cache/pkg/config"
	"github.com/aopal/go-cache/pkg/fetch"
)

type Handler struct {
	cfg *config.Config

	cache   *cache.Cache
	fetcher *fetch.Fetcher
}

func New(cfg *config.Config) (*Handler, error) {
	h := &Handler{
		cfg:     cfg,
		cache:   cache.New(cfg),
		fetcher: fetch.New(cfg),
	}

	return h, nil
}
