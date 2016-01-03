package main

import (
	"bufio"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/yosssi/ace"
	"html/template"
	"net/http"
	"os"
	pa "path"
	"strings"
)

// constants
const version = "0.3"
const license = "MIT"

// Command line interface
var usage string = `ace - Command line utility for the Ace HTML template engine.

Usage:
  ace [-i | --inner=<FILE>] [-m | --map=<FILE>] [-s | --separator=<CHAR>] [-p | --stdout] [-o | --output=<FILE>] [-r | --path=<PATH>] [-w | --httpd] <FILE>
  ace [-h | --help]
  ace [-v | --version]
Options:
  -i --inner=<FILE>     Path to the inner.ace file.
  -m --map=<FILE>       Path to the mappings.map file.
  -s --separator=<CHAR> Separator for key/value map file.
  -p --stdout           Print to stdout.
  -o --output=<FILE>    Write to custom file.
  -w --httpd            Start temporary webserver.
  -r --path=<PATH>	Webserver includes this path.
  -h --help             Show this help.
  -v --version          Display version.
Info:
  Author:       	Antonino Catinello
  Version:      	` + version + `
  License:      	` + license

// Errors:
// 1 = usage
// 2 = template generation
// 3 = FileToMap
// 4 = webserver
func main() {
	// handle options
	args, err := docopt.Parse(usage, nil, true, version, false)

	//fmt.Println(err.Error())
	if err != nil || args["<FILE>"] == nil {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	// middle dot U+00B7 (unicode character)
	// keystroke: alt gr + ,
	var separator string = "\u00B7"

	if len(args["--separator"].([]string)) > 0 {
		separator = args["--separator"].([]string)[0]
	}

	// variables
	var base, inner, output, path string

	base = strings.Split(args["<FILE>"].(string), ".ace")[0]

	if len(args["--inner"].([]string)) > 0 {
		inner = strings.Split(args["--inner"].([]string)[0], ".ace")[0]
	} else {
		inner = ""
	}

	if len(args["--output"].([]string)) > 0 {
		output = args["--output"].([]string)[0]
	} else {
		output = pa.Base(base) + ".html"
	}

	// load, execute, generate ace templates and data
	tpl, err := ace.Load(base, inner, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}

	var data map[string]interface{}

	if len(args["--map"].([]string)) > 0 {
		data = FileToMap(args["--map"].([]string)[0], separator)
	} else {
		data = make(map[string]interface{})
	}

	if len(args["--path"].([]string)) > 0 {
		path = args["--path"].([]string)[0]
	} else {
		path, err = os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(4)
		}
	}

	if args["--stdout"].(bool) {
		if err := tpl.Execute(os.Stdout, data); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(2)
		}
	} else if args["--httpd"].(bool) {
		// handle for static files in ${PWD} eg. css/js/images
		http.Handle("/include/", http.StripPrefix("/include/", http.FileServer(http.Dir(path))))

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, tpl, data)
		})

		if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(4)
		}
	} else {
		w, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0655)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(2)
		}
		defer w.Close()

		if err := tpl.Execute(os.NewFile(w.Fd(), output), data); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(2)
		}
	}

}

// Webserver handle which executes the template by request.
func handler(w http.ResponseWriter, r *http.Request, tpl *template.Template, data map[string]interface{}) {
	if err := tpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// FileToMap opens the fileName and parses the content per line
// and separates key/values by the given unicode separator
// to return a map with the content. Keys are of string and values
// are considered string or []string.
func FileToMap(fileName, separator string) map[string]interface{} {
	// hash table variable
	var data map[string]interface{}
	data = make(map[string]interface{})

	// handle file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
	defer file.Close()

	// scan through file
	scanner := bufio.NewScanner(file)

	// loop
	for scanner.Scan() {
		line := scanner.Text()

		// is the line long enough to be considered?
		if len(line) >= len(separator)+2 {
			// if the line contains a separator, then work with it
			if strings.Contains(line, separator) {
				parts := strings.Split(line, separator)

				// is it multivalue? -> []string
				if len(parts) > 2 {
					data[parts[0]] = []string{}

					for i, v := range parts {
						// don't add the keyname to value
						if i != 0 {
							data[parts[0]] = append(data[parts[0]].([]string), v)
						}
					}
				} else {
					// it is a single value -> string
					data[parts[0]] = parts[1]
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}

	return data
}
