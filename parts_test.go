package main

import (
	"reflect"
	"testing"
)

func TestParseKeysList(t *testing.T) {
	list := "host,status"
	expect := []string{"host", "status"}
	actual := ParseKeysList(list)
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("ParseKeysList(%q) = %q, want %q",
			list, actual, expect)
	}
}
