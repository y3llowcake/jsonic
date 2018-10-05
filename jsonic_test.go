package jsonic

import (
	"testing"
)

func TestJsonic(t *testing.T) {
	j := MustNewString(`{"hello":"world"}`)
	_ = j
	if j.MustAt("hello").MustString() != "world" {
		t.Fatalf("expected world string")
	}
}
