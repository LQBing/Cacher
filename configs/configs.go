package configs

import (
	"cacher/utils"
	"errors"
)

type CacheType string

const (
	LocalCache CacheType = "local"
	RedisCache CacheType = "redis"
)

var Config map[string]string

func Init() error {
	Config = make(map[string]string)
	Config["debug"] = utils.Getenv("DEBUG", "")
	Config["host"] = utils.Getenv("HOST", "0.0.0.0")
	Config["port"] = utils.Getenv("PORT", "80")
	Config["cache_type"] = utils.Getenv("CACHE_TYPE", "local")
	// validate cache type
	if Config["cache_type"] != string(LocalCache) && Config["cache_type"] != string(RedisCache) {
		return errors.New("invalid cache type: " + Config["cache_type"])
	}
	Config["save_path"] = utils.Getenv("SAVE_PATH", "./save")
	Config["rule_file"] = utils.Getenv("RULE_FILE", "rules.json")
	Config["cache_file"] = utils.Getenv("CACHE_FILE", "cache.json")
	Config["log_file"] = utils.Getenv("LOG_FILE", "log.txt")
	Config["API_PATH"] = utils.Getenv("API_PATH", "/a")
	return nil
}
