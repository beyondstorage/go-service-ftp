package ftp

import (
	"regexp"
	"testing"
)

func TestConnection(t *testing.T) {
	name := "hello"
	want := regexp.MustCompile(`\b` + name + `\b`)
	msg := "Hello"
	if !want.MatchString(msg) {
		t.Fatalf("Wrong")
	}
}
