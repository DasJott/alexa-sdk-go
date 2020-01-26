package alexa

import (
	"encoding/json"
	"strconv"
)

// Attr takes a generic value and converts it with provided functions
type Attr struct {
	name   string
	val    interface{}
	exists bool
}

// Exists determines whether the attribute ectually exists within the collection
func (u *Attr) Exists() bool {
	return u.exists
}

// String gets attribute as string by given key
// Can convert from float32, float64, int, int64
func (u *Attr) String() string {
	switch u.val.(type) {
	case string:
		return u.val.(string)
	case float64:
		return strconv.FormatFloat(u.val.(float64), 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(u.val.(float32)), 'f', -1, 32)
	case int64:
		return strconv.FormatInt(u.val.(int64), 10)
	case int:
		return strconv.FormatInt(int64(u.val.(int)), 10)
	}
	return ""
}

// Bool gets the attribute as bool by given key
func (u *Attr) Bool() bool {
	if b, ok := u.val.(bool); ok {
		return b
	}
	return false
}

// Int gets attribute as int by given key
func (u *Attr) Int() int {
	switch u.val.(type) {
	case string:
		n, _ := strconv.Atoi(u.val.(string))
		return n
	case float64:
		return int(u.val.(float64))
	case float32:
		return int(u.val.(float32))
	case int64:
		return int(u.val.(int64))
	case int:
		return u.val.(int)
	}
	return 0
}

// Int64 gets attribute as int by given key
func (u *Attr) Int64() int64 {
	switch u.val.(type) {
	case string:
		n, _ := strconv.Atoi(u.val.(string))
		return int64(n)
	case float64:
		return int64(u.val.(float64))
	case float32:
		return int64(u.val.(float32))
	case int64:
		return u.val.(int64)
	case int:
		return int64(u.val.(int))
	}
	return 0
}

// Float32 gets attribute as float32 by given key
func (u *Attr) Float32() float32 {
	switch u.val.(type) {
	case string:
		n, _ := strconv.Atoi(u.val.(string))
		return float32(n)
	case float64:
		return float32(u.val.(float64))
	case float32:
		return u.val.(float32)
	case int64:
		return float32(u.val.(int64))
	case int:
		return float32(u.val.(int))
	}
	return 0.0
}

// Float64 gets attribute as float64 by given key
func (u *Attr) Float64() float64 {
	switch u.val.(type) {
	case string:
		n, _ := strconv.Atoi(u.val.(string))
		return float64(n)
	case float64:
		return u.val.(float64)
	case float32:
		return float64(u.val.(float32))
	case int64:
		return float64(u.val.(int64))
	case int:
		return float64(u.val.(int))
	}
	return 0.0
}

// Interface gets attribute as interface{} by given key
func (u *Attr) Interface() interface{} {
	return u.val
}

// Unmarshal gets attribute by given key and unmarshals it into data
func (u *Attr) Unmarshal(data interface{}) error {
	b, err := json.Marshal(u.val)
	if err == nil {
		err = json.Unmarshal(b, data)
	}
	return err
}

// StringArr gets attribute as string array by given key
func (u *Attr) StringArr() []string {
	var sarr []string

	if arr, ok := u.val.([]interface{}); ok {
		if count := len(arr); count > 0 {
			sarr = make([]string, 0, count)
			if _, ok := arr[0].(string); ok {
				for _, a := range arr {
					sarr = append(sarr, a.(string))
				}
			}
		}
	}

	return sarr
}

// R returns a map suitable for TR of localization
// If no key is provided, the attribute name is taken.
// If one or more key is provided, all the keys are added with the attributes value.
func (u *Attr) R(key ...string) map[string]string {
	if count := len(key); count > 0 {
		m := make(map[string]string, count)
		for _, k := range key {
			m[k] = u.String()
		}
		return m
	}
	return map[string]string{u.name: u.String()}
}

// attributes stores the attributes
type attributes map[string]interface{}

// Attr gets or sets attributes. Set more than one value and it will become an array
func (a attributes) Attr(key string, values ...interface{}) *Attr {
	if count := len(values); count > 1 {
		arr := make([]interface{}, count)
		for _, val := range values {
			arr = append(arr, val)
		}
		a[key] = arr
	} else if count > 0 {
		if values[0] == nil {
			if _, e := a[key]; e {
				delete(a, key)
			}
		} else {
			a[key] = values[0]
		}
	}

	val, exists := a[key]
	return &Attr{
		name:   key,
		val:    val,
		exists: exists,
	}
}
