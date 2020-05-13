# GTL - CLI Template engine

GTL is a simple command line utility which processes text templates using Go's text/template package.

It has the following features:

* Can use templates from the filesystem, stdin or inline
* Exposes all environment variables to the template, in `.Env`
* Allows to expose JSON files to the template, in `.Data`
* Define several helper functions to be used in templates

The template format documentation is available [here](https://golang.org/pkg/text/template/).

## Install

You can install it using the `go` command:

```bash
go get github.com/morelj/gtl
```

Or just download a pre-compiled binary for your system.

## Usage

### Command line usage

Usage: `gtl [options]`

Options:

* `-h` - Shows the documentation
* `-D string` - An inline JSON to expose to the template as data
* `-d string` - A list of JSON files to load as data, separated with `:`
* `-i string` - Source template file (`-` for stdin), defaults to stdin (mutually exlusive with `-t`)
* `-o string` - Output file (`-` for stdout), defaults to stdout (default `-`)
* `-t string` - Specify an inline template (mutually exclusive with `-i`)
* `-version` - Show the version number and quit

### Template syntax

Please see the official Go documentation for the syntax of the templates (https://golang.org/pkg/text/template/)

The value of . (dot) exposed to the template is a struct with the following content:

* `.Env` - A map contaning all evironment variables (e.g. `.Env.HOME`)
* `.Data` - The data provided using the `-d` and `-D` command line flags

In addition to the default features provided by the Go templating language, the following functions are provided:

* `split <sep string> <value string>` - Splits value on sep and returns a slice containing each part
* `exists <value interface{}>` - Return true if value is not nil, false otherwise
* `has_value <value interface{}>` - Same as exists but also returns fals if value is an empty string
* `default <default interface{}> <value interface{}>`- If has_value value returns true, returns value otherwise returns default
* `concat <str1 string> ... <strN string>` - Returns all its arguments concatenated
* `trim_prefix <prefix string> <s string>` - Removes the prefix from s. Do nothing if s does not start with prefix
* `trim_suffix <suffix string> <s string>` - Removes the suffix from s. Do nothing if s does not end with suffix
* `make_slice <val1 interface{}> ... <valN interface{}>` - Returns a slice containing all the arguments
* `append <s []interface{}> <val1 interface{}> ... <valN interface{}>` - Appends val1 to valN to the slice s, and returns the resulting slice
* `map <key1 string> <val1 interface{}> ... <keyN string> <valN interface{}>` - Builds a new map with the given keys and values
* `set <m map[string]interface{}> <key1 string> <val1 interface{}> ... <keyN string> <valN interface{}>` - Sets the given keys and values to the map m, and returns it
* `to_upper <value string>` - Converts value to upper case
* `to_lower <value string>` - Converts value to lower case
* `to_upper_first <value string>` - Converts the first character of value to upper case and leave the rest untouched
* `to_lower_first <value string>` - Converts the first character of value to lower case and leave the rest untouched
* `base64[_url][_raw]_encode <val string>` - Encodes val in Base64. This function comes in several variants by adding the `_url` and `_raw` tags.
  `_raw` variants remove the `=` padding characters, and `_url` variants use the alternate URL compliant alphabet.
* `base64[_url][_raw]_decode <val string>` - Decodes val from Base64. This function comes in several variants by adding the `_url` and `_raw` tags.
  `_raw` variants remove the `=` padding characters, and `_url` variants use the alternate URL compliant alphabet.
* `filter <v map[string]interface{}|[]interface{}> <filter1 FilterFunc> ... <filterN FilterFunc>` - Returns a new map/slice containing the elements matching the filters. Filters are built using `filter_*` functions
* `first_match <v map[string]interface{}|[]interface{}> <filter1 FilterFunc> ... <filterN FilterFunc>` - Returns the first value of v which matches all the filters. Filters are build using `filter_*` functions

Filter functions (to be used with `filter` or `first_match`):

* `filter_map_value <key string> <filter1 FilterFunc> ... <filterN FilterFunc>` - Returns a FilterFunc which applies filters to one value of the map
* `filter_slice_value <index int> <filter1 FilterFunc> ... <filterN FilterFunc>` - Returns a FilterFunc which applies filters to one value of the slice
* `filter_eq <v interface{}>` - Returns a FilterFunc which checks whether the value equals v
* `filter_not <filter FilterFunc>` - Returns a FilterFunc which negates filter
* `filter_or <filter1 FilterFunc> ... <filterN FilterFunc>` - Returns a FilterFunc which checks if at least one filter matches
* `filter_and <filter1 FilterFunc> ... <filterN FilterFunc>` - Returns a FilterFunc which checks if all filters match
* `filter_to_int <filter1 FilterFunc> ... <filterN FilterFunc>` - Returns a FilterFunc which applies filters using the value converted to an int
* `filter_to_string <filter1 FilterFunc> ... <filterN FilterFunc>` - Returns a FilterFunc which applies filters using the value converted to a string
