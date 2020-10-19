package hl7_test

import (
	"os"
	"testing"

	"github.com/freemed/hl7"
)

func TestDecode(t *testing.T) {
	fname := "./testdata/msg.hl7"
	file, err := os.Open(fname)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	st := my7{}
	msgs, err := hl7.NewDecoder(file).Messages()
	if err != nil {
		t.Error(err)
	}
	if len(msgs) != 1 {
		t.Fatalf("Expected 1 message got %v\n", len(msgs))
	}
	if err := msgs[0].Unmarshal(&st); err != nil {
		t.Fatal(err)
	}
	if st.FirstName != "John" {
		t.Errorf("Expected John got %s\n", st.FirstName)
	}
	if st.LastName != "Jones" {
		t.Errorf("Expected Jones got %s\n", st.LastName)
	}
}
