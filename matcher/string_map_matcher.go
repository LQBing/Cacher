package matcher

import "regexp"

type StringMapMatcher struct {
	Key        string `json:"key"`
	KeyRegex   bool   `json:"keyRegex,omitempty"`
	Value      string `json:"value,omitempty"`
	ValueRegex bool   `json:"valueRegex,omitempty"`
	Exist      bool   `json:"exist,omitempty"`
	Exclude    bool   `json:"exclude,omitempty"`
}

func NewStringMapMatcher(key string, key_regex bool, value string, value_regex bool, exist bool, exclude bool) *StringMapMatcher {
	return &StringMapMatcher{
		Key:        key,
		KeyRegex:   key_regex,
		Value:      value,
		ValueRegex: value_regex,
		Exist:      exist,
		Exclude:    exclude,
	}
}

func (m *StringMapMatcher) Match(value map[string]string) bool {
	for k, v := range value {
		if m.KeyRegex {
			match, _ := regexp.MatchString(m.Key, k)
			if match {
				if m.Exist {
					return !m.Exclude
				} else {
					if m.ValueRegex {
						match, _ := regexp.MatchString(m.Value, v)
						if match {
							return !m.Exclude
						}
					}
				}
			}
		} else {
			if k == m.Key {
				if m.Exist {
					return !m.Exclude
				} else {
					if m.ValueRegex {
						match, _ := regexp.MatchString(m.Value, v)
						if match {
							return !m.Exclude
						}
					}
				}
			}
		}
	}
	return m.Exclude
}

func MatchAllOfStringMapMatcherList(matcher_list []*StringMapMatcher, value map[string]string) bool {
	for _, matcher := range matcher_list {
		if !matcher.Match(value) {
			return false
		}
	}
	return true
}

func MatchAnyOfNewStringMapMatcherList(matcher_list []*StringMapMatcher, value map[string]string) bool {
	for _, matcher := range matcher_list {
		if matcher.Match(value) {
			return true
		}
	}
	return false
}
