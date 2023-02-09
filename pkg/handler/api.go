package handler

import (
	"fmt"
	"net/http"

	"github.com/aopal/go-cache/pkg/cache"
)

func (h *Handler) PurgeByTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Purge requests must be POSTs"))
		return
	}

	tag := r.URL.Query().Get("tag")
	if tag == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Tag not provided"))
		return
	}

	purgeMethod := cache.Purge
	method := r.URL.Query().Get("method")
	if method == "" || method == "purge" {
		purgeMethod = cache.Purge
	} else if method == "invalidate" {
		purgeMethod = cache.Invalidate
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unrecognized purge method '%s'", method)))
		return
	}

	h.cache.PurgeByTag(tag, purgeMethod)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Purge complete"))
}
