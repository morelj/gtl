package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// ReleaseName is the name of the release (injected at compile time)
var ReleaseName = "unknown"

// ReleaseDate is the date of the release (injected at compile time)
var ReleaseDate = "unknown"

// Environment contains the data exposed to the template as the dot
type Environment struct {
	Data map[string]interface{}
	Env  map[string]string
}

func parseEnvironment() (map[string]string, error) {
	envRegexp, err := regexp.Compile(`^([[:word:]]+)=(.+)$`)
	if err != nil {
		return nil, err
	}

	env := make(map[string]string)
	for _, value := range os.Environ() {
		match := envRegexp.FindStringSubmatch(value)
		if match != nil {
			env[match[1]] = match[2]
		}
	}
	return env, nil
}

func loadJSON(path string, target interface{}) {
	file, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(file); err != nil {
		panic(err.Error())
	}

	if err = json.Unmarshal(buf.Bytes(), target); err != nil {
		panic(err.Error())
	}
}

func buildEnvironment(dataFiles, dataInline string) *Environment {
	env := Environment{Data: make(map[string]interface{})}
	var err error

	// Environment variables
	if env.Env, err = parseEnvironment(); err != nil {
		panic(err.Error())
	}

	// Data, if any
	if dataFiles != "" {
		files := strings.Split(dataFiles, string(os.PathListSeparator))
		for i := range files {
			loadJSON(files[i], &env.Data)
		}
	}
	if dataInline != "" {
		if err = json.Unmarshal([]byte(dataInline), &env.Data); err != nil {
			panic(err.Error())
		}
	}

	return &env
}

func createTemplate(name string) *template.Template {
	return template.New(name).Funcs(template.FuncMap{
		"split":                 templateSplit,
		"exists":                templateExists,
		"has_value":             templateHasValue,
		"default":               templateDefault,
		"concat":                templateConcat,
		"trim_prefix":           templateTrimPrefix,
		"trim_suffix":           templateTrimSuffix,
		"slice":                 templateSlice,
		"append":                templateAppend,
		"map":                   templateMap,
		"set":                   templateSet,
		"base64_encode":         templateBase64Encode,
		"base64_raw_encode":     templateBase64RawEncode,
		"base64_url_encode":     templateBase64URLEncode,
		"base64_raw_url_encode": templateBase64RawURLEncode,
		"base64_decode":         templateBase64Decode,
		"base64_raw_decode":     templateBase64RawDecode,
		"base64_url_decode":     templateBase64URLDecode,
		"base64_raw_url_decode": templateBase64RawURLDecode,
		"first_match":           templateFirstMatch,
		"filter":                templateFilter,
		"filter_map_value":      templateFilterMapValue,
		"filter_slice_value":    templateFilterSliceValue,
		"filter_eq":             templateFilterEq,
		"filter_not":            templateFilterNot,
		"filter_or":             templateFilterOr,
		"filter_and":            templateFilterAnd,
		"filter_to_int":         templateFilterToInt,
		"filter_to_string":      templateFilterToString,
	})
}

func loadTemplate(source string) *template.Template {
	name := "stdin"
	if source != "-" && source != "" {
		name = filepath.Base(source)
	}
	tmpl := createTemplate(name)

	if source == "-" || source == "" {
		// Load from stdin
		var buf bytes.Buffer
		buf.ReadFrom(os.Stdin)
		tmpl.Parse(buf.String())
		return tmpl
	}

	// Load from the file
	tmpl, err := tmpl.ParseFiles(source)
	if err != nil {
		panic(err.Error())
	}

	return tmpl
}

