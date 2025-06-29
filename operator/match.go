package operator

import (
	"cacher/matcher"
	"cacher/ramcache"
	"cacher/rules"
)

func MatchMatchRule(request *ramcache.RequestCache, rule rules.MatchRule) bool {
	// check path
	if len(rule.Url) > 0 {
		if !matcher.MatchAnyOfStringList(rule.Url, request.Url, true) {
			return false
		}
	}
	// check method
	if len(rule.Method) > 0 {
		if !matcher.MatchAnyOfStringList(rule.Method, request.Method, false) {
			return false
		}
	}
	// check header
	if len(rule.Header) > 0 {
		if !matcher.MatchAllOfHeaderMatcherList(rule.Header, request.Header) {
			return false
		}
	}
	// check body
	if len((rule.Body)) > 0 {
		if !matcher.MatchAllOfJsonMatcherList(rule.Body, &request.Body) {
			return false
		}
	}
	return true
}

func MatchMatchRuleWithRequest(request *ramcache.RequestCache) (rules.Rule, bool) {
	for _, rule := range rules.RULES {
		if MatchMatchRule(request, rule.MatchRule) {
			return rule, true
		}
	}
	return rules.Rule{}, false
}
