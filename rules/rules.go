package rules

import (
	"encoding/json"
	"fmt"
	"os"
)

var RULES []Rule

func Load(rule_file string) error {
	_, err := os.Stat(rule_file)
	if err != nil {
		return err
	}
	data, _ := os.ReadFile(rule_file)
	json.Unmarshal(data, &RULES)
	fmt.Println("rules loaded: ")
	fmt.Println(RULES)
	return nil
}
