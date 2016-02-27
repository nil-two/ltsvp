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
	version = "0.5.0"

	flagset         = pflag.NewFlagSet(name, pflag.ContinueOnError)
	list            = flagset.StringP("keys", "k", "", "")
	outputDelimiter = flagset.StringP("output-delimiter", "D", "\t", "")
	remainLTSV      = flagset.BoolP("remain-ltsv", "r", false, "")
	isHelp          = flagset.BoolP("help", "h", false, "")
	isVersion       = flagset.BoolP("version", "v", false, "")
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

func printErr(err interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
}

func guideToHelp() {
	fmt.Fprintf(os.Stderr, "Try '%s --help' for more information.\n", name)
}

func do(l *LTSVScanner) error {
	for l.Scan() {
		fmt.Println(l.Text())
	}
	return l.Err()
}

func main() {
	flagset.SetOutput(ioutil.Discard)
	if err := flagset.Parse(os.Args[1:]); err != nil {
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

	specifiedList := false
	flagset.Visit(func(f *pflag.Flag) {
		if f.Name == "keys" {
			specifiedList = true
		}
	})
	if !specifiedList {
		printErr("no specify LIST")
		guideToHelp()
		os.Exit(2)
	}
	keys := ParseKeysList(*list)

	r, err := argf.From(flagset.Args())
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
