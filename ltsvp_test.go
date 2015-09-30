package main

import (
	"reflect"
	"strings"
	"testing"

	"github.com/ymotongpoo/goltsv"
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

type ScanResult struct {
	scan bool
	text string
	err  error
}

var ScanTests = []struct {
	description string
	src         string
	keys        []string
	result      []ScanResult
}{
	{
		description: "regular LTSV",
		keys:        []string{"host"},
		src: `
host:192.168.0.1	status:200
host:172.16.0.12	status:404
`[1:],
		result: []ScanResult{
			{scan: true, text: "192.168.0.1", err: nil},
			{scan: true, text: "172.16.0.12", err: nil},
			{scan: false, text: "", err: nil},
		},
	},
	{
		description: "invalid LTSV",
		keys:        []string{"host"},
		src: `,
host:192.168.0.1	status:200
a	b	c
host:172.16.0.12	status:404
`[1:],
		result: []ScanResult{
			{scan: true, text: "192.168.0.1", err: nil},
			{scan: false, text: "", err: goltsv.ErrLabelName},
			{scan: false, text: "", err: goltsv.ErrLabelName},
		},
	},
}

func TestScan(t *testing.T) {
	for _, test := range ScanTests {
		reader := strings.NewReader(test.src)
		l := NewLTSVScanner(test.keys, reader)
		for i := 0; i < len(test.result); i++ {
			scan := l.Scan()
			expect := test.result[i]
			actual := ScanResult{
				scan: scan,
				text: l.Text(),
				err:  l.Err(),
			}
			if !reflect.DeepEqual(actual, expect) {
				t.Errorf("%s: %v: got %v, want %v",
					test.description, i+1,
					actual, expect)
			}
		}
	}
}

var DelimiterTests = []struct {
	description string
	keys        []string
	delimiter   string
	src         string
	dst         []string
}{
	{
		description: "with comma",
		keys:        []string{"host", "status"},
		delimiter:   ",",
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
		description: "with double dash",
		keys:        []string{"host", "status"},
		delimiter:   "--",
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
			t.Errorf("%s: got %q, want %q",
				test.description, actual, expect)
		}
	}
}

var RemainLTSVTests = []struct {
	description string
	keys        []string
	delimiter   string
	src         string
	dst         []string
}{
	{
		description: "one key",
		keys:        []string{"host"},
		delimiter:   "\t",
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
		description: "two keys",
		keys:        []string{"status", "host"},
		delimiter:   "\t",
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
		description: "ignore delimiter",
		keys:        []string{"status", "host"},
		delimiter:   "---",
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
			t.Errorf("%s: got %q, want %q",
				test.description, actual, expect)
		}
	}
}

func TestBytes(t *testing.T) {
	keys := []string{"host"}
	reader := strings.NewReader(`
host:192.168.0.1	status:200
host:172.16.0.12	status:404
`[1:])
	l := NewLTSVScanner(keys, reader)

	expects := [][]byte{
		[]byte("192.168.0.1"),
		[]byte("172.16.0.12"),
	}
	for i := 0; i < len(expects); i++ {
		l.Scan()
		expect := expects[i]
		actual := l.Bytes()
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("%v: got %v, want %v",
				i+1, actual, expect)
		}
	}
}

func BenchmarkNew(b *testing.B) {
	src := strings.Repeat("key1:value\tkey2:value2\tkey3:value3\n", 10000)
	for i := 0; i < b.N; i++ {
		keys := []string{"key2", "key3"}
		reader := strings.NewReader(src)
		NewLTSVScanner(keys, reader)
	}
}

func BenchmarkScan(b *testing.B) {
	src := strings.Repeat("key1:value\tkey2:value2\tkey3:value3\n", 10000)
	for i := 0; i < b.N; i++ {
		keys := []string{"key2", "key3"}
		reader := strings.NewReader(src)
		l := NewLTSVScanner(keys, reader)
		for l.Scan() {
		}
	}
}
