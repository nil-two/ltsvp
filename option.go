package main

import (
	"flag"
	"io/ioutil"
)

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
