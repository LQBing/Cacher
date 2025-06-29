package matcher

import (
	"regexp"
)

type StringMatcher struct {
	Pattern string `json:"pattern,omitempty"`
	Regex   bool   `json:"regex,omitempty"`
	Exclude bool   `json:"exclude,omitempty"`
}

func NewStringMatcher(pattern string, regex bool, exclude bool) *StringMatcher {
	return &StringMatcher{
		Pattern: pattern,
		Regex:   regex,
		Exclude: exclude,
	}
}

func (m *StringMatcher) Match(value string) bool {
	if m.Regex {
		match, _ := regexp.MatchString(m.Pattern, value)
		if match {
			return !m.Exclude
		} else {
			return m.Exclude
		}
	} else {
		if value == m.Pattern {
			return !m.Exclude
		} else {
			return m.Exclude
		}
	}
}

func MatchAllOfStringMatcherList(matcherList []StringMatcher, value string) bool {
	for _, matcher := range matcherList {
		if !matcher.Match(value) {
			return false
		}
	}
	return true
}
func MatchAnyOfStringMatcherList(matcherList []StringMatcher, value string) bool {
	for _, matcher := range matcherList {
		if matcher.Match(value) {
			return true
		}
	}
	return false
}
func MatchAnyOfStringList(pattern_list []string, value string, regex bool) bool {
	for _, pattern := range pattern_list {
		if regex {
			match, _ := regexp.MatchString(pattern, value)
			if match {
				return true
			}
		} else {
			if pattern == value {
				return true
			}
		}
	}
	return false
}
