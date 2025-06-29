package rules

import (
	"cacher/ramcache"
	"errors"
)

// type OperationType string

// const (
// 	ADD     OperationType = "add"        // If Key not exist, add. Otherwise do nothing.
// 	DELETE  OperationType = "del"        // If Key match, delete. If Pattern is empty, delete. If Pattern is not empty but not match, do nothing.
// 	UPDATE  OperationType = "update"     // If Key match, update, If Key not exist and KeyRegex is false, add. If Pattern is empty or match with value update, if Pattern is not empty but not match, do nothing.
// 	REPLACE OperationType = "replace"    // If Key match, replace. If Pattern is empty or match with value, replace value with Value. If Pattern is empty but not match, do nothing.
// )

type Operation struct {
	Operation    string `json:"op"`
	Property     string `json:"prop"`
	Key          string `json:"key"`
	KeyRegex     bool   `json:"keyRegex,omitempty"`
	KeyJsonpath  string `json:"keyJsonpath,omitempty"`
	Pattern      string `json:"pattern,omitempty"`
	PatternRegex bool   `json:"valueRegex,omitempty"`
	Value        any    `json:"value,omitempty"`
	AsJson       bool   `json:"asJson,omitempty"` // if true, value will be parsed as json. Only use for property body.
}

func NewOperation(operation string, property string, key string, key_regex bool, key_jsonpath string, pattern string, pattern_regex bool, value any, as_json bool) *Operation {
	return &Operation{
		Operation:    operation,
		Property:     property,
		Key:          key,
		KeyRegex:     key_regex,
		KeyJsonpath:  key_jsonpath,
		Pattern:      pattern,
		PatternRegex: pattern_regex,
		Value:        value,
		AsJson:       as_json,
	}
}

func (o *Operation) OperateRequest(request *ramcache.RequestCache) error {
	switch o.Property {
	case "header":
		o.operateHeader(&request.Header)
	case "cookie":
		o.operateCookie(&request.Cookies)
	case "body":
		o.operateBody(&request.Body)
	default:
		return errors.New("unsupported request property " + o.Property)
	}
	return nil
}
func (o *Operation) OperateResponse(response ramcache.ResponseCache) error {
	switch o.Property {
	case "header":
		o.operateHeader(&response.Header)
	case "body":
		o.operateBody(&response.Body)
	default:
		return errors.New("unsupported resoponse property " + o.Property)
	}
	return nil
}
