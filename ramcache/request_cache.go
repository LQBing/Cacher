package ramcache

import (
	"net/http"
)

type RequestCache struct {
	Url     string            `json:"path"`
	Method  string            `json:"method"`
	Header  http.Header       `json:"headers"`
	Cookies map[string]string `json:"cookies"`
	Body    []byte            `json:"body"`
}

func NewRequestCache(url string, method string, header http.Header, cookies map[string]string, body []byte) *RequestCache {
	return &RequestCache{
		Url:     url,
		Method:  method,
		Header:  header,
		Cookies: cookies,
		Body:    body,
	}
}
