package rules

import (
	"errors"
	"log"
	"reflect"
	"regexp"

	"github.com/jefferyjob/go-easy-utils/anyUtil"
)

func (o *Operation) checkCookieOption() error {
	if o.Key == "" {
		return errors.New("cookie key is empty")
	}
	_, ok := o.Value.(string)
	if !ok {
		log.Println(o.Value)
		return errors.New("invalid cookie value type" + reflect.TypeOf(o.Value).Kind().String())
	}
	switch o.Operation {
	case "add":
	case "del":
	case "update":
	case "replace":
		if o.Pattern == "" {
			return errors.New("pattern can not be empty for cookie operation")
		}
	default:
		return errors.New("invlid cookie operation " + o.Operation)
	}
	return nil
}
func CheckCookieOptions(options []Operation) error {
	for i := 0; i < len(options); i++ {
		if err := options[i].checkCookieOption(); err != nil {
			return err
		}
	}
	return nil
}

func (o *Operation) operateCookie(cookie *map[string]string) error {
	switch o.Operation {
	case "add":
		if _, ok := (*cookie)[o.Key]; !ok {
			(*cookie)[o.Key] = anyUtil.AnyToStr(o.Value)
		}
	case "del":
		if o.KeyRegex {
			for key := range *cookie {
				if regexp.MustCompile(o.Key).MatchString(key) {
					if o.Pattern == "" {
						delete(*cookie, key)
					} else {
						if o.PatternRegex {
							if regexp.MustCompile(o.Pattern).MatchString((*cookie)[key]) {
								delete(*cookie, key)
							}
						} else {
							if o.Pattern == (*cookie)[o.Key] {
								delete(*cookie, key)
							}
						}
					}
				}
			}

		} else {
			if o.Pattern == "" {
				delete(*cookie, o.Key)
			} else {
				if regexp.MustCompile(o.Pattern).MatchString((*cookie)[o.Key]) {
					delete(*cookie, o.Key)
				}
			}
		}
	case "update":
		if o.KeyRegex {
			for key, value := range *cookie {
				if regexp.MustCompile(o.Key).MatchString(key) {
					if o.Pattern == "" {
						(*cookie)[key] = anyUtil.AnyToStr(o.Value)
					} else {
						if o.PatternRegex {
							if regexp.MustCompile(o.Pattern).MatchString(value) {
								(*cookie)[key] = anyUtil.AnyToStr(o.Value)
							}
						} else {
							if value == o.Pattern {
								(*cookie)[key] = anyUtil.AnyToStr(o.Value)
							}
						}
					}
				}
			}
		} else {
			if o.Pattern == "" {
				(*cookie)[o.Key] = anyUtil.AnyToStr(o.Value)
			} else {
				for key, value := range *cookie {
					if o.PatternRegex {
						if regexp.MustCompile(o.Pattern).MatchString(value) {
							(*cookie)[key] = anyUtil.AnyToStr(o.Value)
						}
					} else {
						if key == o.Pattern {
							(*cookie)[key] = anyUtil.AnyToStr(o.Value)
							break
						}
					}
				}
			}
		}
	case "replace":
		if o.KeyRegex {
			for key, value := range *cookie {
				if regexp.MustCompile(key).MatchString(key) {
					if o.PatternRegex {
						if regexp.MustCompile(o.Pattern).MatchString(value) {
							(*cookie)[key] = regexp.MustCompile(o.Pattern).ReplaceAllString(o.Pattern, anyUtil.AnyToStr(o.Value))
						}
					} else {
						if value == o.Pattern {
							(*cookie)[key] = regexp.MustCompile(o.Pattern).ReplaceAllString(o.Pattern, anyUtil.AnyToStr(o.Value))
						}
					}
				}
			}
		} else {
			for key, value := range *cookie {
				if o.PatternRegex {
					if regexp.MustCompile(o.Pattern).MatchString(value) {
						(*cookie)[key] = regexp.MustCompile(o.Pattern).ReplaceAllString(o.Pattern, anyUtil.AnyToStr(o.Value))
					}
				} else {
					if key == o.Pattern {
						(*cookie)[key] = regexp.MustCompile(o.Pattern).ReplaceAllString(o.Pattern, anyUtil.AnyToStr(o.Value))
						break
					}
				}
			}
		}
	default:
		log.Println("unsupported cookie operation", o.Operation)
		return errors.New("unsupported cookie operation" + o.Operation)
	}
	return nil
}