func main() {
	// Terminate gracefully with an error code when panicking
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, r)
			os.Exit(1)
		}
	}()

	flag.Usage = func() {
		fmt.Printf("%s - Processes Go templates from the command line\n\n", filepath.Base(os.Args[0]))
		fmt.Printf("Usage: %s [options]\n\n", os.Args[0])
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Println("\nPlease see the official Go documentation for the syntax of the templates (https://golang.org/pkg/text/template/)")
		fmt.Println("\nThe value of . (dot) exposed to the template is a struct with the following content:")
		fmt.Println("    .Env  - A map contaning all evironment variables (e.g. .Env.HOME)")
		fmt.Println("    .Data - The data provided using the -d and -D command line flags")
		fmt.Println("\nIn addition to the default features provided by the Go templating language, the following functions are provided:")
		fmt.Println("\n    split <sep string> <value string>")
		fmt.Println("        Splits value on sep and returns a slice containing each part")
		fmt.Println("    exists <value interface{}>")
		fmt.Println("        Return true if value is not nil, false otherwise")
		fmt.Println("    has_value <value interface{}>")
		fmt.Println("        Same as exists but also returns fals if value is an empty string")
		fmt.Println("    default <default interface{}> <value interface{}>")
		fmt.Println("        If has_value value returns true, returns value otherwise returns default")
		fmt.Println("    concat <str1 string> ... <strN string>")
		fmt.Println("        Returns all its arguments concatenated")
		fmt.Println("    trim_prefix <prefix string> <s string>")
		fmt.Println("        Removes the prefix from s. Do nothing if s does not start with prefix")
		fmt.Println("    trim_suffix <suffix string> <s string>")
		fmt.Println("        Removes the suffix from s. Do nothing if s does not end with suffix")
		fmt.Println("    slice <val1 interface{}> ... <valN interface{}>")
		fmt.Println("        Returns a slice containing all the arguments")
		fmt.Println("    append <s []interface{}> <val1 interface{}> ... <valN interface{}>")
		fmt.Println("        Appends val1 to valN to the slice s, and returns the resulting slice")
		fmt.Println("    map <key1 string> <val1 interface{}> ... <keyN string> <valN interface{}>")
		fmt.Println("        Builds a new map with the given keys and values")
		fmt.Println("    set <m map[string]interface{}> <key1 string> <val1 interface{}> ... <keyN string> <valN interface{}>")
		fmt.Println("        Sets the given keys and values to the map m, and returns it")
		fmt.Println("    filter <v map[string]interface{}|[]interface{}> <filter1 FilterFunc> ... <filterN FilterFunc>")
		fmt.Println("        Returns a new map/slice containing the elements matching the filters. Filters are built using filter_* functions")
		fmt.Println("    first_match <v map[string]interface{}|[]interface{}> <filter1 FilterFunc> ... <filterN FilterFunc>")
		fmt.Println("        Returns the first value of v which matches all the filters. Filters are build using filter_* functions")
		fmt.Println("\nAvailable filter functions, for use with filter or first_match:")
		fmt.Println("    filter_map_value <key string> <filter1 FilterFunc> ... <filterN FilterFunc>")
		fmt.Println("        Use with filter or first_match. Returns a FilterFunc which applies filters to one value of the map")
		fmt.Println("    filter_slice_value <index int> <filter1 FilterFunc> ... <filterN FilterFunc>")
		fmt.Println("        Use with filter or first_match. Returns a FilterFunc which applies filters to one value of the slice")
		fmt.Println("    filter_eq <v interface{}>")
		fmt.Println("        Use with filter or first_match. Returns a FilterFunc which checks whether the value equals v")
		fmt.Println("    filter_not <filter FilterFunc>")
		fmt.Println("        Use with filter or first_match. Returns a FilterFunc which negates filter")
		fmt.Println("    filter_or <filter1 FilterFunc> ... <filterN FilterFunc>")
		fmt.Println("        Use with filter or first_match. Returns a FilterFunc which checks if at least one filter matches")
		fmt.Println("    filter_and <filter1 FilterFunc> ... <filterN FilterFunc>")
		fmt.Println("        Use with filter or first_match. Returns a FilterFunc which checks if all filters match")
		fmt.Println("    filter_to_int <filter1 FilterFunc> ... <filterN FilterFunc>")
		fmt.Println("        Use with filter or first_match. Returns a FilterFunc which applies filters using the value converted to an int")
		fmt.Println("    filter_to_string <filter1 FilterFunc> ... <filterN FilterFunc>")
		fmt.Println("        Use with filter or first_match. Returns a FilterFunc which applies filters using the value converted to a string")
		fmt.Printf("\n")
	}

	templateFile := flag.String("i", "", "Source template file (- for stdin), defaults to stdin (mutually exlusive with -t)")
	templateInline := flag.String("t", "", "Specify an inline template (mutually exclusive with -i)")
	outputFile := flag.String("o", "-", "Output file (- for stdout), defaults to stdout")
	dataFiles := flag.String("d", "", fmt.Sprintf("A list of JSON files to load as data, separated with %c", os.PathListSeparator))
	dataInline := flag.String("D", "", "An inline JSON to expose to the template as data")
	version := flag.Bool("version", false, "Show the version number and quit")
	flag.Parse()

	if *version {
		fmt.Printf("gtl build %s on %s\n", ReleaseName, ReleaseDate)
		os.Exit(0)
	}

	if *templateFile != "" && *templateInline != "" {
		panic("-i and -t are mutually exclusive")
	}

	env := buildEnvironment(*dataFiles, *dataInline)

	var tmpl *template.Template
	if *templateInline != "" {
		tmpl = createTemplate("inline")
		tmpl.Parse(*templateInline)
	} else {
		tmpl = loadTemplate(*templateFile)
	}

	out := os.Stdout
	if *outputFile != "-" {
		file, err := os.Create(*outputFile)
		if err != nil {
			panic(err.Error())
		}
		out = file
		defer file.Close()
	}

	if err := tmpl.Execute(out, env); err != nil {
		panic(err.Error())
	}
}
