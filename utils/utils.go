package utils

import (
	"net/http"
	"os"
	"regexp"

	"github.com/ohler55/ojg/jp"
)

func Getenv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
func StringInList(list []string, str string) bool {
	for _, item := range list {
		if str == item {
			return true
		}
	}
	return false
}

func CompareStringMapWithKeys(a map[string]string, b map[string]string, keys []string) bool {
	for _, key := range keys {
		if a[key] != b[key] { // compare the values of each key in both maps
			return false
		}
	}
	return true
}

func GetJsonValue(jsonpath string, data any) []any {
	x, _ := jp.ParseString(jsonpath)
	return x.Get(data)
}
func CreateFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func GetMatchStringListFromStringListWithString(list []string, pattern string, regex bool) []string {
	var matchedStrings []string
	for _, str := range list {
		if regex {
			match, _ := regexp.MatchString(pattern, str)
			if match {
				matchedStrings = append(matchedStrings, str)
			}
		} else {
			if str == pattern {
				matchedStrings = append(matchedStrings, str)
			}
		}
	}
	return matchedStrings
}
func GetMatchStringListFromStringListWithStringList(list []string, patterns []string, regex bool) []string {
	var matchedStrings []string
	for _, str := range list {
		for _, pattern := range patterns {
			if regex {
				match, _ := regexp.MatchString(pattern, str)
				if match {
					matchedStrings = append(matchedStrings, str)
				}
			} else {
				if str == pattern {
					matchedStrings = append(matchedStrings, str)
				}
			}
		}
	}
	return matchedStrings
}
func GetFilteredHeader(h http.Header) http.Header {
	h.Del("cookie")
	h.Del("user-agent")
	h.Del("accept-encoding")
	h.Del("referer")
	return h
}
func GetTotalKeysFrom2StringListWithMatchList(list1 []string, list2 []string, patterns []string, regex bool) map[string]struct{} {
	r := make(map[string]struct{})
	p_exist := len(patterns) > 0
	for _, i := range list1 {
		if p_exist {
			for _, pattern := range patterns {
				if regex {
					match, _ := regexp.MatchString(pattern, i)
					if match {
						r[i] = struct{}{}
					}
				} else {
					if pattern == i {
						r[i] = struct{}{}
					}
				}
			}
		} else {
			r[i] = struct{}{}
		}
	}
	for _, i := range list2 {
		if p_exist {
			for _, pattern := range patterns {
				if regex {
					match, _ := regexp.MatchString(pattern, i)
					if match {
						r[i] = struct{}{}
					}
				} else {
					if pattern == i {
						r[i] = struct{}{}
					}
				}
			}
		} else {
			r[i] = struct{}{}
		}
	}
	return r
}
