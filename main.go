package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
)

func help() {
	fmt.Fprintf(os.Stderr, `usage: %s { help | keys | copy <KEY> }

oishi is a simple clipboard program using CSV

  help - print this help message
  keys - print all keys to standard output, '\n' will be escaped
  copy - print the value of <KEY> to standard output

add key-value CSV records to $XDG_DATA_HOME/oishi/oishi.csv
add comments with '#'
`, os.Args[0])
}

func die(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	help()
	os.Exit(1)
}

func readFile() [][]string {
	var (
		file   *os.File
		dir    string
		reader *csv.Reader
		recs   [][]string
		dup    map[string]bool
		i      int
		err    error
	)

	dir = os.Getenv("XDG_DATA_HOME")
	if dir == "" {
		dir = os.Getenv("HOME")
		if dir == "" {
			panic(errors.New("$HOME not set"))
		}

		dir += "/.local/share"
	}

	dir += "/oishi"
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}

	file, err = os.OpenFile(dir+"/oishi.csv", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	reader = csv.NewReader(file)
	reader.FieldsPerRecord = 2
	reader.Comment = '#'

	recs, err = reader.ReadAll()
	if err != nil {
		panic(err)
	}

	if len(recs) == 0 {
		die("%s: oishi.csv is empty\n", os.Args[0])
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}

	dup = make(map[string]bool, len(recs))
	for i = range recs {
		if dup[recs[i][0]] {
			panic(fmt.Errorf("key %q is repeated", recs[i][0]))
		}

		dup[recs[i][0]] = true
	}

	return recs
}

func main() {
	var (
		argc    int
		recs    [][]string
		builder strings.Builder
		i       int
	)

	argc = len(os.Args)

	if argc == 1 {
		die("%s: missing argument\n", os.Args[0])
	}

	switch os.Args[1] {
	case "help":
		help()
	case "keys":
		recs = readFile()

		for i = range recs {
			builder.WriteString(strings.ReplaceAll(recs[i][0], "\n", "\\n"))
			builder.WriteByte('\n')
		}

		fmt.Print(builder.String())
	case "copy":
		if argc < 3 {
			die("%s: missing <KEY> argument\n", os.Args[0])
		}

		recs = readFile()
		for i = range recs {
			if recs[i][0] == os.Args[2] {
				fmt.Println(recs[i][1])
				os.Exit(0)
			}
		}

		die("%s: key %q not found\n", os.Args[0], os.Args[2])
	default:
		die("%s: unknown argument %q\n", os.Args[0], os.Args[1])
	}
}
