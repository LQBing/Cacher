package rules

import "cacher/comparator"

// match cache. ignore some properties, return the same cache. all empty value will return all match
type CompareRule struct {
	Headers comparator.Comparator `json:"headers,omitempty"`
	Cookies comparator.Comparator `json:"cookies,omitempty"`
	Body    comparator.Comparator `json:"body,omitempty"`
}

func NewCompareRule(headers comparator.Comparator, cookies comparator.Comparator, body comparator.Comparator) *CompareRule {
	return &CompareRule{
		Headers: headers,
		Cookies: cookies,
		Body:    body,
	}
}
