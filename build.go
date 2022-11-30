package hl7

import (
	"fmt"
	"time"
)

// MsgInfo describes the basic message fields
type MsgInfo struct {
	ServiceSymbol     string `hl7:"MSH.1"`
	SendingApp        string `hl7:"MSH.2"`
	SendingFacility   string `hl7:"MSH.3"`
	ReceivingApp      string `hl7:"MSH.4"`
	ReceivingFacility string `hl7:"MSH.5"`
	MsgDate           string `hl7:"MSH.6"`  // if blank will generate
	MessageType       string `hl7:"MSH.8"`  // Required example ORM^001
	ControlID         string `hl7:"MSH.9"` // if blank will generate
	ProcessingID      string `hl7:"MSH.10"` // default P
	VersionID         string `hl7:"MSH.11"` // default 2.4
}

// NewMsgInfo returns a MsgInfo with controlID, message date, Processing Id, and Version set
// Version = 2.4
// ProcessingID = P
func NewMsgInfo() *MsgInfo {
	info := MsgInfo{}
	now := time.Now()
	t := now.Format("20060102150405")
	info.MsgDate = t
	info.ControlID = fmt.Sprintf("MSGID%s%d", t, now.Nanosecond())
	info.ProcessingID = "P"
	info.VersionID = "2.4"
	return &info
}

// NewMsgInfoAck returns a MsgInfo ACK based on the MsgInfo passed in
func NewMsgInfoAck(mi *MsgInfo) *MsgInfo {
	info := NewMsgInfo()
	info.MessageType = "ACK"
	info.ReceivingApp = mi.SendingApp
	info.ReceivingFacility = mi.SendingFacility
	info.SendingApp = mi.ReceivingApp
	info.SendingFacility = mi.ReceivingFacility
	info.ProcessingID = mi.ProcessingID
	info.VersionID = mi.VersionID
	return info
}

// StartMessage returns a Message with an MSH segment based on the MsgInfo struct
func StartMessage(info MsgInfo) (*Message, error) {
	if info.MessageType == "" {
		return nil, fmt.Errorf("Message Type is required")
	}
	now := time.Now()
	t := now.Format("20060102150405")
	if info.MsgDate == "" {
		info.MsgDate = t
	}
	if info.ControlID == "" {
		info.ControlID = fmt.Sprintf("MSGID%s%d", t, now.Nanosecond())
	}
	if info.ProcessingID == "" {
		info.ProcessingID = "P"
	}
	if info.VersionID == "" {
		info.VersionID = "2.4"
	}
	msg := NewMessage([]byte{})
	Marshal(msg, &info)
	return msg, nil
}
