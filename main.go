package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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

type Option struct {
	List       string
	Delimiter  string
	RemainLTSV bool
	IsHelp     bool
	IsVersion  bool
	Files      []string
}

func ParseOption(args []string) (opt *Option, err error) {
	f := flag.NewFlagSet("ltsvp", flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)

	opt = &Option{}
	f.StringVar(&opt.List, "k", "", "")
	f.StringVar(&opt.List, "keys", "", "")
	f.StringVar(&opt.Delimiter, "d", "\t", "")
	f.StringVar(&opt.Delimiter, "delimiter", "\t", "")
	f.BoolVar(&opt.RemainLTSV, "r", false, "")
	f.BoolVar(&opt.RemainLTSV, "remain-ltsv", false, "")
	f.BoolVar(&opt.IsHelp, "h", false, "")
	f.BoolVar(&opt.IsHelp, "help", false, "")
	f.BoolVar(&opt.IsVersion, "v", false, "")
	f.BoolVar(&opt.IsVersion, "version", false, "")

	if err := f.Parse(args); err != nil {
		return nil, err
	}
	opt.Files = f.Args()
	return opt, nil
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
