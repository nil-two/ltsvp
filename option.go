package main

type Option struct {
	List       string
	Delimiter  string
	RemainLTSV bool
	IsHelp     bool
	IsVersion  bool
	Files      []string
}
