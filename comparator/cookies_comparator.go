package comparator

import (
	"reflect"
	"regexp"
)

func preProcessCookies(cookies *map[string]string, ignore_list []string, regex bool) {
	if regex {
		for c := range *cookies {
			for _, ignore := range ignore_list {
				if regexp.MustCompile(ignore).MatchString(c) {
					delete(*cookies, c)
				}
			}
		}
	} else {
		for c := range *cookies {
			for _, ignore := range ignore_list {
				if c == ignore {
					delete(*cookies, c)
				}
			}
		}
	}
}

func CompareCookiesWithComparator(pattern_cookies map[string]string, value_cookies map[string]string, comparator Comparator) bool {
	if comparator.IgnoreAll {
		return true
	}
	if len(comparator.Match) == 0 {
		comparator.MatchRegex = false
	}
	// remove ignore cookies
	preProcessCookies(&pattern_cookies, comparator.Ignore, comparator.IgnoreRegex)
	preProcessCookies(&value_cookies, comparator.Ignore, comparator.IgnoreRegex)
	// list all cookies
	a_map := make(map[string]struct{})
	for k := range pattern_cookies {
		if len(comparator.Match) == 0 {
			a_map[k] = struct{}{}
		} else {
			for _, m := range comparator.Match {
				if comparator.MatchRegex {
					match, _ := regexp.MatchString(m, k)
					if match {
						a_map[k] = struct{}{}
					}
				} else {
					if k == m {
						a_map[k] = struct{}{}
					}
				}
			}
		}
	}
	for k := range value_cookies {
		if len(comparator.Match) == 0 {
			a_map[k] = struct{}{}
		} else {
			for _, m := range comparator.Match {
				if comparator.MatchRegex {
					match, _ := regexp.MatchString(m, k)
					if match {
						a_map[k] = struct{}{}
					}
				} else {
					if k == m {
						a_map[k] = struct{}{}
					}
				}
			}
		}
	}
	// compare cookies
	for k := range a_map {
		_, ok1 := pattern_cookies[k]
		_, ok2 := value_cookies[k]
		if ok1 != ok2 {
			return false
		}
		if !reflect.DeepEqual(pattern_cookies[k], value_cookies[k]) {
			return false
		}
	}
	return true
}
