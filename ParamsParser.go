package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// DefaultParams is just sugar for passing a map string interface
type DefaultParams map[string]interface{}

// Params is a wrapper for all params
type Params struct {
	Params map[string]*Param
}

// Add adds a key/value pair for a param to the params map
// Allows for chaining
func (p *Params) Add(key string, value *Param) *Params {
	p.Params[key] = value
	return p
}

// Get returns a param based on string key
func (p *Params) Get(key string) *Param {
	if param, ok := p.Params[key]; ok {
		return param
	}
	return &Param{} // Return an empty param
}

// Has checks if a param exists based on string key
func (p *Params) Has(key string) bool {
	_, ok := p.Params[key]
	return ok
}

// ParamsParser generates a slice of param's with default values
func ParamsParser(maps ...map[string]interface{}) *Params {
	result := &Params{map[string]*Param{}}
	for _, params := range maps {
		for key, value := range params {
			if !result.Has(key) {
				result.Add(key, &Param{value})
			}
		}
	}
	return result
}

// Param is the param passed to a task
type Param struct {
	Value interface{}
}

// String returns the param as a string
func (p *Param) String() string {
	switch value := p.Value.(type) {
	case string:
		return value
	case []string:
		return strings.Join(value, "")
	case nil, interface{}, []interface{}: // TODO: Figure out parsing to string
		return ""
	default:
		return fmt.Sprintf("%v", value)
	}
}

// Strings returns the param as a slice of strings
func (p *Param) Strings() []string {
	strings := []string{}
	switch value := p.Value.(type) {
	case []interface{}:
		for _, v := range value {
			switch value := v.(type) {
			case string:
				strings = append(strings, value)
			case []string:
				strings = append(strings, value...)
			case int, int64, float32, float64:
				strings = append(strings, fmt.Sprintf("%v", value))
			}
		}
	case []string:
		return value
	}
	return strings
}

// Ints returns the param as a slice of int64
func (p *Param) Ints() []int64 {
	ints := []int64{}
	switch value := p.Value.(type) {
	case []interface{}:
		for _, v := range value {
			switch value := v.(type) {
			case string:
				if i, err := strconv.Atoi(value); err == nil {
					ints = append(ints, int64(i))
				}
			case int:
				ints = append(ints, int64(value))
			case int64:
				ints = append(ints, value)
			case float32:
				ints = append(ints, int64(value))
			case float64:
				ints = append(ints, int64(value))
			}
		}
	case []string:
		for _, v := range value {
			if i, err := strconv.Atoi(v); err == nil {
				ints = append(ints, int64(i))
			}
		}
	case []int64:
		return value
	}
	return ints
}

// Floats returns the param as a slice of float64
func (p *Param) Floats() []float64 {
	floats := []float64{}
	switch value := p.Value.(type) {
	case []interface{}:
		for _, v := range value {
			switch value := v.(type) {
			case string:
				if f, err := strconv.ParseFloat(value, 64); err == nil {
					floats = append(floats, f)
				}
			case int:
				floats = append(floats, float64(value))
			case int64:
				floats = append(floats, float64(value))
			case float32:
				floats = append(floats, float64(value))
			case float64:
				floats = append(floats, value)
			}
		}
	case []string:
		for _, v := range value {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				floats = append(floats, f)
			}
		}
	case []float64:
		return value
	}
	return floats
}

// Int64 returns the value as an int64 or defaults back to
// the provided default int
func (p *Param) Int64() int64 {
	switch value := p.Value.(type) {
	case string:
		if i, err := strconv.Atoi(value); err == nil {
			return int64(i)
		}
	case int:
		return int64(value)
	case int64:
		return value
	case float32:
		return int64(value)
	case float64:
		return int64(value)
	}
	return 0
}

// Float64 returns the value as an float64 or defaults back to
// the provided default float64
func (p *Param) Float64() float64 {
	switch value := p.Value.(type) {
	case string:
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	case int:
		return float64(value)
	case int64:
		return float64(value)
	case float32:
		return float64(value)
	case float64:
		return value
	}
	return 0
}

// Bool returns the value as bool or defaults back to
// the provided default bool
func (p *Param) Bool() bool {
	switch value := p.Value.(type) {
	case string:
		if f, err := strconv.ParseBool(value); err == nil {
			return f
		}
	case bool:
		return value
	}
	return false
}
