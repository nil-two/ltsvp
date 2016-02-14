package main

import (
	"io"
	"strings"

	"github.com/ymotongpoo/goltsv"
)

type LTSVScanner struct {
	OutputDelimiter string
	RemainLTSV      bool
	keys            []string
	text            string
	err             error
	reader          *goltsv.LTSVReader
}

func NewLTSVScanner(keys []string, r io.Reader) *LTSVScanner {
	return &LTSVScanner{
		OutputDelimiter: "\t",
		keys:            keys,
		reader:          goltsv.NewReader(r),
	}
}

func (l *LTSVScanner) Err() error {
	if l.err == io.EOF {
		return nil
	}
	return l.err
}

func (l *LTSVScanner) Bytes() []byte {
	return []byte(l.text)
}

func (l *LTSVScanner) Text() string {
	return l.text
}

func (l *LTSVScanner) Scan() bool {
	if l.err != nil {
		return false
	}

	recode, err := l.reader.Read()
	if err != nil {
		l.err = err
		l.text = ""
		return false
	}

	switch {
	case l.RemainLTSV:
		var fields []string
		for _, key := range l.keys {
			fields = append(fields, key+":"+recode[key])
		}
		l.text = strings.Join(fields, "\t")
	default:
		var values []string
		for _, key := range l.keys {
			values = append(values, recode[key])
		}
		l.text = strings.Join(values, l.OutputDelimiter)
	}

	return true
}
