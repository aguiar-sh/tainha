package proxy

import (
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(target string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	return proxy, nil
}
