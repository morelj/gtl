package function

import (
	"fmt"
	"text/template"
)

const mapSliceCategory = "Maps and slices"

var mapSliceFuncs = []FunctionSet{
	{
		Category:    mapSliceCategory,
		Syntax:      "make_slice <val1 interface{}> ... <valN interface{}>",
		Description: []string{"Returns a slice containing all the arguments"},
		Functions: template.FuncMap{"make_slice": func(e ...interface{}) []interface{} {
			return e
		}},
	},
	{
		Category:    mapSliceCategory,
		Syntax:      "append <s []interface{}> <val1 interface{}> ... <valN interface{}>",
		Description: []string{"Appends val1 to valN to the slice s, and returns the resulting slice"},
		Functions: template.FuncMap{"append": func(s []interface{}, e ...interface{}) []interface{} {
			return append(s, e...)
		}},
	},
	{
		Category:    mapSliceCategory,
		Syntax:      "map <key1 string> <val1 interface{}> ... <keyN string> <valN interface{}>",
		Description: []string{"Builds a new map with the given keys and values"},
		Functions: template.FuncMap{"map": func(kv ...interface{}) map[string]interface{} {
			ret := make(map[string]interface{})
			mapSet(ret, kv...)
			return ret
		}},
	},
	{
		Category:    mapSliceCategory,
		Syntax:      "set <m map[string]interface{}> <key1 string> <val1 interface{}> ... <keyN string> <valN interface{}>",
		Description: []string{"Sets the given keys and values to the map m, and returns it"},
		Functions:   template.FuncMap{"set": mapSet},
	},
	{
		Category:    mapSliceCategory,
		Syntax:      "filter <v map[string]interface{}|[]interface{}> <filter1 FilterFunc> ... <filterN FilterFunc>",
		Description: []string{"Returns a new map/slice containing the elements matching the filters. Filters are built using filter_* functions"},
		Functions: template.FuncMap{"filter": func(e interface{}, filters ...FilterFunc) interface{} {
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
		}},
	},
	{
		Category:    mapSliceCategory,
		Syntax:      "first_match <v map[string]interface{}|[]interface{}> <filter1 FilterFunc> ... <filterN FilterFunc>",
		Description: []string{"Returns the first value of v which matches all the filters. Filters are build using filter_* functions"},
		Functions: template.FuncMap{"first_match": func(e interface{}, filters ...FilterFunc) interface{} {
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
		}},
	},
}

func mapSet(m map[string]interface{}, kv ...interface{}) map[string]interface{} {
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
