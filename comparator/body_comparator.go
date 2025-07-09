package comparator

import (
	"log"
	"reflect"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
)

func CompareBodyWithComparator(body1 any, body2 any, comparator Comparator) bool {
	if comparator.IgnoreAll {
		return true
	}
	if reflect.TypeOf(body1) != reflect.TypeOf(body2) { // not same type, return false
		return false
	}
	bb1 := string(body1.([]byte))
	bb2 := string(body2.([]byte))
	// remove ignore items
	b1, _ := oj.ParseString(bb1)
	b2, _ := oj.ParseString(bb2)
	if len(comparator.Ignore) > 0 {
		for _, ignore := range comparator.Ignore {
			x, err := jp.ParseString(ignore)
			if err != nil {
				log.Panicln(err.Error())
			}
			x.Del(b1)
			x.Del(b2)
		}
	}
	// compare body
	if len(comparator.Match) == 0 { // match regex
		if reflect.DeepEqual(b1, b2) {
			return true
		} else {
			return false
		}
	} else {
		for _, match := range comparator.Match {
			x, err := jp.ParseString(match)
			if err != nil {
				log.Panicln(err.Error())
			}
			bb1 := x.Get(b1)
			bb2 := x.Get(b2)
			if !reflect.DeepEqual(bb1, bb2) {
				return false
			}
		}
	}
	return true
}
