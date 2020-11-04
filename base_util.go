package deta

// append util
type appendUtil struct {
	value []interface{}
}

// prepend util
type prependUtil struct {
	value []interface{}
}

// increment util
type incrementUtil struct {
	value interface{}
}

type trimUtil struct{}

// utility struct
type util struct{}

// Append utility
func (u *util) Append(value []interface{}) *appendUtil {
	return &appendUtil{
		value: value,
	}
}

// Prepend utility
func (u *util) Prepend(value []interface{}) *prependUtil {
	return &prependUtil{
		value: value,
	}
}

// Increment utility
func (u *util) Increment(value interface{}) *incrementUtil {
	return &incrementUtil{
		value: value,
	}
}

// Trim utility
func (u *util) TrimUtil() *trimUtil {
	return &trimUtil{}
}
