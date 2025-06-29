package operator

import (
	"bytes"
	"cacher/configs"
	"cacher/ramcache"
	"cacher/rules"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"

	"github.com/andybalholm/brotli"
)

func GetCachePath() string {
	return path.Join(configs.Config["save_path"], configs.Config["cache_file"])
}
func GetRulePath() string {
	return path.Join(configs.Config["save_path"], configs.Config["rule_file"])
}

func GetCacheItem(request *ramcache.RequestCache, rule rules.CompareRule) (*ramcache.CacheItem, bool) {
	cache_paths, ok := ramcache.RAM_CACHE[request.Url]
	if !ok {
		return nil, false
	}
	cache_method, ok := cache_paths[request.Method]
	if !ok {
		return nil, false
	}
	for _, cache := range cache_method {
		if MatchCompareRuleWithRequest(request, cache.Request, rule) {
			return &cache, true
		}
	}
	return nil, false
}

func RequestFromSourceUrl(host string, request *ramcache.RequestCache, rule rules.Rule) (*ramcache.ResponseCache, error) {
	if len(rule.RequestOperations) > 0 {
		err := OperateRequests(rule.RequestOperations, request)
		if err != nil {
			log.Println("operate request failed: ", err)
			return nil, err
		}
	}
	req_url := request.Url
	if rule.MatchRule.UrlPattern != "" {
		req_url = regexp.MustCompile(rule.MatchRule.UrlPattern).ReplaceAllString(request.Url, rule.MatchRule.UrlValue)
	}
	u, _ := url.Parse(host)
	u.Path = path.Join(u.Path, req_url)
	req, _ := http.NewRequest(request.Method, u.String(), bytes.NewReader(request.Body))
	// add header
	for key, value := range request.Header {
		for _, v := range value {
			req.Header.Add(key, v)
		}
	}
	// add cookies
	for key, value := range request.Cookies {
		req.AddCookie(&http.Cookie{
			Name:  key,
			Value: value,
		})
	}
	// act request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// decode brotli
	if res.Header.Get("content-encoding") == "br" {
		reader := brotli.NewReader(bytes.NewReader(body))
		body, _ = io.ReadAll(reader)
	}
	response_cache := ramcache.NewResponseCache(res.StatusCode, res.Header, body)
	if len(rule.ResponseOperations) > 0 {
		err = OperateResponses(rule.ResponseOperations, *response_cache)
		if err != nil {
			log.Println("operate response failed: ", err)
			return nil, err
		}
	}
	return response_cache, nil
}
