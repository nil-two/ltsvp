package main

import (
	"io"
	"regexp"
	"strings"

	"github.com/ymotongpoo/goltsv"
)

var KEYS_LIST = regexp.MustCompile(`(?:[^,\\]|\\.)*`)

func ParseKeysList(list string) []string {
	keys := KEYS_LIST.FindAllString(list, -1)
	for i := 0; i < len(keys); i++ {
		keys[i] = strings.Replace(keys[i], `\,`, `,`, -1)
		keys[i] = strings.Replace(keys[i], `\\`, `\`, -1)
	}
	return keys
}

type LTSVScanner struct {
	keys   []string
	line   string
	err    error
	reader *goltsv.LTSVReader
}

func NewLTSVScanner(keys []string, r io.Reader) *LTSVScanner {
	return &LTSVScanner{
		keys:   keys,
		reader: goltsv.NewReader(r),
	}
}

func (l *LTSVScanner) Scan() bool {
	if l.err != nil {
		return false
	}

	recode, err := l.reader.Read()
	if err != nil {
		l.err = err
		return false
	}

	var values []string
	for _, key := range l.keys {
		values = append(values, recode[key])
	}
	l.line = strings.Join(values, "\t")

	return true
}

func (l *LTSVScanner) Text() string {
	return l.line
}
