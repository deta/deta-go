package base

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

// Append utility, can be used in Updates to append items to a list
//
// value can be single or a list of items
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

// Prepend utility, can be used in Updates to prepend items to a list
//
// value can be single or a list of items
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

// Increment utility, can be used in Updates to increment/decrement a numerical value
//
// if value is negative, the value is subtracted in the update operation
func (u *util) Increment(value interface{}) *incrementUtil {
	return &incrementUtil{
		value: value,
	}
}

// Trim utility, can be used in Updates to trim a field from an item
func (u *util) Trim() *trimUtil {
	return &trimUtil{}
}
