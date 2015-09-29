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

func TestScan(t *testing.T) {
	keys := []string{"host"}
	reader := strings.NewReader(`
host:192.168.0.1	status:200
host:172.16.0.12	status:404
`[1:])
	l := NewLTSVScanner(keys, reader)
	expects := []bool{true, true, false}
	for i := 0; i < len(expects); i++ {
		expect := expects[i]
		actual := l.Scan()
		if actual != expect {
			t.Errorf("Scan[%v]: got %v, want %v", i, actual, expect)
		}
	}
}
