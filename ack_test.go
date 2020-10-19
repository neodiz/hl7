package hl7

import (
	"errors"
	"os"
	"testing"
)

func TestAcknowledge(t *testing.T) {
	fname := "./testdata/msg.hl7"
	file, err := os.Open(fname)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	msgs, err := NewDecoder(file).Messages()
	if err != nil {
		t.Fatal(err)
	}
	mi, err := msgs[0].Info()
	ack, _ := Acknowledge(mi, nil)
	if ack == nil {
		t.Fatal("Expected ACK message got nil")
	}
	// for _, s := range ack.Segments {
	// 	for _, f := range s.Fields {
	// 		fmt.Println(string(f.Value))
	// 	}
	// }
	ack, _ = Acknowledge(mi, errors.New("This is a test error"))
	if ack == nil {
		t.Fatal("Expected ACK message got nil")
	}
	m := NewMsgInfo()
	m.ReceivingApp = "ORG_REC_APP"
	m.ReceivingFacility = "ORG_REC_FAC"
	m.SendingApp = "ORG_SEND_APP"
	m.SendingFacility = "ORG_SEND_FAC"
	ack, _ = Acknowledge(*m, errors.New("Fatal error"))
	if ack == nil {
		t.Fatal("Expected ACK message got nil")
	}
}
