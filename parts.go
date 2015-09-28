package main

import (
	"regexp"
)

var KEYS_LIST = regexp.MustCompile(`(?:[^,\\]|\\.)*`)

func ParseKeysList(list string) []string {
	return KEYS_LIST.FindAllString(list, -1)
}
