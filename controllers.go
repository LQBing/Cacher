package main

import (
	"cacher/configs"
	"cacher/operator"
	"cacher/ramcache"
	"cacher/rules"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func cacheHandler(c *gin.Context) {
	if c.Request.URL.Path == "/health" {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	}
	// ineternal apis
	if c.Request.URL.Path == configs.Config["API_PATH"]+"/rules/get" {
		c.JSON(http.StatusOK, rules.RULES)
		return
	}
	if c.Request.URL.Path == configs.Config["API_PATH"]+"/caches/get" {
		path := c.Query("path")
		method := strings.ToUpper(c.Query("method"))
		if path == "" {
			c.JSON(http.StatusOK, ramcache.RAM_CACHE)
			return
		} else {
			if _, ok := ramcache.RAM_CACHE[path]; ok {
				if method == "" || method == "ALL" {
					c.JSON(http.StatusOK, ramcache.RAM_CACHE[path])
					return
				} else {
					if _, ok := ramcache.RAM_CACHE[path][method]; ok {
						c.JSON(http.StatusOK, ramcache.RAM_CACHE[path][method])
						return
					}
				}
			}
		}
		c.JSON(http.StatusOK, "")
		return
	}
	if c.Request.URL.Path == configs.Config["API_PATH"]+"/caches/list" {
		path := c.Query("path")
		if path == "" {
			var paths []string
			for p := range ramcache.RAM_CACHE {
				paths = append(paths, p)
			}
			c.JSON(http.StatusOK, paths)
			return
		} else {
			if _, ok := ramcache.RAM_CACHE[path]; ok {
				var methods []string
				for m := range ramcache.RAM_CACHE[path] {
					methods = append(methods, m)
				}
				c.JSON(http.StatusOK, methods)
				return
			}
		}
		c.JSON(http.StatusOK, "")
		return
	}
	if c.Request.URL.Path == configs.Config["API_PATH"]+"/caches/del" {
		path := c.Query("path")
		method := strings.ToUpper(c.Query("method"))
		if _, ok1 := ramcache.RAM_CACHE[path]; ok1 {
			if method == "" {
				delete(ramcache.RAM_CACHE, path)
				ramcache.Save(operator.GetCachePath())
			}
			if _, ok2 := ramcache.RAM_CACHE[path][method]; ok2 {
				delete(ramcache.RAM_CACHE[path], path)
				ramcache.Save(operator.GetCachePath())
			}
		}
		c.JSON(http.StatusBadRequest, "not fond caches")
		return
	}
	// create request cache
	req_cookies := make(map[string]string)
	for _, cookie := range c.Request.Cookies() {
		req_cookies[cookie.Name] = cookie.Value
	}
	req_body, _ := c.GetRawData()
	rc := ramcache.NewRequestCache(c.Request.URL.String(), c.Request.Method, c.Request.Header, req_cookies, req_body)
	// match rules
	rule, ok := operator.MatchMatchRuleWithRequest(rc)
	if !ok && rule.MatchRule.SourceUrl == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 Not Found!"})
		return
	}
	// cache hit
	cache_item, ok := operator.GetCacheItem(rc, rule.CompareRule)
	if ok {
		c.Data(http.StatusOK, cache_item.Response.Header.Get("Content-Type"), []byte(cache_item.Response.Body))
		// replace header
		for k := range c.Writer.Header() {
			delete(c.Writer.Header(), k)
		}
		for k, v := range cache_item.Response.Header {
			for i, vv := range v {
				if i == 0 {
					c.Writer.Header().Set(k, vv)
				} else {
					c.Writer.Header().Add(k, vv)
				}
			}
		}
		if configs.Config["debug"] != "" {
			log.Println("hit " + rc.Url)
		}
		if configs.Config["debug_body"] != "" {
			log.Println("req body: " + string(rc.Body))
			log.Println("resp body: " + string(cache_item.Response.Body))
		}
		return
	}
	// no cache, request and cache
	res, err := operator.RequestFromSourceUrl(rule.MatchRule.SourceUrl, rc, rule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if res != nil {
		if res.StatusCode == http.StatusOK {
			// save cache item
			cache_item := ramcache.CacheItem{
				Request:  *rc,
				Response: *res}

			ramcache.AddCacheItem(rc.Url, rc.Method, cache_item)
			ramcache.Save(operator.GetCachePath())
		}
		// return response
		cache_item, _ := operator.GetCacheItem(rc, rule.CompareRule)
		c.Data(http.StatusOK, cache_item.Response.Header.Get("Content-Type"), []byte(cache_item.Response.Body))
		// replace header
		for k := range c.Writer.Header() {
			delete(c.Writer.Header(), k)
		}
		for k, v := range cache_item.Response.Header {
			for i, vv := range v {
				if i == 0 {
					c.Writer.Header().Set(k, vv)
				} else {
					c.Writer.Header().Add(k, vv)
				}
			}
		}

		if configs.Config["debug"] != "" {
			log.Println("call source url " + rc.Url)
		}
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "404 Not Found "})
}
