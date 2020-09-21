package function

import (
	"text/template"
)

const valueCategory = "Value"

var valueFuncs = []FunctionSet{
	{
		Category:    valueCategory,
		Syntax:      "exists <value interface{}>",
		Description: []string{"Return true if value is not nil, false otherwise"},
		Functions: template.FuncMap{"exists": func(value interface{}) bool {
			return value != nil
		}},
	},
	{
		Category:    valueCategory,
		Syntax:      "has_value <value interface{}>",
		Description: []string{"Same as exists but also returns false if value is an empty string"},
		Functions:   template.FuncMap{"has_value": hasValue},
	},
	{
		Category:    valueCategory,
		Syntax:      "default <default interface{}> <value interface{}>",
		Description: []string{"If has_value value returns true, returns value otherwise returns default"},
		Functions: template.FuncMap{"default": func(def, value interface{}) interface{} {
			if hasValue(value) {
				return value
			}
			return def
		}},
	},
}

func hasValue(value interface{}) bool {
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
