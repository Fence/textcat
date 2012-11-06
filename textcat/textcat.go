/*
The program `textcat` is for classifying text by language.

Usage:

    textcat [-f=textfile] [-b|-r] [-p=patternfiles] [-a] [-l] [text]

The text to be classified is the first applicable of these:
1) text from a file, loaded with option: -f=filename;
2) text on the command line, following any options;
3) text read from standard input.

By default, only utf-8 patterns are used. Options to change this are:

    -b : both raw and utf-8 patterns
    -r : raw patterns, instead of utf-8

You can load additional language patterns with option -p:

    -p=language1,language2

Here, both `language1` and `language2` are pattern files create with the
`textpat` program. Note: pattern files are listed with commas in
between, and no spaces.

By default, `textcat` classifies the whole input document as a single
text. To classify individual lines instead, use option -l

Use option -a to get a list of all available languages (after processing
options -b, -r and -p).
*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/pebbe/textcat"
	"github.com/pebbe/util"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var (
	opt_a = flag.Bool("a", false, "list all languages, and exit")
	opt_b = flag.Bool("b", false, "both raw and utf-8 patterns")
	opt_r = flag.Bool("r", false, "raw patterns, instead of utf-8")
	opt_l = flag.Bool("l", false, "classify individual lines instead of whole document")
	opt_f = flag.String("f", "", "file name")
	opt_p = flag.String("p", "", "pattern file names, separated by comma's (no spaces)")
)

func main() {
	flag.Parse()

	if *opt_f == "" && flag.NArg() == 0 && util.IsTerminal(os.Stdin) && !*opt_a {
		fmt.Fprintf(os.Stderr, "\nUsage: %s [args] [text]\n\nargs with default values are:\n\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nIf both -f and text are missing, read from stdin\n\n")
		return
	}

	tc := textcat.NewTextCat()
	if *opt_p != "" {
		for _, i := range strings.Split(*opt_p, ",") {
			e := tc.AddLanguage(path.Base(i), i)
			util.CheckErr(e)
		}
	}
	if *opt_r || *opt_b {
		tc.EnableAllRawLanguages()
	}
	if *opt_b || !*opt_r {
		tc.EnableAllUtf8Languages()
	}

	if *opt_a {
		for _, i := range tc.ActiveLanguages() {
			fmt.Println(i)
		}
		return
	}

	if *opt_l {
		var r *util.Reader
		if *opt_f != "" {
			fp, err := os.Open(*opt_f)
			util.CheckErr(err)
			defer fp.Close()
			r = util.NewReader(fp)
		} else if flag.NArg() > 0 {
			b := bytes.NewBufferString(strings.Join(flag.Args(), " "))
			r = util.NewReader(b)
		} else {
			r = util.NewReader(os.Stdin)
		}
		for {
			line, err := r.ReadLineString()
			if err == io.EOF {
				break
			}
			util.CheckErr(err)
			l, err := tc.Classify(line)
			if err != nil {
				fmt.Print(err)
			} else {
				fmt.Print(strings.Join(l, ","))
			}
			fmt.Println("\t" + line)
		}
		return
	}

	var text string
	if *opt_f != "" {
		t, err := ioutil.ReadFile(*opt_f)
		util.CheckErr(err)
		text = string(t)
	} else if flag.NArg() > 0 {
		text = strings.Join(flag.Args(), " ")
	} else {
		t, err := ioutil.ReadAll(os.Stdin)
		util.CheckErr(err)
		text = string(t)
	}

	l, e := tc.Classify(text)
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(strings.Join(l, "\n"))
	}
}