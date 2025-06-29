package operator

import (
	"cacher/ramcache"
	"cacher/rules"
	"log"
)

func OperateResponses(operations []rules.Operation, response ramcache.ResponseCache) error {
	for _, operation := range operations {
		err := operation.OperateResponse(response)
		if err != nil {
			log.Println(err.Error())
			continue
		}
	}
	return nil
}
