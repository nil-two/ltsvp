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
		keys:   keys,
		reader: goltsv.NewReader(reader),
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

var MultipleKeysTests = []struct {
	keys []string
	src  string
	dst  []string
}{
	{
		keys: []string{"host"},
		src: `
host:192.168.0.1	status:200
host:172.16.0.12	status:404
`[1:],
		dst: []string{
			"192.168.0.1",
			"172.16.0.12",
		},
	},
	{
		keys: []string{"host", "status"},
		src: `
host:192.168.0.1	status:200
host:172.16.0.12	status:404
`[1:],
		dst: []string{
			"192.168.0.1\t200",
			"172.16.0.12\t404",
		},
	},
}

func TestMultipleKeys(t *testing.T) {
	for _, test := range MultipleKeysTests {
		l := NewLTSVScanner(test.keys,
			strings.NewReader(test.src))

		expect := test.dst
		actual := []string{}
		for l.Scan() {
			actual = append(actual, l.Text())
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("(keys: %q) got %q, want %q",
				test.keys, actual, expect)
		}
	}
}
