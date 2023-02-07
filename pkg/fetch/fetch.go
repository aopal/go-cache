package fetch

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/aopal/go-cache/pkg/config"
)

type Fetcher struct {
	cfg    *config.Config
	client *http.Client
	// inProgress map[string]*Response
	inProgress sync.Map
}

func New(cfg *config.Config) *Fetcher {
	return &Fetcher{
		cfg:    cfg,
		client: http.DefaultClient,
	}
}

// func (f *Fetcher) FetchOld(req *http.Request) (*http.Response, error) {
// 	return f.client.Do(f.buildBackendRequest(req))
// }

func (f *Fetcher) Fetch(req *http.Request, cacheKey string) (*Response, error) {
	resp := NewResponse()

	respAny, loaded := f.inProgress.LoadOrStore(cacheKey, resp)
	resp = respAny.(*Response)
	if !loaded {
		resp.Fetch(f.client, f.buildBackendRequest(req))
	}

	resp.wg.Wait()

	return resp, nil
}

func (f *Fetcher) RemoveInProgress(cacheKey string) {
	f.inProgress.Delete(cacheKey)
}

func (f *Fetcher) buildBackendRequest(req *http.Request) *http.Request {
	beURL := fmt.Sprintf("%s%s", f.cfg.Origins[0], req.URL.RequestURI())

	beReq, _ := http.NewRequest(req.Method, beURL, nil)
	beReq.Header = req.Header.Clone()

	return beReq
}
