package hl7

import (
	"fmt"
)

var (
	ErrSegmentNotFound     = fmt.Errorf("Segment not found")
	ErrFieldNotFound       = fmt.Errorf("Field not found")
	ErrComponentOutOfRange = fmt.Errorf("Component out of range")
)
