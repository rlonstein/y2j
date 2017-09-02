//
// y2j -- hack to convert YAML to JSON
//
// released under the MIT License
//
// Copyright (c) 2017, Ross Lonstein
//

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/docopt/docopt-go"
	"gopkg.in/yaml.v2"
)

func main() {
	usage := `Convert yaml to json

Usage:
  y2j [-pn] [FILE] ...

Options:
  -h --help         Show this screen
  -p --pretty       Pretty print (indent) json
  -n --newline      Newline after each document produced
`
	arguments, _ := docopt.Parse(usage, nil, true, "", false)
	pretty_print := arguments["--pretty"].(bool)
	print_newline := arguments["--newline"].(bool)
	filenames := arguments["FILE"].([]string)

	files := filenames[:]

	if len(files) == 0 {
		files = append(files, "-")
	}

	for _, fn := range files {
		var src *os.File
		var err error

		if fn == "-" {
			src = os.Stdin
		} else {
			src, err = os.Open(fn)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		reader := bufio.NewReader(src)
		buf, err := ioutil.ReadAll(reader)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		err = src.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			// but keep going...
		}

		var doc interface{}
		if err := yaml.Unmarshal(buf, &doc); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		doc = convert(doc)

		var js []byte
		if pretty_print {
			js, err = json.MarshalIndent(doc, "", "  ")
		} else {
			js, err = json.Marshal(doc)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("%s", js)
		if print_newline {
			fmt.Println()
		}
	}
}

func convert(i interface{}) interface{} {
	switch obj := i.(type) {
	case map[interface{}]interface{}:
		m := map[string]interface{}{}
		for k, v := range obj {
			m[k.(string)] = convert(v)
		}
		return m
	case []interface{}:
		for i, v := range obj {
			obj[i] = convert(v)
		}
	}
	return i
}
