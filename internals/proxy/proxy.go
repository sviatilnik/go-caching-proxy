package proxy

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
)

type Proxy struct {
	pattern string
	re      *regexp.Regexp
	proxy   *httputil.ReverseProxy
}

func NewProxy(pattern string, target string) (*Proxy, error) {
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
	}, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		if match := p.re.MatchString(r.URL.Path); match {
			// TODO cache function here
			slog.Info("Get request from cache")
			return
		}
	}

	slog.Info("Proxy request")
	p.proxy.ServeHTTP(w, r)
}
