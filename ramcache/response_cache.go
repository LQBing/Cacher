package ramcache

import "net/http"

type ResponseCache struct {
	StatusCode int         `json:"code"`
	Header     http.Header `json:"header"`
	Body       []byte      `json:"body"`
}

func NewResponseCache(code int, header http.Header, body []byte) *ResponseCache {
	return &ResponseCache{
		StatusCode: code,
		Header:     header,
		Body:       body,
	}
}
