package main

import (
	"fmt"
	"os"

	"github.com/ogier/pflag"
	"github.com/yuya-takeyama/argf"
)

var (
	flag            = pflag.NewFlagSet("ltsvp", pflag.ContinueOnError)
	list            = pflag.StringP("keys", "k", "", "")
	outputDelimiter = pflag.StringP("output-delimiter", "d", "\t", "")
	remainLTSV      = pflag.BoolP("remain-ltsv", "r", false, "")
	isHelp          = pflag.BoolP("help", "h", false, "")
	isVersion       = pflag.BoolP("version", "v", false, "")
)

func printUsage() {
	os.Stderr.WriteString(`
Usage: ltsvp OPTION... [FILE]...
Print selected parts of LTSV from each FILE to standard output.

Options:
  -k, --keys=LIST
                 select only these keys (required)
  -D, --output-delimiter=STRING
                 use STRING to separate parts (default: \t)
  -r, --remain-ltsv
                 print selected parts as LTSV
  --help
                 display this help text and exit
  --version
                 output version information and exit

LIST is made up of keys separated by commas.
  host           # Select host
  host,time,ua   # Select host, time, and ua
`[1:])
}

func printVersion() {
	os.Stderr.WriteString(`
0.3.0
`[1:])
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
	if err := flag.Parse(os.Args[1:]); err != nil {
		printErr(err)
		guideToHelp()
		return 2
	}
	switch {
	case *isHelp:
		printUsage()
		return 0
	case *isVersion:
		printVersion()
		return 0
	}

	keys := ParseKeysList(*list)
	r, err := argf.From(flag.Args())
	if err != nil {
		printErr(err)
		return 2
	}

	l := NewLTSVScanner(keys, r)
	l.OutputDelimiter = *outputDelimiter
	l.RemainLTSV = *remainLTSV
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
