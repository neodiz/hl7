package hl7_test

import (
	"os"
	"testing"

	"github.com/lenaten/hl7"
)

func TestValid(t *testing.T) {

	fname := "./testdata/msg.hl7"
	file, err := os.Open(fname)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	msgs, err := hl7.NewDecoder(file).Messages()
	if err != nil {
		t.Fatal(err)
	}

	valid, failures := msgs[0].IsValid(hl7.NewValidMSH24())
	if valid == false {
		t.Error("Expected valid MSH got invalid. Failures:")
		for i, f := range failures {
			t.Errorf("%d %+v\n", i, f)
		}
	}

	valid, failures = msgs[0].IsValid(hl7.NewValidPID24())
	if valid == false {
		t.Error("Expected valid PID got invalid. Failures:")
		for i, f := range failures {
			t.Errorf("%d %+v\n", i, f)
		}
	}
}
