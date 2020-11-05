package deta

import "reflect"

// append util
type appendUtil struct {
	value interface{}
}

// prepend util
type prependUtil struct {
	value interface{}
}

// increment util
type incrementUtil struct {
	value interface{}
}

type trimUtil struct{}

// utility struct
type util struct{}

// Append utility
func (u *util) Append(value interface{}) *appendUtil {
	switch reflect.ValueOf(value).Kind() {
	case reflect.Array, reflect.Slice:
		return &appendUtil{
			value: value,
		}
	case reflect.Ptr:
		switch reflect.Indirect(reflect.ValueOf(value)).Kind() {
		case reflect.Array, reflect.Slice:
			return &appendUtil{
				value: value,
			}
		default:
			return &appendUtil{
				value: []interface{}{value},
			}
		}
	default:
		return &appendUtil{
			value: []interface{}{value},
		}
	}
}

// Prepend utility
func (u *util) Prepend(value interface{}) *prependUtil {
	switch reflect.ValueOf(value).Kind() {
	case reflect.Array, reflect.Slice:
		return &prependUtil{
			value: value,
		}
	case reflect.Ptr:
		switch reflect.Indirect(reflect.ValueOf(value)).Kind() {
		case reflect.Array, reflect.Slice:
			return &prependUtil{
				value: value,
			}
		default:
			return &prependUtil{
				value: []interface{}{value},
			}
		}
	default:
		return &prependUtil{
			value: []interface{}{value},
		}
	}
}

// Increment utility
func (u *util) Increment(value interface{}) *incrementUtil {
	return &incrementUtil{
		value: value,
	}
}

// Trim utility
func (u *util) Trim() *trimUtil {
	return &trimUtil{}
}
