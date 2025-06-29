package rules

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"regexp"

	"github.com/jefferyjob/go-easy-utils/v2/anyUtil"
)

func (o *Operation) checkHeaderOption() error {
	if o.Key == "" {
		return errors.New("header key is empty")
	}
	_, ok1 := o.Value.(string)
	_, ok2 := o.Value.([]string)
	if !ok1 && !ok2 {
		log.Println(o.Value)
		return errors.New("invalid header value type " + reflect.TypeOf(o.Value).Kind().String())
	}
	switch o.Operation {
	case "add":
	case "del":
	case "update":
	case "replace":
		if o.Pattern == "" {
			return errors.New("pattern can not beempty for replace operation")
		}
		if !ok1 {
			return errors.New("invalid header value type " + reflect.TypeOf(o.Value).Kind().String() + " for replace operation")
		}
	default:
		return errors.New("invlid header operation " + o.Operation)
	}
	return nil
}
func CheckHeaderOptions(options []Operation) error {
	for i := 0; i < len(options); i++ {
		if err := options[i].checkHeaderOption(); err != nil {
			return err
		}
	}
	return nil
}

func (o *Operation) operateHeader(header *http.Header) error {
	switch o.Operation {
	case "add":
		if _, ok := o.Value.(string); ok {
			header.Add(o.Key, o.Value.(string))
		} else if _, ok := o.Value.([]string); ok {
			for _, vv := range o.Value.([]string) {
				header.Add(o.Key, vv)
			}
		}
	case "del":
		if o.KeyRegex {
			for key := range *header {
				if regexp.MustCompile(o.Key).MatchString(key) {
					if o.Pattern == "" {
						// do delete
						header.Del(key)
					} else {
						for _, v := range (*header)[key] {
							// do delete
							header.Del(key)
							if o.PatternRegex {
								if regexp.MustCompile(o.Pattern).MatchString(v) {
									header.Add(key, v)
								}
							} else {
								if v == o.Pattern {
									header.Add(key, v)
								}
							}
						}
					}
				}
			}
		} else {
			if o.Pattern == "" {
				// do delete
				header.Del(o.Key)
			} else {
				for key := range *header {
					if regexp.MustCompile(o.Key).MatchString(key) {
						// do delete
						header.Del(key)
						for _, v := range (*header)[key] {
							if o.PatternRegex {
								if !regexp.MustCompile(o.Pattern).MatchString(v) {
									header.Add(key, v)
								}
							} else {
								if o.Pattern != v {
									header.Add(key, v)
								}
							}
						}
					}
				}
			}
		}
	case "update":
		_, ok1 := o.Value.(string)
		_, ok2 := o.Value.([]string)
		if o.KeyRegex {
			for key, value := range *header {
				if regexp.MustCompile(o.Key).MatchString(key) {
					if o.Pattern == "" {
						// key re, pattern empty
						// do update
						if ok1 || ok2 {
							if ok1 {
								header.Add(key, o.Value.(string))
							}
							if ok2 {
								for _, vv := range o.Value.([]string) {
									header.Add(key, vv)
								}
							}
						}
					} else {
						for _, v := range value {
							if o.PatternRegex {
								if regexp.MustCompile(o.Pattern).MatchString(v) {
									// key re, pattern re
									// do update
									if ok1 || ok2 {
										header.Del(key)
										if ok1 {
											for _, vv := range value {
												if !regexp.MustCompile(o.Pattern).MatchString(v) {
													header.Add(key, vv)
												}
											}
											header.Add(key, o.Value.(string))
										}
										if ok2 {
											for _, vv := range value {
												if !regexp.MustCompile(o.Pattern).MatchString(v) {
													header.Add(key, vv)
												}
											}
											for _, vv := range o.Value.([]string) {
												header.Add(key, vv)
											}
										}
									}
									break
								}
							} else {
								if v == o.Pattern {
									// key re, pattern not re
									// do update
									if ok1 || ok2 {
										header.Del(key)
										if ok1 {
											for _, vv := range value {
												if vv != o.Pattern {
													header.Add(key, vv)
												}
											}
											header.Add(key, o.Value.(string))
										}
										if ok2 {
											for _, vv := range value {
												if vv != o.Pattern {
													header.Add(key, vv)
												}
											}
											for _, vv := range o.Value.([]string) {
												header.Add(key, vv)
											}
										}
									}
									break
								}
							}
						}
					}

				}
			}
		} else {
			key := o.Key
			// key exist
			if value, ok := (*header)[o.Key]; ok {
				if o.Pattern == "" {
					// key not re, pattern empty
					// do update
					if ok1 || ok2 {
						if ok1 {
							header.Add(key, o.Value.(string))
						}
						if ok2 {
							for _, vv := range o.Value.([]string) {
								header.Add(key, vv)
							}
						}
					}
				} else {
					for _, v := range value {
						if o.PatternRegex {
							if regexp.MustCompile(o.Pattern).MatchString(v) {
								// key not re, pattern re
								// do update
								if ok1 || ok2 {
									header.Del(key)
									if ok1 {
										for _, vv := range value {
											if !regexp.MustCompile(o.Pattern).MatchString(v) {
												header.Add(key, vv)
											}
										}
										header.Add(key, o.Value.(string))
									}
									if ok2 {
										for _, vv := range value {
											if !regexp.MustCompile(o.Pattern).MatchString(v) {
												header.Add(key, vv)
											}
										}
										for _, vv := range o.Value.([]string) {
											header.Add(key, vv)
										}
									}
								}
								break
							}
						} else {
							if v == o.Pattern {
								// key not re, pattern not re
								// do update
								if ok1 || ok2 {
									header.Del(key)
									if ok1 {
										for _, vv := range value {
											if vv != o.Pattern {
												header.Add(key, vv)
											}
										}
										header.Add(key, o.Value.(string))
									}
									if ok2 {
										for _, vv := range value {
											if vv != o.Pattern {
												header.Add(key, vv)
											}
										}
										for _, vv := range o.Value.([]string) {
											header.Add(key, vv)
										}
									}
								}
								break
							}

						}
					}
				}
			} else {
				// key not exist
				// do update as add
				if _, ok := o.Value.(string); ok {
					header.Add(key, o.Value.(string))
				} else if _, ok := o.Value.([]string); ok {
					for _, vv := range o.Value.([]string) {
						header.Add(key, vv)
					}
				}
			}
		}
	case "replace":
		if o.KeyRegex {
			for key, value := range *header {
				if regexp.MustCompile(o.Key).MatchString(key) {
					// do replace
					header.Del(key)
					for _, v := range value {
						if o.PatternRegex {
							if regexp.MustCompile(o.Pattern).MatchString(v) {
								header.Add(key, regexp.MustCompile(o.Pattern).ReplaceAllString(v, anyUtil.AnyToStr(o.Value)))
							} else {
								header.Add(key, v)
							}
						} else {
							if o.Pattern == v {
								header.Add(key, regexp.MustCompile(o.Pattern).ReplaceAllString(v, anyUtil.AnyToStr(o.Value)))
							} else {
								header.Add(key, v)
							}
						}
					}
				}
			}
		} else {
			if value, ok := (*header)[o.Key]; ok {
				key := o.Key
				// do replace
				header.Del(key)
				for _, v := range value {
					if o.PatternRegex {
						if regexp.MustCompile(o.Pattern).MatchString(v) {
							header.Add(key, regexp.MustCompile(o.Pattern).ReplaceAllString(v, anyUtil.AnyToStr(o.Value)))
						} else {
							header.Add(key, v)
						}
					} else {
						if o.Pattern == v {
							header.Add(key, regexp.MustCompile(o.Pattern).ReplaceAllString(v, anyUtil.AnyToStr(o.Value)))
						} else {
							header.Add(key, v)
						}
					}
				}
			}
		}

	default:
		log.Println("unsupported header operation: ", o.Operation)
		return errors.New("unsupported header operation " + o.Operation)
	}
	return nil
}
