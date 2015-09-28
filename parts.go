package main

import (
	"regexp"
	"strings"
)

var KEYS_LIST = regexp.MustCompile(`(?:[^,\\]|\\.)*`)

func ParseKeysList(list string) []string {
	keys := KEYS_LIST.FindAllString(list, -1)
	for i := 0; i < len(keys); i++ {
		keys[i] = strings.Replace(keys[i], `\,`, `,`, -1)
	}
	return keys
}
