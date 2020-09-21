package function

import (
	"regexp"
	"strings"
	"text/template"
)

const stringCategory = "String"

var lettersRegexp = regexp.MustCompile(`[[:alpha:]]+`)
var nonLettersRegexp = regexp.MustCompile(`[[:^alpha]]+`)
var camelCaseRegexp = regexp.MustCompile(`[[:upper:]][[:lower:]]*`)

var stringFuncs = []FunctionSet{
	{
		Category:    stringCategory,
		Syntax:      "split <sep string> <value string>",
		Description: []string{"Splits value on sep and returns a slice containing each part"},
		Functions: template.FuncMap{"split": func(sep, s string) []string {
			return strings.Split(s, sep)
		}},
	},
	{
		Category:    stringCategory,
		Syntax:      "concat <str1 string> ... <strN string>",
		Description: []string{"Returns all its arguments concatenated"},
		Functions: template.FuncMap{"concat": func(str string, rest ...string) string {
			ret := str
			for i := range rest {
				ret += rest[i]
			}
			return ret
		}},
	},
	{
		Category:    stringCategory,
		Syntax:      "trim_prefix <prefix string> <s string>",
		Description: []string{"Removes the prefix from s. Do nothing if s does not start with prefix"},
		Functions: template.FuncMap{"trim_prefix": func(suffix, s string) string {
			return strings.TrimPrefix(s, suffix)
		}},
	},
	{
		Category:    stringCategory,
		Syntax:      "trim_suffix <suffix string> <s string>",
		Description: []string{" Removes the suffix from s. Do nothing if s does not end with suffix"},
		Functions: template.FuncMap{"trim_suffix": func(suffix, s string) string {
			return strings.TrimSuffix(s, suffix)
		}},
	},
	{
		Category:    stringCategory,
		Syntax:      "to_upper <value string>",
		Description: []string{"Converts value to upper case"},
		Functions:   template.FuncMap{"to_upper": strings.ToUpper},
	},
	{
		Category:    stringCategory,
		Syntax:      "to_lower <value string>",
		Description: []string{"Converts value to lower case"},
		Functions:   template.FuncMap{"to_lower": strings.ToLower},
	},
	{
		Category:    stringCategory,
		Syntax:      "to_upper_first <value string>",
		Description: []string{"Converts the first character of value to upper case and leave the rest untouched"},
		Functions: template.FuncMap{"to_upper_first": func(v string) string {
			if v == "" {
				return ""
			}
			return strings.ToUpper(v[0:1]) + v[1:]
		}},
	},
	{
		Category:    stringCategory,
		Syntax:      "to_lower_first <value string>",
		Description: []string{"Converts the first character of value to lower case and leave the rest untouched"},
		Functions: template.FuncMap{"to_lower_first": func(v string) string {
			if v == "" {
				return ""
			}
			return strings.ToLower(v[0:1]) + v[1:]
		}},
	},
	{
		Category:    stringCategory,
		Syntax:      "replace <old string> <new string> <n int> <s string>",
		Description: []string{"Returns a copy of the string s with the first n non-overlapping instances of old replaced by new."},
		Functions: template.FuncMap{"replace": func(old, new string, n int, s string) string {
			return strings.Replace(s, old, new, n)
		}},
	},
	{
		Category:    stringCategory,
		Syntax:      "replace_all <old string> <new string> <s string>",
		Description: []string{"Returns a copy of the string s with all non-overlapping instances of old replaced by new."},
		Functions: template.FuncMap{"replace_all": func(old, new, s string) string {
			return strings.ReplaceAll(s, old, new)
		}},
	},
	{
		Category:    stringCategory,
		Syntax:      "to_camel_case <s string>",
		Description: []string{"Converts a snake_case string to CamelCase"},
		Functions: template.FuncMap{"to_camel_case": func(s string) string {
			s = lettersRegexp.ReplaceAllStringFunc(s, func(match string) string {
				if len(match) > 0 {
					return strings.ToUpper(match[0:1]) + strings.ToLower(match[1:])
				}
				return ""
			})
			return nonLettersRegexp.ReplaceAllString(s, "")
		}},
	},
	{
		Category:    stringCategory,
		Syntax:      "to_snake_case <s string>",
		Description: []string{"Converts a CamelCase string to snake_case"},
		Functions: template.FuncMap{"to_snake_case": func(s string) string {
			matches := camelCaseRegexp.FindAllString(s, -1)
			for i := range matches {
				matches[i] = strings.ToLower(matches[i])
			}
			return strings.Join(matches, "_")
		}},
	},
}
