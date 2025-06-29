package operator

import (
	"cacher/ramcache"
	"cacher/rules"
	"log"
)

func OperateRequests(operations []rules.Operation, request *ramcache.RequestCache) error {
	request.Header.Del("Cookie")
	for _, operation := range operations {
		err := operation.OperateRequest(request)
		if err != nil {
			log.Println(err.Error())
			continue
		}
	}
	return nil
}
