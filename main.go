package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/yuya-takeyama/argf"
)

func usage() {
	os.Stderr.WriteString(`
Usage: ltsvp OPTION... [FILE]...
Print selected parts of LTSV from each FILE to standard output.

Options:
  -k, --keys=LIST          select only these keys (required)
  -d, --delimiter=STRING   use STRING to separate parts (default: \t)
  -r, --remain-ltsv        print selected parts as LTSV
      --help               display this help text and exit
      --version            output version information and exit

LIST is made up of keys separated by commas.
  host           # Select host
  host,time,ua   # Select host, time, and ua
`[1:])
}

func version() {
	os.Stderr.WriteString(`
0.2.0
`[1:])
}

type Option struct {
	List       string `short:"k" long:"keys" required:"true"`
	Delimiter  string `short:"d" long:"delimiter" default:"\t"`
	RemainLTSV bool   `short:"r" long:"remain-ltsv"`
	IsHelp     bool   `          long:"help"`
	IsVersion  bool   `          long:"version"`
	Files      []string
}

func parseOption(args []string) (opt *Option, err error) {
	opt = &Option{}
	flag := flags.NewParser(opt, flags.PassDoubleDash)

	opt.Files, err = flag.ParseArgs(args)
	if err != nil && !opt.IsHelp && !opt.IsVersion {
		return nil, err
	}
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

func printErr(err error) {
	fmt.Fprintln(os.Stderr, "ltsvp:", err)
}

func guideToHelp() {
	os.Stderr.WriteString(`
Try 'ltsvp --help' for more information.
`[1:])
}

func _main() int {
	opt, err := parseOption(os.Args[1:])
	if err != nil {
		printErr(err)
		guideToHelp()
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
		guideToHelp()
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
