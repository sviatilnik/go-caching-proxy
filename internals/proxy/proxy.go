package proxy

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/sviatilnik/go-caching-proxy/internals/cache"
	"golang.org/x/sync/singleflight"
)

type Proxy struct {
	pattern      string
	re           *regexp.Regexp
	proxy        *httputil.ReverseProxy
	cache        *cache.Cache
	target       *url.URL
	singleflight singleflight.Group
}

func NewProxy(pattern string, target string, cache *cache.Cache) (*Proxy, error) {
	parseURL, err := url.Parse(target)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	prx := httputil.NewSingleHostReverseProxy(parseURL)
	prx.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		slog.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return &Proxy{
		pattern: pattern,
		proxy:   prx,
		re:      re,
		cache:   cache,
		target:  parseURL,
	}, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet || p.cache == nil {
		slog.Info("Proxy request")
		p.proxy.ServeHTTP(w, r)
		return
	}

	if match := p.re.MatchString(r.URL.Path); match {

		var err error
		var resp *cache.Response

		if p.cache.Has(r) {
			resp, err = p.cache.Get(r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			if resp != nil {
				slog.Info("Get response from cache")

				resp.AddCacheHitHeader()
			}
		}

		if resp == nil {
			path := p.getProxyPath(r)
			v, err, _ := p.singleflight.Do(path, func() (interface{}, error) {
				return p.sendRequest(r)
			})

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			resp = v.(*cache.Response)

			resp.AddCacheMissHeader()
		}

		resp.CopyToWriter(w)

		return
	}

}

func (p *Proxy) sendRequest(r *http.Request) (*cache.Response, error) {
	// TODO Вынести Timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	path := p.getProxyPath(r)

	clientReq, err := http.NewRequest(r.Method, path, r.Body)
	if err != nil {
		return nil, err
	}
	clientReq.Header = r.Header

	httpResponse, err := client.Do(clientReq)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	var headers []cache.Header
	for key, head := range httpResponse.Header {
		for _, v := range head {
			headers = append(headers, cache.Header{
				Name:  key,
				Value: v,
			})
		}
	}

	resp := &cache.Response{
		Status:  httpResponse.StatusCode,
		Body:    string(body),
		Headers: headers,
	}

	// Статус код 200 семейства, поэтому пишем ответ в кеш
	if httpResponse.StatusCode < 300 {
		if err := p.cache.Save(r, resp); err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (p *Proxy) getProxyPath(r *http.Request) string {
	return strings.TrimRight(p.target.String(), "/") + p.re.ReplaceAllString(r.URL.Path, "") + r.URL.RawQuery
}
