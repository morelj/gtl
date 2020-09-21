package function

import (
	"io/ioutil"
	"text/template"
)

const ioCategory = "I/O"

var ioFuncs = []FunctionSet{
	{
		Category:    ioCategory,
		Syntax:      "read_file <filename string>",
		Description: []string{"Reads the given filename and returns its content as a string. Panics if an error occurs"},
		Functions: template.FuncMap{"read_file": func(filename string) string {
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				panic(err)
			}
			return string(data)
		}},
	},
}
