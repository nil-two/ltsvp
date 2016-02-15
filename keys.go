package main

import (
	"regexp"
)

var (
	KEYS_LIST = regexp.MustCompile(`(?:[^,\\]|\\.)*`)
	BACKSLASH = regexp.MustCompile(`\\(.)`)
)

func ParseKeysList(list string) []string {
	keys := KEYS_LIST.FindAllString(list, -1)
	for i := 0; i < len(keys); i++ {
		keys[i] = BACKSLASH.ReplaceAllString(keys[i], "$1")
	}
	return keys
}
