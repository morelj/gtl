package function

import (
	"text/template"
)

const valueCategory = "Value"

var valueFuncs = []FunctionSet{
	{
		Category:    valueCategory,
		Syntax:      "exists <value any>",
		Description: []string{"Return true if value is not nil, false otherwise"},
		Functions: template.FuncMap{"exists": func(value any) bool {
			return value != nil
		}},
	},
	{
		Category:    valueCategory,
		Syntax:      "has_value <value any>",
		Description: []string{"Same as exists but also returns false if value is an empty string"},
		Functions:   template.FuncMap{"has_value": hasValue},
	},
	{
		Category:    valueCategory,
		Syntax:      "default <default any> <value any>",
		Description: []string{"If has_value value returns true, returns value otherwise returns default"},
		Functions: template.FuncMap{"default": func(def, value any) any {
			if hasValue(value) {
				return value
			}
			return def
		}},
	},
}

func hasValue(value any) bool {
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
