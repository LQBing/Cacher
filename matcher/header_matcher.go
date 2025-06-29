package matcher

import (
	"regexp"
	"strings"
)

type HeaderMatcher struct {
	Key     string   `json:"key"`
	Values  []string `json:"values,omitempty"`
	Regex   bool     `json:"regex,omitempty"`
	Exist   bool     `json:"exist,omitempty"`
	Exclude bool     `json:"exclude,omitempty"`
}

func NewHeaderMatcher(key string, values []string, regex bool, exist bool, exclude bool) *HeaderMatcher {
	return &HeaderMatcher{
		Key:     key,
		Values:  values,
		Regex:   regex,
		Exist:   exist,
		Exclude: exclude,
	}
}

func (m *HeaderMatcher) Match(header map[string][]string) bool {
	for k, l := range header {
		if strings.EqualFold(k, m.Key) {
			if m.Exist {
				return !m.Exclude
			} else {
				for _, v := range l {
					for _, value := range m.Values {
						if m.Regex {
							if match, _ := regexp.MatchString(value, v); match {
								return !m.Exclude
							}
						} else {
							if v == value {
								return !m.Exclude
							}
						}
					}
				}
			}
		}
	}
	return m.Exclude
}
func MatchAllOfHeaderMatcherList(matchers []HeaderMatcher, header map[string][]string) bool {
	for _, m := range matchers {
		if !m.Match(header) {
			return false
		}
	}
	return true
}
func MatchAnyOfHeaderMatcherList(matchers []HeaderMatcher, header map[string][]string) bool {
	for _, m := range matchers {
		if m.Match(header) {
			return true
		}
	}
	return false
}
