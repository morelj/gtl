package function

import (
	"text/template"
)

const mathCategory = "Math"

var mathFuncs = []FunctionSet{
	{
		Category:    mathCategory,
		Syntax:      "add|sub|mul|div <v1 int> ... <vN int>",
		Description: []string{"Returns the result of the addition/substraction/multiplication/division of the ints."},
		Functions: template.FuncMap{
			"add": arith(func(acc, v int) int { return acc + v }),
			"sub": arith(func(acc, v int) int { return acc - v }),
			"mul": arith(func(acc, v int) int { return acc * v }),
			"div": arith(func(acc, v int) int { return acc / v }),
		},
	},
}

func arith(f func(acc, v int) int) func(vals ...int) int {
	return func(vals ...int) int {
		if len(vals) == 0 {
			return 0
		}
		res := vals[0]
		for i := 1; i < len(vals); i++ {
			res = f(res, vals[i])
		}
		return res
	}
}
