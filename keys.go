package main

import (
	"regexp"
	"strings"
)

var (
	keysList  = regexp.MustCompile(`(?:[^,\\]|\\.)*`)
	backslash = regexp.MustCompile(`\\(.)`)
	trailing  = regexp.MustCompile(`\\+$`)
)

func ParseKeysList(list string) []string {
	list = trailing.ReplaceAllStringFunc(list, func(s string) string {
		return strings.Repeat(`\\`, len(s)/2)
	})
	if list == "" {
		return make([]string, 0)
	}

	keys := keysList.FindAllString(list, -1)
	for i := 0; i < len(keys); i++ {
		keys[i] = backslash.ReplaceAllString(keys[i], "$1")
	}
	return keys
}
