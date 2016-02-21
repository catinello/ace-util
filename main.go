package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	pa "path"
	"strings"

	"github.com/yosssi/ace"
)

var (
	c = new(cli)
	version string
)

// constants
const (
        author = "Antonino Catinello"
        license = "MIT"
        year = "2016"
        copyright = "\u00A9"
)


type cli struct {
	inner string
	mapname string
	separator string
	stdout bool
	output string
	path string
	webserver bool
	file string
	help bool
	version bool
}

func help() {
	var o string = os.Args[0] + ` - Command line utility for the Ace HTML template engine.

Usage:
  ` + os.Args[0] + ` [-i | --inner <FILE>] [-m | --map <FILE>] [-s | --separator <CHAR>] [-t | --stdout] [-o | --output <FILE>] [-p | --path <PATH>] [-w | --httpd] <FILE>
  ` + os.Args[0] + ` [-h | --help]
  ` + os.Args[0] + ` [-v | --version]
Options:
  -i | --inner <FILE>     Path to the inner.ace file.
  -m | --map <FILE>       Path to the mappings.map file.
  -s | --separator <CHAR> Separator for key/value map file.
  -t | --stdout           Print to stdout.
  -o | --output <FILE>    Write to custom file.
  -p | --path <PATH>      Webserver includes this path.
  -w | --httpd            Start temporary webserver.
  -h | --help             Show this help.
  -v | --version          Display version.
Info:
  ` + license + ` license ` + copyright + ` ` + year + ` ` + author + `
  Version: ` + version

	fmt.Fprintln(os.Stderr, o)
        os.Exit(1)
}

func info() {
	fmt.Fprintln(os.Stdout, version)
	os.Exit(0)
}

func init() {
	if len(os.Args) < 2 {
                help()
        }

	flags(os.Args[1:])
}

// Errors:
// 1 = usage
// 2 = template generation
// 3 = FileToMap
// 4 = webserver
func main() {
	if len(c.file) == 0 {
		help()
	}

	// middle dot U+00B7 (unicode character)
	// keystroke: alt gr + ,
	var separator string = "\u00B7"

	if len(c.separator) > 0 {
		separator = c.separator
	}

	// variables
	var base, inner, output, path string
	var tpl *template.Template
	var err error

	base = strings.Split(c.file, ".ace")[0]

	if len(c.inner) > 0 {
		inner = strings.Split(c.inner, ".ace")[0]
	} else {
		inner = ""
	}

	if len(c.output) > 0 {
		output = c.output
	} else {
		output = pa.Base(base) + ".html"
	}

	if c.webserver == false {
		// load, execute, generate ace templates and data
		tpl, err = ace.Load(base, inner, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(2)
		}
	}

	var data map[string]interface{}

	if len(c.mapname) > 0 {
		data = FileToMap(c.mapname, separator)
	} else {
		data = make(map[string]interface{})
	}

	if len(c.path) > 0 {
		path = c.path
	} else {
		path, err = os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(4)
		}
	}

	if c.stdout {
		if err := tpl.Execute(os.Stdout, data); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(2)
		}
	} else if c.webserver {
		// handle for static files in ${PWD} eg. css/js/images
		http.Handle("/include/", http.StripPrefix("/include/", http.FileServer(http.Dir(path))))

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, base, inner, data)
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

func flags(opts []string) {
        na := len(opts)-1

        for i, v := range opts {
                //fmt.Println(i,v)
                switch v {
                case "-v", "--version":
                        c.version = true
                case "-h", "--help":
                        c.help = true
                case "-t", "--stdout":
                        c.stdout = true
                case "-w", "--httpd":
                        c.webserver = true
                case "-i", "--inner":
			if na < i+1 {
				help()
			}
                        c.inner = opts[i+1]
                case "-m", "--map":
			if na < i+1 {
				help()
			}
                        c.mapname = opts[i+1]
                case "-o", "--output":
			if na < i+1 {
				help()
			}
                        c.output = opts[i+1]
                case "-p", "--path":
			if na < i+1 {
				help()
			}
                        c.path = opts[i+1]
               case "-s", "--separator":
			if na < i+1 {
				help()
			}
                        c.separator = opts[i+1]
                default:
                        if i == na {
                                c.file = v
                        }
                }
        }

	if c.help {
		help()
	}

	if c.version {
		info()
	}

	if c.webserver && c.stdout {
		help()
	}
}

// Webserver handle which executes the template by request.
func handler(w http.ResponseWriter, r *http.Request, base, inner string, data map[string]interface{}) {
	// load, execute, generate ace templates and data with dynamic reload
	tpl, err := ace.Load(base, inner, &ace.Options{DynamicReload: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
