package ramcache

type CacheItem struct {
	Request  RequestCache  `json:"req,omitempty"`
	Response ResponseCache `json:"res,omitempty"`
}

func NewCacheItem(request RequestCache, response ResponseCache) *CacheItem {
	return &CacheItem{
		Request:  request,
		Response: response,
	}
}
