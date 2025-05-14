package function

import (
	"fmt"
	"text/template"
)

const mapSliceCategory = "Maps and slices"

var mapSliceFuncs = []FunctionSet{
	{
		Category:    mapSliceCategory,
		Syntax:      "make_slice <val1 any> ... <valN any>",
		Description: []string{"Returns a slice containing all the arguments"},
		Functions: template.FuncMap{"make_slice": func(e ...any) []any {
			return e
		}},
	},
	{
		Category:    mapSliceCategory,
		Syntax:      "append <s []any> <val1 any> ... <valN any>",
		Description: []string{"Appends val1 to valN to the slice s, and returns the resulting slice"},
		Functions: template.FuncMap{"append": func(s []any, e ...any) []any {
			return append(s, e...)
		}},
	},
	{
		Category:    mapSliceCategory,
		Syntax:      "map <key1 string> <val1 any> ... <keyN string> <valN any>",
		Description: []string{"Builds a new map with the given keys and values"},
		Functions: template.FuncMap{"map": func(kv ...any) map[string]any {
			ret := make(map[string]any)
			mapSet(ret, kv...)
			return ret
		}},
	},
	{
		Category:    mapSliceCategory,
		Syntax:      "set <m map[string]any> <key1 string> <val1 any> ... <keyN string> <valN any>",
		Description: []string{"Sets the given keys and values to the map m, and returns it"},
		Functions:   template.FuncMap{"set": mapSet},
	},
	{
		Category:    mapSliceCategory,
		Syntax:      "filter <v map[string]any|[]any> <filter1 FilterFunc> ... <filterN FilterFunc>",
		Description: []string{"Returns a new map/slice containing the elements matching the filters. Filters are built using filter_* functions"},
		Functions: template.FuncMap{"filter": func(e any, filters ...FilterFunc) any {
			switch v := e.(type) {
			case []any:
				filtered := make([]any, 0, len(v))
				for i := range v {
					if filterAnd(v[i], filters) {
						filtered = append(filtered, v[i])
					}
				}
				return filtered

			case map[string]any:
				filtered := make(map[string]any)
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
		Syntax:      "first_match <v map[string]any|[]any> <filter1 FilterFunc> ... <filterN FilterFunc>",
		Description: []string{"Returns the first value of v which matches all the filters. Filters are build using filter_* functions"},
		Functions: template.FuncMap{"first_match": func(e any, filters ...FilterFunc) any {
			switch v := e.(type) {
			case []any:
				for i := range v {
					if filterAnd(v[i], filters) {
						return v[i]
					}
				}
				return nil

			case map[string]any:
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

func mapSet(m map[string]any, kv ...any) map[string]any {
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
