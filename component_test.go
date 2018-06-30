package hl7_test

import (
	"testing"

	"github.com/lenaten/hl7"
)

func TestCompParse(t *testing.T) {
	val := []byte("v1&v2&v3&&v5")
	seps := hl7.NewDelimeters()
	cmp := &hl7.Component{Value: val}
	cmp.Parse(seps)
	if len(cmp.SubComponents) != 5 {
		t.Errorf("Expected 5 subcomponents got %d\n", len(cmp.SubComponents))
	}
}

func TestCompSet(t *testing.T) {
	seps := hl7.NewDelimeters()
	loc := "ZZZ.1.0.5"
	l := hl7.NewLocation(loc)
	cmp := &hl7.Component{}
	err := cmp.Set(l, "TEST", seps)
	if err != nil {
		t.Error(err)
	}
	if len(cmp.SubComponents) != 6 {
		t.Fatalf("Expected 6 got %d\n", len(cmp.SubComponents))
	}
	if string(cmp.SubComponents[5].Value) != "TEST" {
		t.Errorf("Expected TEST got %s\n", cmp.SubComponents[5].Value)
	}
}
