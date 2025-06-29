package main

import (
	"cacher/configs"
	"cacher/operator"
	"cacher/ramcache"
	"cacher/rules"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	errConfig := configs.Init()
	if errConfig != nil {
		fmt.Println(errConfig)
		return
	}
	// load rules
	errRule := rules.Load(operator.GetRulePath())
	if errRule != nil {
		fmt.Println(errRule)
		return
	}
	// load cache
	errCache := ramcache.Load(operator.GetCachePath())
	if errCache != nil {
		fmt.Println(errCache)
		return
	}

	r := gin.Default()
	r.Any("/*path", cacheHandler)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
