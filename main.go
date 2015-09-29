package main

import (
	"fmt"
	"os"

	"github.com/yuya-takeyama/argf"
)

func usage() {
	os.Stderr.WriteString(`
Usage: ltsvp OPTION... [FILE]...
Print selected parts of ltsv from each FILE to standard output.

Options:
  -k, --keys=LIST          select only these keys
                           LIST is made up of keys separated by commas
                             host           # Select host
                             host,time,ua   # Select host, time, and ua
  -d, --delimiter=STRING   use STRING to separate parts (default: \t)
  -r, --remain-ltsv        print selected parts as LTSV
  -h, --help               display this help text and exit
  -v, --version            display version information and exit
`[1:])
}

func version() {
	os.Stderr.WriteString(`
v0.1.0
`[1:])
}

func printErr(err error) {
	fmt.Fprintln(os.Stderr, "ltsvp:", err)
}

func newLTSVScannerFromOption(opt *Option) (l *LTSVScanner, err error) {
	keys := ParseKeysList(opt.List)
	reader, err := argf.From(opt.Files)
	if err != nil {
		return nil, err
	}

	l = NewLTSVScanner(keys, reader)
	l.Delimiter = opt.Delimiter
	l.RemainLTSV = opt.RemainLTSV
	return l, nil
}

func do(l *LTSVScanner) error {
	for l.Scan() {
		fmt.Println(l.Text())
	}
	return l.Err()
}

func _main() int {
	opt, err := ParseOption(os.Args[1:])
	if err != nil {
		printErr(err)
		return 2
	}
	switch {
	case opt.IsHelp:
		usage()
		return 0
	case opt.IsVersion:
		version()
		return 0
	}

	l, err := newLTSVScannerFromOption(opt)
	if err != nil {
		printErr(err)
		return 2
	}
	if err := do(l); err != nil {
		printErr(err)
		return 1
	}
	return 0
}

func main() {
	e := _main()
	os.Exit(e)
}
