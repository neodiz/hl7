package hl7

import (
	"os"
	"testing"
)

type MyHL7Message struct {
	SendingApp        string `hl7:"MSH.3"`
	SendingFacility   string `hl7:"MSH.4"`
	ReceivingApp      string `hl7:"MSH.5"`
	ReceivingFacility string `hl7:"MSH.6"`
	MsgDate           string `hl7:"MSH.7"`
	MessageType       string `hl7:"MSH.9"`
	ControlID         string `hl7:"MSH.10"`
	ProcessingID      string `hl7:"MSH.11"`
	VersionID         string `hl7:"MSH.12"`
}

type PIDSegment struct {
	FirstName string `hl7:"PID.5.1"`
	LastName  string `hl7:"PID.5.0"`
}

func TestCreateHL7WithMsgInfo(t *testing.T) {
	mi := MsgInfo{
		SendingApp:        "MyApp",
		SendingFacility:   "MyPlace",
		ReceivingApp:      "EMR",
		ReceivingFacility: "MedicalPlace",
		MessageType:       "ORM^001",
	}

	msg, err := StartMessage(mi)
	if err != nil {
		t.Error("Error creating initial message:", err.Error())
		return
	}

	ps := PIDSegment{FirstName: "Davin", LastName: "Hills"}
	bytes, err := Marshal(msg, &ps)
	if err != nil {
		t.Error("Error marshaling message to bytes:", err.Error())
		return
	}

	err = os.WriteFile("test1.hl7", bytes, os.ModeAppend)
	if err != nil {
		t.Error("Error writing to file:", err.Error())
		return
	}
}

func TestCreateHL7WithCustomStruct(t *testing.T) {
	my := MyHL7Message{
		SendingApp:        "MyApp",
		SendingFacility:   "MyPlace",
		ReceivingApp:      "EMR",
		ReceivingFacility: "MedicalPlace",
		MessageType:       "ORM^001",
		MsgDate:           "20151209154606",
		ControlID:         "MSGID1",
		ProcessingID:      "P",
		VersionID:         "2.4",
	}

	ps := PIDSegment{FirstName: "Davin", LastName: "Hills"}

	file, err := os.OpenFile("test2.hl7", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Error("Error opening file:", err.Error())
		return
	}

	err = NewEncoder(file).Encode(&my)
	if err != nil {
		t.Error("Error encoding initial message:", err.Error())
		return
	}

	err = NewEncoder(file).Encode(&ps)
	if err != nil {
		t.Error("Error encoding PID segment:", err.Error())
		return
	}
}
