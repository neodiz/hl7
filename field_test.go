package hl7_test

import (
	"testing"

	"github.com/lenaten/hl7"
)

func TestFieldParse(t *testing.T) {
	val := []byte("520 51st Street^^Denver^CO^80020^USA")
	seps := hl7.NewDelimeters()
	fld := &hl7.Field{Value: val}
	fld.Parse(seps)
	if len(fld.Components) != 6 {
		t.Errorf("Expected 6 components got %d\n", len(fld.Components))
	}
}

func TestFieldSet(t *testing.T) {
	seps := hl7.NewDelimeters()
	fld := &hl7.Field{}
	loc := "ZZZ.1.10"
	l := hl7.NewLocation(loc)
	err := fld.Set(l, "TEST", seps)
	if err != nil {
		t.Error(err)
	}
	if len(fld.Components) != 11 {
		t.Fatalf("Expected 11 got %d\n", len(fld.Components))
	}
	if string(fld.Components[10].SubComponents[0].Value) != "TEST" {
		t.Errorf("Expected TEST got %s\n", fld.Components[10].SubComponents[0].Value)
	}
}
