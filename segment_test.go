package hl7_test

import (
	"testing"

	"github.com/lenaten/hl7"
)

func TestSegParse(t *testing.T) {
	val := []byte("PID|||12001||Jones^John^^^Mr.||19670824|M|||123 West St.^^Denver^CO^80020^USA~520 51st Street^^Denver^CO^80020^USA|||||||")
	seps := hl7.NewDelimeters()
	seg := &hl7.Segment{Value: val}
	seg.Parse(seps)
	if len(seg.Fields) != 20 {
		t.Errorf("Expected 20 fields got %d\n", len(seg.Fields))
	}
}

func TestSegSet(t *testing.T) {
	seps := hl7.NewDelimeters()
	loc := "ZZZ.10"
	l := hl7.NewLocation(loc)
	seg := &hl7.Segment{}
	err := seg.Set(l, "TEST", seps)
	if err != nil {
		t.Error(seg)
	}
	str, err := seg.Get(l)
	if err != nil {
		t.Error(err)
	}
	if str != "TEST" {
		t.Errorf("Expected TEST got %s\n", str)
	}
}
