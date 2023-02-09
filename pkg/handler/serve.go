package handler

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/aopal/go-cache/pkg/fetch"
)

func (h *Handler) Serve(w http.ResponseWriter, r *http.Request) {
	h.normalizeRequest(r)
	cacheKey := h.getCacheKey(r)

	response, err := h.cache.Lookup(cacheKey, r)
	if err != nil {
		log.Printf("error while looking up in cache: %+v", err)
	}

	if response != nil {
		w.Header().Add("X-Cache-Status", "HIT")
		_, err = h.streamResponse(response, w)
		if err != nil {
			log.Printf("error while streaming cached response: %+v", err)
		}
		return
	}

	response, err = h.fetcher.Fetch(r, cacheKey)
	if err != nil {
		log.Printf("error from fetching: v", err)
	}

	w.Header().Add("X-Cache-Status", "MISS")
	h.normalizeResponse(response)
	_, err = h.streamResponse(response, w)
	if err != nil {
		log.Printf("error from streaming response: %+v", err)
	}

	go func() {
		h.cache.Insert(cacheKey, r, response)
		h.fetcher.RemoveInProgress(cacheKey)
	}()
}

func (h *Handler) streamResponse(resp *fetch.Response, w http.ResponseWriter) (int64, error) {
	for name, values := range resp.Header {
		for _, val := range values {
			w.Header().Add(name, val)
		}
	}

	w.WriteHeader(resp.Status)

	return io.Copy(w, bytes.NewReader(resp.Body))
}

// do things like normalizing accept headers
func (h *Handler) normalizeRequest(r *http.Request) {
	accept := r.Header.Get("Accept")

	if strings.Contains(accept, "image/jxl") {
		r.Header.Set("Accept", "image/jxl,image/avif,image/webp")
	} else if strings.Contains(accept, "image/avif") {
		r.Header.Set("Accept", "image/avif,image/webp")
	} else if strings.Contains(accept, "image/webp") {
		r.Header.Set("Accept", "image/webp")
	} else if strings.Contains(accept, "image/") {
		r.Header.Set("Accept", "image/*")
	}
}

func (h *Handler) normalizeResponse(resp *fetch.Response) {
	// resp.Header.Del("Cache-Tag")
	// resp.Header.Del("Surrogate-Key")
}

// generate a cache key from the request,
func (h *Handler) getCacheKey(req *http.Request) string {
	return req.URL.String()
}
