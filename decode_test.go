package hl7

import (
	"os"
	"testing"
)

func TestDecode(t *testing.T) {
	fname := "./testdata/msg.hl7"
	file, err := os.Open(fname)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	ps := PIDSegment{}
	msgs, err := NewDecoder(file).Messages()
	if err != nil {
		t.Error(err)
	}
	if len(msgs) != 1 {
		t.Fatalf("Expected 1 message got %v\n", len(msgs))
	}
	if err := msgs[0].Unmarshal(&ps); err != nil {
		t.Fatal(err)
	}
	if ps.FirstName != "John" {
		t.Errorf("Expected John got %s\n", ps.FirstName)
	}
	if ps.LastName != "Jones" {
		t.Errorf("Expected Jones got %s\n", ps.LastName)
	}
}
