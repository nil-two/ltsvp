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
	Delimiter string
	keys      []string
	line      string
	err       error
	reader    *goltsv.LTSVReader
}

func NewLTSVScanner(keys []string, r io.Reader) *LTSVScanner {
	return &LTSVScanner{
		Delimiter: "\t",
		keys:      keys,
		reader:    goltsv.NewReader(r),
	}
}

func (l *LTSVScanner) Scan() bool {
	if l.err != nil {
		return false
	}

	recode, err := l.reader.Read()
	if err != nil {
		l.err = err
		l.line = ""
		return false
	}

	var values []string
	for _, key := range l.keys {
		values = append(values, recode[key])
	}
	l.line = strings.Join(values, l.Delimiter)

	return true
}

func (l *LTSVScanner) Err() error {
	if l.err == io.EOF {
		return nil
	}
	return l.err
}

func (l *LTSVScanner) Text() string {
	return l.line
}
