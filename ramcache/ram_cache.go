package ramcache

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var RAM_CACHE map[string]map[string][]CacheItem

func Load(cache_file string) error {

	if _, err := os.Stat(cache_file); os.IsNotExist(err) {
		file, err := os.Create(cache_file)
		if err != nil {
			return errors.New("cache file not found")
		}
		defer file.Close()
	}
	data, _ := os.ReadFile(cache_file)
	json.Unmarshal(data, &RAM_CACHE)
	return nil
}

func AddCacheItem(path string, method string, cache CacheItem) error {
	if RAM_CACHE == nil {
		RAM_CACHE = make(map[string]map[string][]CacheItem)
	}
	if _, ok := RAM_CACHE[path]; !ok {
		RAM_CACHE[path] = make(map[string][]CacheItem)
	}
	if _, ok := RAM_CACHE[path][method]; !ok {
		RAM_CACHE[path][method] = []CacheItem{}
	}
	RAM_CACHE[path][method] = append(RAM_CACHE[path][method], cache)
	return nil
}
func Save(cache_file string) error {
	data, _ := json.Marshal(RAM_CACHE)
	err := os.WriteFile(cache_file, data, 0644)
	if err != nil {
		fmt.Println(err)
		return errors.New("save cache failed")
	}
	return nil
}
