package main

import (
	"regexp"
	"strings"
)

var (
	KEYS_LIST          = regexp.MustCompile(`(?:[^,\\]|\\.)*`)
	BACKSLASH          = regexp.MustCompile(`\\(.)`)
	TRAILING_BACKSLASH = regexp.MustCompile(`\\+$`)
)

func ParseKeysList(list string) []string {
	list = TRAILING_BACKSLASH.ReplaceAllStringFunc(list, func(s string) string {
		return strings.Repeat(`\\`, len(s)/2)
	})
	keys := KEYS_LIST.FindAllString(list, -1)
	for i := 0; i < len(keys); i++ {
		keys[i] = BACKSLASH.ReplaceAllString(keys[i], "$1")
	}
	return keys
}
