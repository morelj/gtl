package function

import "text/template"

type FunctionSet struct {
	Category    string
	Syntax      string
	Description []string
	Functions   template.FuncMap
}

type Library []FunctionSet

var Functions Library

func (l Library) ByCategory() [][]*FunctionSet {
	indices := map[string]int{}
	groups := [][]*FunctionSet{}

	for i := range l {
		index, ok := indices[l[i].Category]
		var group []*FunctionSet
		if !ok {
			index = len(groups)
			group = []*FunctionSet{}
			groups = append(groups, group)
			indices[l[i].Category] = index
		} else {
			group = groups[index]
		}

		group = append(group, &l[i])
		groups[index] = group
	}

	return groups
}

func init() {
	Functions = append(Functions, valueFuncs...)
	Functions = append(Functions, stringFuncs...)
	Functions = append(Functions, regexpFuncs...)
	Functions = append(Functions, mathFuncs...)
	Functions = append(Functions, base64Funcs...)
	Functions = append(Functions, ioFuncs...)
	Functions = append(Functions, mapSliceFuncs...)
	Functions = append(Functions, filterFuncs...)
}
