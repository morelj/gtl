package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/morelj/gtl/internal/function"
)

// ReleaseName is the name of the release (injected at compile time)
var ReleaseName = "unknown"

// ReleaseCommit is the commit of the release (injected at compile time)
var ReleaseCommit = "unknown"

// ReleaseDate is the date of the release (injected at compile time)
var ReleaseDate = "unknown"

// Environment contains the data exposed to the template as the dot
type Environment struct {
	Data map[string]any
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

func loadJSON(path string, target any) {
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
	env := Environment{Data: make(map[string]any)}
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
	funcs := template.FuncMap{}
	for i := range function.Functions {
		maps.Copy(funcs, function.Functions[i].Functions)
	}
	return template.New(name).Funcs(funcs)
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
		fmt.Printf("\nIn addition to the default features provided by the Go templating language, the following functions are provided:\n\n")

		for _, group := range function.Functions.ByCategory() {
			if len(group) > 0 {
				fmt.Printf("%s functions\n\n", group[0].Category)
				for _, f := range group {
					fmt.Printf("  %s\n", f.Syntax)
					for _, line := range f.Description {
						fmt.Printf("    %s\n", line)
					}
				}
				fmt.Printf("\n")
			}
		}
	}

	templateFile := flag.String("i", "", "Source template file (- for stdin), defaults to stdin (mutually exlusive with -t)")
	templateInline := flag.String("t", "", "Specify an inline template (mutually exclusive with -i)")
	outputFile := flag.String("o", "-", "Output file (- for stdout), defaults to stdout")
	dataFiles := flag.String("d", "", fmt.Sprintf("A list of JSON files to load as data, separated with %c", os.PathListSeparator))
	dataInline := flag.String("D", "", "An inline JSON to expose to the template as data")
	version := flag.Bool("version", false, "Show the version number and quit")
	flag.Parse()

	if *version {
		fmt.Printf("gtl release %s (commit %s) on %s\n", ReleaseName, ReleaseCommit, ReleaseDate)
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
