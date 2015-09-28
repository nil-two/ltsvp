package main

import (
	"reflect"
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
