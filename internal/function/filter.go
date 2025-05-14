package function

import (
	"fmt"
	"strconv"
	"text/template"
)

type FilterFunc func(v any) bool

const filterCategory = "Filter"

var filterFuncs = []FunctionSet{
	{
		Category:    filterCategory,
		Syntax:      "filter_map_value <key string> <filter1 FilterFunc> ... <filterN FilterFunc>",
		Description: []string{"Use with filter or first_match. Returns a FilterFunc which applies filters to one value of the map"},
		Functions: template.FuncMap{"filter_map_value": func(k string, filters ...FilterFunc) FilterFunc {
			return func(v any) bool {
				m := v.(map[string]any)
				return filterAnd(m[k], filters)
			}
		}},
	},
	{
		Category:    filterCategory,
		Syntax:      "filter_slice_value <index int> <filter1 FilterFunc> ... <filterN FilterFunc>",
		Description: []string{"Use with filter or first_match. Returns a FilterFunc which applies filters to one value of the slice"},
		Functions: template.FuncMap{"fliter_slice_value": func(i int, filters ...FilterFunc) FilterFunc {
			return func(v any) bool {
				s := v.([]any)
				return filterAnd(s[i], filters)
			}
		}},
	},
	{
		Category:    filterCategory,
		Syntax:      "filter_eq <v any>",
		Description: []string{"Use with filter or first_match. Returns a FilterFunc which checks whether the value equals v"},
		Functions: template.FuncMap{"filter_eq": func(v1 any) FilterFunc {
			return func(v2 any) bool {
				return v1 == v2
			}
		}},
	},
	{
		Category:    filterCategory,
		Syntax:      "filter_not <filter FilterFunc>",
		Description: []string{"Use with filter or first_match. Returns a FilterFunc which negates filter"},
		Functions: template.FuncMap{"filter_not": func(filter FilterFunc) FilterFunc {
			return func(v any) bool {
				return !filter(v)
			}
		}},
	},
	{
		Category:    filterCategory,
		Syntax:      "filter_or <filter1 FilterFunc> ... <filterN FilterFunc>",
		Description: []string{"Use with filter or first_match. Returns a FilterFunc which checks if at least one filter matches"},
		Functions: template.FuncMap{"filter_or": func(filters ...FilterFunc) FilterFunc {
			return func(v any) bool {
				for _, filter := range filters {
					if filter(v) {
						return true
					}
				}
				return false
			}
		}},
	},
	{
		Category:    filterCategory,
		Syntax:      "filter_and <filter1 FilterFunc> ... <filterN FilterFunc>",
		Description: []string{"Use with filter or first_match. Returns a FilterFunc which checks if all filters match"},
		Functions: template.FuncMap{"filter_and": func(filters ...FilterFunc) FilterFunc {
			return func(v any) bool {
				return filterAnd(v, filters)
			}
		}},
	},
	{
		Category:    filterCategory,
		Syntax:      "filter_to_int <filter1 FilterFunc> ... <filterN FilterFunc>",
		Description: []string{"Use with filter or first_match. Returns a FilterFunc which applies filters using the value converted to an int"},
		Functions: template.FuncMap{"filter_to_int": func(filters ...FilterFunc) FilterFunc {
			return func(v any) bool {
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
		}},
	},
	{
		Category:    filterCategory,
		Syntax:      "filter_to_string <filter1 FilterFunc> ... <filterN FilterFunc>",
		Description: []string{"Use with filter or first_match. Returns a FilterFunc which applies filters using the value converted to a string"},
		Functions: template.FuncMap{"filter_to_string": func(filters ...FilterFunc) FilterFunc {
			return func(v any) bool {
				return filterAnd(fmt.Sprintf("%v", v), filters)
			}
		}},
	},
}

func filterAnd(v any, filters []FilterFunc) bool {
	for _, filter := range filters {
		if !filter(v) {
			return false
		}
	}
	return true
}
