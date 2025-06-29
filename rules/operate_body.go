package rules

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
)

func (o *Operation) checkBodyOption() error {
	if o.AsJson {
		switch o.Operation {
		case "del":
		case "update":
		case "replace":
			if o.Pattern == "" {
				log.Println("body option check failed: pattern is empty")
				return errors.New("pattern is empty")
			}
			if _, ok := o.Value.(string); !ok {
				log.Println("body option check failed: invalid value type")
				return errors.New("invalid value type")
			}
		default:
			log.Println("body option check failed: invalid operation type")
			return errors.New("invalid operation type")
		}
	} else {
		switch o.Operation {
		case "del":
		case "update":
			if _, ok := o.Value.(string); !ok {
				log.Println("body option check failed: invalid value type")
				return errors.New("invalid value type")
			}
		case "replace":
			if o.Pattern == "" {
				log.Println("body option check failed: pattern is empty")
				return errors.New("pattern is empty")
			}
			if _, ok := o.Value.(string); !ok {
				log.Println("body option check failed: invalid value type")
				return errors.New("invalid value type")
			}
		default:
			log.Println("body option check failed: invalid operation type")
			return errors.New("invalid operation type")
		}
	}
	return nil
}
func CheckBodyOptions(options []Operation) error {
	for i := 0; i < len(options); i++ {
		if err := options[i].checkBodyOption(); err != nil {
			log.Println("body option check failed: ", options[i].Operation, options[i].Key, options[i].Value)
			return err
		}
	}
	return nil
}
func (o *Operation) operateBody(body_byte *[]byte) error {
	if o.AsJson {
		body, err := oj.ParseString(string(*body_byte))
		if err != nil {
			log.Println(err.Error())
			return err
		}
		var key string
		if o.Key == "" {
			key = "$"
		} else {
			key = o.Key
		}

		x, err := jp.ParseString(key)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		if len(x) == 0 {
			return nil
		}
		switch o.Operation {
		case "del":
			err := x.Del(body)
			if err != nil {
				log.Println(err.Error())
				return err
			}
		case "update":
			err := x.Set(body, o.Value)
			if err != nil {
				log.Println(err.Error())
				return err
			}
		case "replace":
			results := x.Get(body)
			for _, result := range results {
				if result == o.Pattern {
					err := x.Set(body, o.Value)
					if err != nil {
						log.Println(err.Error())
						return err
					}
				}
			}
		default:
			log.Println("operateBody: invalid operation type")
			return errors.New("invalid operation type")
		}
		nb, err := json.Marshal(body)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		*body_byte = nb
		return nil
	} else {
		body := string(*body_byte)
		fmt.Print(body)
		switch o.Operation {
		case "del":
			*body_byte = make([]byte, 0)
		case "update":
			value, ok := o.Value.(string)
			if ok {
				*body_byte = []byte(value)
			} else {
				return errors.New("invalid value type")
			}
		case "replace":
			value, ok := o.Value.(string)
			if ok {
				v := regexp.MustCompile(o.Pattern).ReplaceAllString(body, value)
				*body_byte = []byte(v)
			} else {
				return errors.New("invalid body value type" + reflect.TypeOf(o.Value).Kind().String())
			}
		default:
			return errors.New("unsupported body operation " + o.Operation)
		}
		return nil
	}
}
