package xaddr

import "testing"

func TestParseQUICK(t *testing.T) {
	_, e := ParseOFFLINE("google.hotels2")
	if e == nil {
		t.Error("ParseOFFLINE test failed")
	}
}
