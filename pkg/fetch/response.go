package fetch

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type Response struct {
	Status int
	Header http.Header
	Body   []byte

	// beReq  *http.Request
	// client *http.Client
	wg sync.WaitGroup
	// sync.Mutex
}

func NewResponse() *Response {
	r := &Response{}
	r.wg.Add(1)

	return r
}

func (resp *Response) Fetch(client *http.Client, beReq *http.Request) {
	log.Printf("Fetching from origin: %s\n", beReq.URL.String())
	beResp, err := client.Do(beReq)
	if err != nil {
		log.Printf("error while making origin request: %+v", err)
		// do something
	}

	body, err := ioutil.ReadAll(beResp.Body)
	if err != nil {
		log.Printf("error while reading origin response body: %+v", err)
	}

	resp.Status = beResp.StatusCode
	resp.Header = beResp.Header.Clone()
	resp.Body = body

	resp.wg.Done()
}

func (resp *Response) Clone() *Response {
	return &Response{
		Status: resp.Status,
		Body:   resp.Body,
		Header: resp.Header.Clone(),
	}
}
