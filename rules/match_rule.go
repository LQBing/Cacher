package rules

import "cacher/matcher"

// match request, if not match, just proxy request.
type MatchRule struct {
	Url        []string                   `json:"url,omitempty"`
	Method     []string                   `json:"method,omitempty"` //
	Header     []matcher.HeaderMatcher    `json:"headers,omitempty"`
	Cookie     []matcher.StringMapMatcher `json:"cookies,omitempty"`
	Body       []matcher.JsonMatcher      `json:"body,omitempty"`
	SourceUrl  string                     `json:"sourceUrl,omitempty"`
	UrlPattern string                     `json:"urlPattern,omitempty"` // regex
	UrlValue   string                     `json:"urlValue,omitempty"`   // regex
}

func NewMatchRule(url []string, method []string, header []matcher.HeaderMatcher, cookie []matcher.StringMapMatcher, body []matcher.JsonMatcher, source_url string, url_pattern string, url_value string) *MatchRule {
	return &MatchRule{
		Url:        url,
		Method:     method,
		Header:     header,
		Cookie:     cookie,
		Body:       body,
		SourceUrl:  source_url,
		UrlPattern: url_pattern,
		UrlValue:   url_value,
	}
}
