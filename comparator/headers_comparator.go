package comparator

import (
	"net/http"
	"reflect"
	"regexp"
)

func preProcessHeader(header *http.Header, ignore_list []string, regex bool) {
	header.Del("host")
	header.Del("cookie")
	header.Del("origin")
	header.Del("referer")
	if regex {
		for h := range *header {
			for _, ignore := range ignore_list {
				if regexp.MustCompile(ignore).MatchString(h) {
					header.Del(h)
				}
			}
		}
	} else {
		for _, ignore := range ignore_list {
			if header.Get(ignore) != "" {
				header.Del(ignore)
			}
		}
	}
}

func CompareHeadersWithComparator(pattern_headers http.Header, value_headers http.Header, comparator Comparator) bool {
	if comparator.IgnoreAll {
		return true
	}
	if len(comparator.Match) == 0 {
		comparator.MatchRegex = false
	}
	// remove ignore headers
	preProcessHeader(&pattern_headers, comparator.Ignore, comparator.IgnoreRegex)
	preProcessHeader(&value_headers, comparator.Ignore, comparator.IgnoreRegex)
	// list all headers
	a_map := make(map[string]struct{})
	for k := range pattern_headers {
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
	for k := range value_headers {
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
	// compare headers
	for k := range a_map {
		_, ok1 := pattern_headers[k]
		_, ok2 := value_headers[k]
		if ok1 != ok2 {
			return false
		}
		if !reflect.DeepEqual(pattern_headers[k], value_headers[k]) {
			return false
		}
	}
	return true
}
