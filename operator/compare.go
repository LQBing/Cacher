package operator

import (
	"cacher/comparator"
	"cacher/ramcache"
	"cacher/rules"
)

func MatchCompareRuleWithRequest(request *ramcache.RequestCache, cache ramcache.RequestCache, rule rules.CompareRule) bool {
	// match method
	if request.Method != cache.Method {
		return false
	}
	// match headers
	if !comparator.CompareHeadersWithComparator(request.Header, cache.Header, rule.Headers) {
		return false
	}
	// match cookies
	if !comparator.CompareCookiesWithComparator(request.Cookies, cache.Cookies, rule.Cookies) {
		return false
	}
	// match body
	if !comparator.CompareBodyWithComparator(request.Body, cache.Body, rule.Body) {
		return false
	}
	// match status code
	return true
}
