package main

import (
	"github.com/ymotongpoo/goltsv"
	"reflect"
	"strings"
	"testing"
)

var ParseKeysListTests = []struct {
	list string
	keys []string
}{
	// normal
	{`host`, []string{`host`}},
	{`host,status`, []string{`host`, `status`}},
	{`host,status,size`, []string{`host`, `status`, `size`}},

	// include empty keys
	{``, []string{``}},
	{`,`, []string{``, ``}},
	{`,,`, []string{``, ``, ``}},
	{`,host`, []string{``, `host`}},
	{`,,host`, []string{``, ``, `host`}},
	{`host,`, []string{`host`, ``}},
	{`host,,`, []string{`host`, ``, ``}},
	{`,,host,,status,,`, []string{``, ``, `host`, ``, `status`, ``, ``}},

	// include escaped comma
	{`a\,b`, []string{`a,b`}},
	{`a\,\,b`, []string{`a,,b`}},
	{`a\,,b\,`, []string{`a,`, `b,`}},
	{`\,a,\,b`, []string{`,a`, `,b`}},
	{`\,a\,,\,b\,`, []string{`,a,`, `,b,`}},
	{`a\,b,c\,d\,e`, []string{`a,b`, `c,d,e`}},
	{`a\,b,c\,d\,e,f\,g\,h\,i`, []string{`a,b`, `c,d,e`, `f,g,h,i`}},

	// include escaped backslash
	{`a\\b`, []string{`a\b`}},
	{`a\\\\b`, []string{`a\\b`}},
	{`a\\,b\\`, []string{`a\`, `b\`}},
	{`\\a,\\b`, []string{`\a`, `\b`}},
	{`\\a\\,\\b\\`, []string{`\a\`, `\b\`}},
	{`a\\b,c\\d\\e`, []string{`a\b`, `c\d\e`}},
	{`a\\b,c\\d\\e,f\\g\\h\\i`, []string{`a\b`, `c\d\e`, `f\g\h\i`}},
}

func TestParseKeysList(t *testing.T) {
	for _, test := range ParseKeysListTests {
		expect := test.keys
		actual := ParseKeysList(test.list)
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("ParseKeysList(%q) = %q, want %q",
				test.list, actual, expect)
		}
	}
}

func TestNewLTSVScanner(t *testing.T) {
	keys := []string{"host", "time"}
	reader := strings.NewReader(``)

	expect := &LTSVScanner{
		Delimiter:  "\t",
		RemainLTSV: false,
		keys:       keys,
		reader:     goltsv.NewReader(reader),
	}
	actual := NewLTSVScanner(keys, reader)
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("NewLTSVScanner(%q, %q) = %q, want %q",
			keys, reader, actual, expect)
	}
}

type ScanResult struct {
	scan bool
	text string
	err  error
}

func TestScan(t *testing.T) {
	keys := []string{"host"}
	reader := strings.NewReader(`
host:192.168.0.1	status:200
host:172.16.0.12	status:404
`[1:])
	l := NewLTSVScanner(keys, reader)

	expects := []ScanResult{
		{scan: true, text: "192.168.0.1", err: nil},
		{scan: true, text: "172.16.0.12", err: nil},
		{scan: false, text: "", err: nil},
	}
	for i := 0; i < len(expects); i++ {
		expect := expects[i]
		actual := ScanResult{}
		actual.scan = l.Scan()
		actual.text = l.Text()
		actual.err = l.Err()
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("Scan %v time: got %v, want %v",
				i+1, actual, expect)
		}
	}
}

func TestScanError(t *testing.T) {
	keys := []string{"host"}
	reader := strings.NewReader(`
host:192.168.0.1	status:200
a	b	c
host:172.16.0.12	status:404
`[1:])
	l := NewLTSVScanner(keys, reader)

	expects := []ScanResult{
		{scan: true, text: "192.168.0.1", err: nil},
		{scan: false, text: "", err: goltsv.ErrLabelName},
		{scan: false, text: "", err: goltsv.ErrLabelName},
	}
	for i := 0; i < len(expects); i++ {
		expect := expects[i]
		actual := ScanResult{}
		actual.scan = l.Scan()
		actual.text = l.Text()
		actual.err = l.Err()
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("Scan %v time: got %v, want %v",
				i+1, actual, expect)
		}
	}
}

var DelimiterTests = []struct {
	keys      []string
	delimiter string
	src       string
	dst       []string
}{
	{
		keys:      []string{"host", "status"},
		delimiter: ",",
		src: `
host:192.168.0.1	status:200
host:172.16.0.12	status:404
`[1:],
		dst: []string{
			"192.168.0.1,200",
			"172.16.0.12,404",
		},
	},
	{
		keys:      []string{"host", "status"},
		delimiter: "--",
		src: `
host:192.168.0.1	status:200
host:172.16.0.12	status:404
`[1:],
		dst: []string{
			"192.168.0.1--200",
			"172.16.0.12--404",
		},
	},
}

func TestDelimiter(t *testing.T) {
	for _, test := range DelimiterTests {
		l := NewLTSVScanner(test.keys, strings.NewReader(test.src))
		l.Delimiter = test.delimiter

		expect := test.dst
		actual := []string{}
		for l.Scan() {
			actual = append(actual, l.Text())
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("(keys: %q, delimiter: %q) got %q, want %q",
				test.keys, test.delimiter, actual, expect)
		}
	}
}

var RemainLTSVTests = []struct {
	keys      []string
	delimiter string
	src       string
	dst       []string
}{
	{
		keys:      []string{"host"},
		delimiter: "\t",
		src: `
host:192.168.0.1	status:200
host:172.16.0.12	status:404
`[1:],
		dst: []string{
			"host:192.168.0.1",
			"host:172.16.0.12",
		},
	},
	{
		keys:      []string{"status", "host"},
		delimiter: "\t",
		src: `
host:192.168.0.1	status:200
host:172.16.0.12	status:404
`[1:],
		dst: []string{
			"status:200\thost:192.168.0.1",
			"status:404\thost:172.16.0.12",
		},
	},
	{
		keys:      []string{"status", "host"},
		delimiter: "---",
		src: `
host:192.168.0.1	status:200
host:172.16.0.12	status:404
`[1:],
		dst: []string{
			"status:200\thost:192.168.0.1",
			"status:404\thost:172.16.0.12",
		},
	},
}

func TestRemainLTSV(t *testing.T) {
	for _, test := range RemainLTSVTests {
		l := NewLTSVScanner(test.keys, strings.NewReader(test.src))
		l.RemainLTSV = true

		expect := test.dst
		actual := []string{}
		for l.Scan() {
			actual = append(actual, l.Text())
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("(keys: %q, remainLTSV: true) got %q, want %q",
				test.keys, actual, expect)
		}
	}
}
