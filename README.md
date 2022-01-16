# HL7 [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/borisrodman/hl7)

## Overview

	This is a fork of Davin Hills' Go Level 7 library for creation/manipulation of HL7 files.

> 	HL7 segment specs can be found here: https://hl7-definition.caristix.com/v2/HL7v2.4/Segments

## Features

* Decode HL7 messages
* Multiple message support with Split
* Unmarshal into Go structs
* Simple query syntax
* Message validation

Note: Message building is not currently working for MSH segments. Coming soon...

## Installation
	go get github.com/borisrodman/hl7

## Usage

###	Data Location Syntax

	segment-name.field-sequence-number.component.subcomponent
	Segments are specified using the three letter name (MSH)
	Fields are specifified by the sequence number of the HL7 specification (Field 0 is always the segment name)
	Components and Subcomponents are 0 based indexes
	"" returns the message
	"PID" returns the PID segment
	"PID.5" returns the 5th field of the PID segment
	"PID.5.1" returns the 1st component of the 5th field of the PID segment
	"PID.5.1.2" returns the 2nd subcomponent of the 1st component of the 5th field of the PID
    
### Message building

```go
type MSHSegment struct {
	FieldSeparator     string `hl7:"MSH.1" hl7default:"^~\\&"`
	EncodingCharacters string `hl7:"MSH.2" hl7default:""`
    SendingApp        string `hl7:"MSH.3"`
    SendingFacility   string `hl7:"MSH.4"`
    ReceivingApp      string `hl7:"MSH.5"`
    ReceivingFacility string `hl7:"MSH.6"`
    MsgDate           string `hl7:"MSH.7"`
    MessageType       string `hl7:"MSH.9"`
    ControlID         string `hl7:"MSH.10"`
    ProcessingID      string `hl7:"MSH.11"`
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

func CreateHL7() {
    msh := MSHSegment{
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
    pv1 := datamodel.PV1Segment{SetID: "0001", PatientClass: "I", AdmissionType: "O"}

    file, err := os.OpenFile("test.hl7", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    defer os.close(file)

    err = hl7.NewEncoder(file).Encode(&msh)
    err = hl7.NewEncoder(file).Encode(&pid)
    err = hl7.NewEncoder(file).Encode(&pv1)
}
```

###	Data Extraction / Unmarshal

```go
data := []byte(...) // raw message
type PIDSegment struct {
	FirstName string `hl7:"PID.5.1"`
	LastName  string `hl7:"PID.5.0"`
}
ps := PIDSegment{}

err := hl7.Unmarshal(data, &ps)

// from an io.Reader

hl7.NewDecoder(reader).Decode(&ps)
```

### Message Query

```go
msg, err := hl7.NewDecoder(reader).Message()

// First matching value
val, err := msg.Find("PID.5.1")

// All matching values
vals, err := msg.FindAll("PID.11.1")
```

### Message Validation

Message validation is accomplished using the IsValid function. Create a slice of Validation structs and pass them, with the message, to the IsValid function. The first return value is a pass / fail bool. The second return value returns the Validation structs that failed.

A number of validation slices are already defined and can be combined to build custom validation criteria. The NewValidMSH24() function is one example. It returns a set of validations for the MSH segment for version 2.4 of the HL7 specification.

```go
val := []hl7.Validation{
	Validation{Location: "MSH.0", VCheck: SpecificValue, Value: "MSH"},
	Validation{Location: "MSH.1", VCheck: HasValue},
	Validation{Location: "MSH.2", VCheck: HasValue},
}
msg, err := hl7.NewDecoder(reader).Message()
valid, failures := msg.IsValid(val)
```

## License
Copyright 2015, 2016 Davin Hills. All rights reserved.
MIT license. License details can be found in the LICENSE file.


