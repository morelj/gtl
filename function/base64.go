package function

import (
	"encoding/base64"
	"text/template"
)

const base64Category = "Base64"

var base64Funcs = []FunctionSet{
	{
		Category: base64Category,
		Syntax:   "base64[_url][_raw]_encode <val string>",
		Description: []string{
			"Encodes val in Base64. This function comes in several variants by adding the _url and _raw tags.",
			"_raw variants remove the = padding characters, and _url variants use the alternate URL compliant alphabet",
		},
		Functions: template.FuncMap{
			"base64_encode": func(v string) string {
				return base64.StdEncoding.EncodeToString(([]byte)(v))
			},
			"base64_raw_encode": func(v string) string {
				return base64.RawStdEncoding.EncodeToString(([]byte)(v))
			},
			"base64_url_encode": func(v string) string {
				return base64.URLEncoding.EncodeToString(([]byte)(v))
			},
			"base64_raw_url_encode": func(v string) string {
				return base64.RawURLEncoding.EncodeToString(([]byte)(v))
			},
		},
	},
	{
		Category: base64Category,
		Syntax:   "base64[_url][_raw]_decode <val string>",
		Description: []string{
			"Decodes val from Base64. This function comes in several variants by adding the _url and _raw tags.",
			"_raw variants remove the = padding characters, and _url variants use the alternate URL compliant alphabet",
		},
		Functions: template.FuncMap{
			"base64_decode": func(v string) string {
				return base64Decode(v, base64.StdEncoding)
			},
			"base64_raw_decode": func(v string) string {
				return base64Decode(v, base64.RawStdEncoding)
			},
			"base64_url_decode": func(v string) string {
				return base64Decode(v, base64.URLEncoding)
			},
			"base64_raw_url_decode": func(v string) string {
				return base64Decode(v, base64.RawURLEncoding)
			},
		},
	},
}

func base64Decode(v string, e *base64.Encoding) string {
	data, err := e.DecodeString(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}
