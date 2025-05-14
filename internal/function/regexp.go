package function

import (
	"regexp"
	"text/template"
)

const regexpCategory = "Regular Expressions"

var regexpFuncs = []FunctionSet{
	{
		Category: regexpCategory,
		Syntax:   "regexp <regexp string>",
		Description: []string{
			"Compiles a regexp (using regexp.MustCompile) and returns the regexp. Standard regexp methods can then be used on it.",
			"See https://golang.org/pkg/regexp/ for details.",
		},
		Functions: template.FuncMap{"regexp": func(v string) *regexp.Regexp {
			return regexp.MustCompile(v)
		}},
	},
}
