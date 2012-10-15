package main

import (
	"fmt"
	"github.com/pebbe/textcat"
	"github.com/pebbe/util"
	"io/ioutil"
	"os"
)

func main() {
	doUtf8 := true
	doRaw := true

	if isatty(os.Stdin) {
		syntax()
		return
	}

	for _, arg := range os.Args[1:] {
		switch arg {
		case "-r":
			doUtf8 = false
		case "-u":
			doRaw = false
		default:
			syntax()
			return
		}
	}
	if !doUtf8 && !doRaw {
		syntax()
		return
	}

	data, err := ioutil.ReadAll(os.Stdin)
	util.CheckErr(err)
	str := string(data)

	if doRaw {
		fmt.Println("[[[RAW]]]")
		n := 0
		for i, p := range textcat.GetPatterns(str, false) {
			if i == textcat.MaxPatterns {
				break
			}
			n += 1
			fmt.Printf("%s\t%d\n", p.S, p.I)
		}
		if n < textcat.MaxPatterns {
			fmt.Fprintf(os.Stderr, "Warning: there are less than %d raw patterns\n", textcat.MaxPatterns)
		}
	}

	if doUtf8 {
		fmt.Println("[[[UTF8]]]")
		n := 0
		for i, p := range textcat.GetPatterns(str, true) {
			if i == textcat.MaxPatterns {
				break
			}
			n += 1
			fmt.Printf("%s\t%d\n", p.S, p.I)
		}
		if n < textcat.MaxPatterns {
			fmt.Fprintf(os.Stderr, "Warning: there are less than %d utf8 patterns\n", textcat.MaxPatterns)
		}
	}

}

func syntax() {
	fmt.Fprintf(os.Stderr, `
Usage: %s [-r|-u] < sample data

Reads text samples from standard input, write to standard output
text patterns for package github.com/pebbe/textcat

Options:

    -r : raw patterns only
    -u : utf8 patterns only

`, os.Args[0])
}

// This seems to work on Windows and Linux, except on Linux when redirecting to/from device
func isatty(f *os.File) bool {
	s, e := f.Stat()
	if e != nil {
		return true
	}
	m := s.Mode()
	if m&os.ModeDevice != 0 {
		return true
	}
	return false
}
