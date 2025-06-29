package matcher

import "regexp"

type StringListMatcher struct {
	Pattern string `json:"pattern,omitempty"`
	Regex   bool   `json:"regex,omitempty"`
	Exclude bool   `json:"exclude,omitempty"`
}

func NewStringListMatcher(pattern string, regex bool) *StringListMatcher {
	return &StringListMatcher{
		Pattern: pattern,
		Regex:   regex,
	}
}

func (m *StringListMatcher) Match(value []string) bool {
	for _, v := range value {
		if m.Regex {
			match, _ := regexp.MatchString(m.Pattern, v)
			if match {
				return !m.Exclude
			} else {
				return m.Exclude
			}

		} else {
			if v == m.Pattern {
				return !m.Exclude
			}
		}
	}
	return m.Exclude
}

func MatchAllOfStringListMatcherList(matchers []*StringListMatcher, value []string) bool {
	for _, matcher := range matchers {
		if !matcher.Match(value) {
			return false
		}
	}
	return true
}
func MatchAnyOfStringListMatcherList(matchers []*StringListMatcher, value []string) bool {
	for _, matcher := range matchers {
		if matcher.Match(value) {
			return true
		}
	}
	return false
}
