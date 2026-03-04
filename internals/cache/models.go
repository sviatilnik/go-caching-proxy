package cache

import (
	"net/http"
)

type Response struct {
	Body    string   `json:"body"`
	Status  int      `json:"status"`
	Headers []Header `json:"headers"`
}

func (r *Response) CopyToWriter(w http.ResponseWriter) {
	for _, h := range r.Headers {
		w.Header().Add(h.Name, h.Value)
	}

	w.WriteHeader(r.Status)
	w.Write([]byte(r.Body))
}

func (r *Response) AddCacheMissHeader() {
	r.Headers = append(r.Headers, Header{
		Name:  "X-Cache",
		Value: "MISS",
	})
}

func (r *Response) AddCacheHitHeader() {
	r.Headers = append(r.Headers, Header{
		Name:  "X-Cache",
		Value: "HIT",
	})
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
