package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

func readFile(flag int) (*os.File, map[string]string) {
	var (
		file      *os.File
		dir       string
		reader    *csv.Reader
		recsSlice [][]string
		recs      map[string]string
		i         int
		ok        bool
		err       error
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

	file, err = os.OpenFile(dir+"/oishi.csv", flag|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	reader = csv.NewReader(file)
	reader.FieldsPerRecord = 2
	reader.Comment = '#'

	recsSlice, err = reader.ReadAll()
	if err != nil {
		panic(err)
	}

	recs = make(map[string]string, len(recsSlice))
	for i = range recsSlice {
		_, ok = recs[recsSlice[i][0]]
		if ok {
			panic(fmt.Errorf("key %q is repeated", recsSlice[i][0]))
		}

		recs[recsSlice[i][0]] = recsSlice[i][1]
	}

	return file, recs
}

func closeFile(file *os.File) {
	var err error

	err = file.Close()
	if err != nil {
		panic(err)
	}
}

func help() {
	fmt.Fprintf(os.Stderr, `usage: %s { help | copy <KEY> | add <KEY> <VALUE> }

%s is a simple clipboard program.

  help - display this help message
  copy - print the value of <KEY> to standard output
  add  - set the value of <KEY> to <VALUE>
`, os.Args[0], os.Args[0])
}

func die(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	help()
	os.Exit(1)
}

func main() {
	var (
		argc int
		f    *os.File
		recs map[string]string
		v    string
		ok   bool
	)

	argc = len(os.Args)

	if argc == 1 {
		die("%s: missing argument\n", os.Args[0])
	}

	switch os.Args[1] {
	case "help":
		help()
	case "copy":
		if argc != 3 {
			die("%s: expected 3 arguments, got %d\n", os.Args[0], argc)
		}

		f, recs = readFile(os.O_RDONLY)
		closeFile(f)

		v, ok = recs[os.Args[2]]
		if !ok {
			die("%s: key %q not found\n", os.Args[0], os.Args[2])
		}

		fmt.Println(v)
	}
}
