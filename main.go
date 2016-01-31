package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ogier/pflag"
	"github.com/yuya-takeyama/argf"
)

var (
	name    = "ltsvp"
	version = "0.4.2"

	flag            = pflag.NewFlagSet(name, pflag.ContinueOnError)
	list            = flag.StringP("keys", "k", "", "")
	outputDelimiter = flag.StringP("output-delimiter", "D", "\t", "")
	remainLTSV      = flag.BoolP("remain-ltsv", "r", false, "")
	isHelp          = flag.BoolP("help", "h", false, "")
	isVersion       = flag.BoolP("version", "v", false, "")
)

func printUsage() {
	fmt.Fprintf(os.Stderr, `
Usage: %s OPTION... [FILE]...
Print selected parts of LTSV from each FILE to standard output.

Options:
  -k, --keys=LIST
                 select only these keys (required)
  -D, --output-delimiter=STRING
                 use STRING to separate parts (default: \t)
  -r, --remain-ltsv
                 print selected parts as LTSV
  -h, --help
                 display this help text and exit
  -v, --version
                 output version information and exit

LIST is made up of keys separated by commas.
  host           # Select host
  host,time,ua   # Select host, time, and ua
`[1:], name)
}

func printVersion() {
	fmt.Fprintln(os.Stderr, version)
}

func do(l *LTSVScanner) error {
	for l.Scan() {
		fmt.Println(l.Text())
	}
	return l.Err()
}

func printErr(err interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
}

func guideToHelp() {
	fmt.Fprintf(os.Stderr, "Try '%s --help' for more information.\n", name)
}

func main() {
	flag.SetOutput(ioutil.Discard)
	if err := flag.Parse(os.Args[1:]); err != nil {
		printErr(err)
		guideToHelp()
		os.Exit(2)
	}
	switch {
	case *isHelp:
		printUsage()
		os.Exit(0)
	case *isVersion:
		printVersion()
		os.Exit(0)
	}

	if *list == "" {
		printErr("no specify LIST")
		guideToHelp()
		os.Exit(2)
	}
	keys := ParseKeysList(*list)

	r, err := argf.From(flag.Args())
	if err != nil {
		printErr(err)
		guideToHelp()
		os.Exit(2)
	}

	l := NewLTSVScanner(keys, r)
	l.OutputDelimiter = *outputDelimiter
	l.RemainLTSV = *remainLTSV
	if err := do(l); err != nil {
		printErr(err)
		os.Exit(1)
	}
	os.Exit(0)
}
