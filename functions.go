package main

import (
	"fmt"
	"strconv"
	"strings"
)

type FilterFunc func(v interface{}) bool

func templateSplit(sep, s string) []string {
	return strings.Split(s, sep)
}

func templateExists(value interface{}) bool {
	return value != nil
}

func templateHasValue(value interface{}) bool {
	if value == nil {
		return false
	}

	switch t := value.(type) {
	case string:
		return t != ""
	default:
		return true
	}
}

func templateDefault(def, value interface{}) interface{} {
	if templateHasValue(value) {
		return value
	}
	return def
}

func templateConcat(str string, rest ...string) string {
	ret := str
	for i := range rest {
		ret += rest[i]
	}
	return ret
}

func templateTrimSuffix(suffix, s string) string {
	return strings.TrimSuffix(s, suffix)
}

func templateTrimPrefix(prefix, s string) string {
	return strings.TrimPrefix(s, prefix)
}

func templateSlice(e ...interface{}) []interface{} {
	return e
}

func templateAppend(s []interface{}, e ...interface{}) []interface{} {
	return append(s, e...)
}

func templateMap(kv ...interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	templateSet(ret, kv...)
	return ret
}

func templateSet(m map[string]interface{}, kv ...interface{}) map[string]interface{} {
	if len(kv)%2 != 0 {
		panic("Invalid number of arguments")
	}
	for i := 0; i < len(kv); i += 2 {
		k, ok := kv[i].(string)
		if !ok {
			panic("Map keys must be strings")
		}
		m[k] = kv[i+1]
	}
	return m
}

func templateFilter(e interface{}, filters ...FilterFunc) interface{} {
	switch v := e.(type) {
	case []interface{}:
		filtered := make([]interface{}, 0, len(v))
		for i := range v {
			if filterAnd(v[i], filters) {
				filtered = append(filtered, v[i])
			}
		}
		return filtered

	case map[string]interface{}:
		filtered := make(map[string]interface{})
		for k, v := range v {
			if filterAnd(v, filters) {
				filtered[k] = v
			}
		}
		return filtered

	default:
		panic(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func templateFirstMatch(e interface{}, filters ...FilterFunc) interface{} {
	switch v := e.(type) {
	case []interface{}:
		for i := range v {
			if filterAnd(v[i], filters) {
				return v[i]
			}
		}
		return nil

	case map[string]interface{}:
		for _, v := range v {
			if filterAnd(v, filters) {
				return v
			}
		}
	}
	return nil
}

func templateFilterMapValue(k string, filters ...FilterFunc) FilterFunc {
	return func(v interface{}) bool {
		m := v.(map[string]interface{})
		return filterAnd(m[k], filters)
	}
}

func templateFilterSliceValue(i int, filters ...FilterFunc) FilterFunc {
	return func(v interface{}) bool {
		s := v.([]interface{})
		return filterAnd(s[i], filters)
	}
}

func templateFilterEq(v1 interface{}) FilterFunc {
	return func(v2 interface{}) bool {
		return v1 == v2
	}
}

func templateFilterNot(filter FilterFunc) FilterFunc {
	return func(v interface{}) bool {
		return !filter(v)
	}
}

func templateFilterOr(filters ...FilterFunc) FilterFunc {
	return func(v interface{}) bool {
		for _, filter := range filters {
			if filter(v) {
				return true
			}
		}
		return false
	}
}

func templateFilterAnd(filters ...FilterFunc) FilterFunc {
	return func(v interface{}) bool {
		return filterAnd(v, filters)
	}
}

func templateFilterToInt(filters ...FilterFunc) FilterFunc {
	return func(v interface{}) bool {
		switch v := v.(type) {
		case int:
			return filterAnd(v, filters)
		case int8:
			return filterAnd(int(v), filters)
		case int16:
			return filterAnd(int(v), filters)
		case int32:
			return filterAnd(int(v), filters)
		case int64:
			return filterAnd(int(v), filters)
		case float32:
			return filterAnd(int(v), filters)
		case float64:
			return filterAnd(int(v), filters)
		case string:
			iv, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			return filterAnd(iv, filters)
		default:
			panic(fmt.Sprintf("Cannot convert %T to int", v))
		}
	}
}

func templateFilterToString(filters ...FilterFunc) FilterFunc {
	return func(v interface{}) bool {
		return filterAnd(fmt.Sprintf("%v", v), filters)
	}
}

func filterAnd(v interface{}, filters []FilterFunc) bool {
	for _, filter := range filters {
		if !filter(v) {
			return false
		}
	}
	return true
}
