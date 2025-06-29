package matcher

import (
	"reflect"
	"regexp"

	"github.com/jefferyjob/go-easy-utils/anyUtil"
	"github.com/ohler55/ojg/jp"
)

type JsonMatcher struct {
	Path      string `json:"path,omitempty"`
	Pattern   any    `json:"pattern,omitempty"`
	Regex     bool   `json:"regex,omitempty"`
	FullMatch bool   `json:"fullMatch,omitempty"`
	Exist     bool   `json:"exist,omitempty"`
	Exclude   bool   `json:"exclude,omitempty"`
}

func NewJsonMatcher(path string, pattern any, regex bool, fullMatch bool, exist bool, exclude bool) *JsonMatcher {
	return &JsonMatcher{
		Path:      path,
		Pattern:   pattern,
		Regex:     regex,
		FullMatch: fullMatch,
		Exist:     exist,
		Exclude:   exclude,
	}
}
func (m *JsonMatcher) Match(value any) bool {
	if m.Path == "" {
		m.Path = "$"
	}
	x, _ := jp.ParseString(m.Path)
	results := x.Get(value)
	if len(results) == 0 {
		return m.Exclude
	} else {
		if m.Exist {
			return !m.Exclude
		} else {
			for _, result := range results {
				if m.FullMatch {
					match := reflect.DeepEqual(m.Pattern, result)
					if match {
						return !m.Exclude
					} else {
						return m.Exclude
					}
				}
				if m.Regex {
					match, _ := regexp.MatchString(anyUtil.AnyToStr(m.Pattern), anyUtil.AnyToStr(result))
					if match {
						return !m.Exclude
					} else {
						return m.Exclude
					}
				}
			}
		}
	}
	return m.Exclude
}

func MatchAllOfJsonMatcherList(matcher_list []JsonMatcher, value any) bool {
	for _, matcher := range matcher_list {
		if !matcher.Match(value) {
			return false
		}
	}
	return true
}
func MatchAnyOfJsonMatcherList(matcher_list []JsonMatcher, value any) bool {
	for _, matcher := range matcher_list {
		if matcher.Match(value) {
			return true
		}
	}
	return false
}
