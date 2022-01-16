package hl7

import (
	"os"
	"testing"
)

type MyHL7Message struct {
	FieldSeparator     string `hl7:"MSH.1" hl7default:"^~\\&"`
	EncodingCharacters string `hl7:"MSH.2" hl7default:""`
	SendingApp         string `hl7:"MSH.3"`
	SendingFacility    string `hl7:"MSH.4"`
	ReceivingApp       string `hl7:"MSH.5"`
	ReceivingFacility  string `hl7:"MSH.6"`
	MsgDate            string `hl7:"MSH.7"`
	MessageType        string `hl7:"MSH.9"`
	ControlID          string `hl7:"MSH.10"`
	ProcessingID       string `hl7:"MSH.11"`
	VersionID          string `hl7:"MSH.12" hl7default:"2.4"`
}

type PIDSegment struct {
	FirstName string `hl7:"PID.5.1"`
	LastName  string `hl7:"PID.5.0"`
}

type PV1Segment struct {
	SetID                   string `hl7:"PV1.1"`
	PatientClass            string `hl7:"PV1.2"`
	AssignedPatientLocation string `hl7:"PV1.3"`
	AdmissionType           string `hl7:"PV1.4"`
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

	pid := PIDSegment{FirstName: "Davin", LastName: "Hills"}
	pv1 := PV1Segment{SetID: "0001", PatientClass: "I", AdmissionType: "O"}

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

	err = NewEncoder(file).Encode(&pid)
	if err != nil {
		t.Error("Error encoding PID segment:", err.Error())
		return
	}

	err = NewEncoder(file).Encode(&pv1)
	if err != nil {
		t.Error("Error encoding PID segment:", err.Error())
		return
	}
}
