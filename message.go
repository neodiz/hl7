package hl7

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"golang.org/x/net/html/charset"
)

// Message is an HL7 message
type Message struct {
	Segments   []Segment
	Value      []rune
	Delimeters Delimeters
}

// NewMessage returns a new message with the v byte value
func NewMessage(v []byte) *Message {
	var utf8V []byte
	if len(v) != 0 {
		reader, err := charset.NewReader(bytes.NewReader(v), "text/plain")
		if err != nil {
			return nil
		}
		utf8V, err = ioutil.ReadAll(reader)
	} else {
		utf8V = v
	}
	return &Message{
		Value:      []rune(string(utf8V)),
		Delimeters: *NewDelimeters(),
	}
}
func (m *Message) Parse() error {
   return m.parse()
}

func (m *Message) String() string {
	var str string
	for _, s := range m.Segments {
		str += "Message Segment: " + string(s.Value) + "\n"
		str += s.String()
	}
	return str
}

// Segment returns the first matching segmane with name s
func (m *Message) Segment(s string) (*Segment, error) {
	for i, seg := range m.Segments {
		fld, err := seg.Field(0)
		if err != nil {
			continue
		}
		if string(fld.Value) == s {
			return &m.Segments[i], nil
		}
	}
	return nil, ErrSegmentNotFound
}

// AllSegments returns the first matching segmane with name s
func (m *Message) AllSegments(s string) ([]*Segment, error) {
	segs := []*Segment{}
	for i, seg := range m.Segments {
		fld, err := seg.Field(0)
		if err != nil {
			continue
		}
		if string(fld.Value) == s {
			segs = append(segs, &m.Segments[i])
		}
	}
	if len(segs) == 0 {
		return segs, ErrSegmentNotFound
	}
	return segs, nil
}

// Find gets a value from a message using location syntax
// finds the first occurence of the segment and first of repeating fields
// if the loc is not valid an error is returned
func (m *Message) Find(loc string) (string, error) {
	return m.Get(NewLocation(loc))
}

// FindAll gets all values from a message using location syntax
// finds all occurences of the segments and all repeating fields
// if the loc is not valid an error is returned
func (m *Message) FindAll(loc string) ([]string, error) {
	return m.GetAll(NewLocation(loc))
}

// Get returns the first value specified by the Location
func (m *Message) Get(l *Location) (string, error) {
	if l.Segment == "" {
		return string(m.Value), nil
	}
	seg, err := m.Segment(l.Segment)
	if err != nil {
		return "", err
	}
	return seg.Get(l)
}

// GetAll returns all values specified by the Location
func (m *Message) GetAll(l *Location) ([]string, error) {
	vals := []string{}
	if l.Segment == "" {
		vals = append(vals, string(m.Value))
		return vals, nil
	}
	segs, err := m.AllSegments(l.Segment)
	if err != nil {
		return vals, err
	}
	for _, s := range segs {
		vs, err := s.GetAll(l)
		if err != nil {
			return vals, err
		}
		vals = append(vals, vs...)
	}
	return vals, nil
}

// Set will insert a value into a message at Location
func (m *Message) Set(l *Location, val string) error {
	if l.Segment == "" {
		return errors.New("Segment is required")
	}
	seg, err := m.Segment(l.Segment)
	if err != nil {
		s := Segment{}
		s.forceField([]rune(l.Segment), 0)
		s.Set(l, val, &m.Delimeters)
		m.Segments = append(m.Segments, s)
	} else {
		seg.Set(l, val, &m.Delimeters)
	}
	m.Value = m.encode()
	return nil
}

func (m *Message) parse() error {
	if err := m.parseSep(); err != nil {
		return err
	}
	r := strings.NewReader(string(m.Value))
	i := 0
	ii := 0
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			ch = eof
		}
		ii++
		if i >= len(m.Value) {
			return nil
		}
		switch {
		case ch == eof || (ch == endMsg && m.Delimeters.LFTermMsg):
			//just for safety: cannot reproduce this on windows
			start := i
			if m.Value[i] == endMsg {
				start++
			}
			safeii := map[bool]int{true: len(m.Value), false: ii}[ii > len(m.Value)]
			v := m.Value[start:safeii]
			if len(v) > 4 { // seg name + field sep
				seg := Segment{Value: v}
				seg.parse(&m.Delimeters)
				m.Segments = append(m.Segments, seg)
			}
			return nil
		case ch == segTerm:
			start := i
			if m.Value[i] == endMsg {
				start++
			}
			seg := Segment{Value: m.Value[start : ii-1]}
			seg.parse(&m.Delimeters)
			m.Segments = append(m.Segments, seg)
			i = ii
		case ch == m.Delimeters.Escape:
			ii++
			r.ReadRune()
		}
	}
}

func (m *Message) parseSep() error {
	if len(m.Value) < 8 {
		return errors.New("Invalid message length less than 8 bytes")
	}
	if string(m.Value[:3]) != "MSH" {
		return errors.New("Invalid message: Missing MSH segment")
	}

	r := bytes.NewReader([]byte(string(m.Value)))
	for i := 0; i < 8; i++ {
		ch, _, _ := r.ReadRune()
		if ch == eof {
			return fmt.Errorf("Invalid message: eof while parsing MSH")
		}
		switch i {
		case 3:
			m.Delimeters.Field = ch
		case 4:
			m.Delimeters.DelimeterField = string(ch)
			m.Delimeters.Component = ch
		case 5:
			m.Delimeters.DelimeterField += string(ch)
			m.Delimeters.Repetition = ch
		case 6:
			m.Delimeters.DelimeterField += string(ch)
			m.Delimeters.Escape = ch
		case 7:
			m.Delimeters.DelimeterField += string(ch)
			m.Delimeters.SubComponent = ch
		}
	}
	return nil
}

func (m *Message) encode() []rune {
	buf := [][]byte{}
	for _, s := range m.Segments {
		buf = append(buf, []byte(string(s.Value)))
	}
	return []rune(string(bytes.Join(buf, []byte(string(segTerm)))))
}

// IsValid checks a message for validity based on a set of criteria
// it returns valid and any failed validation rules
func (m *Message) IsValid(val []Validation) (bool, []Validation) {
	failures := []Validation{}
	valid := true
	for _, v := range val {
		values, err := m.FindAll(v.Location)
		if err != nil || len(values) == 0 {
			valid = false
			failures = append(failures, v)
		}
		for _, value := range values {
			if value == "" || (v.VCheck == SpecificValue && v.Value != value) {
				valid = false
				failures = append(failures, v)
			}
		}
	}

	return valid, failures
}

// Unmarshal fills a structure from an HL7 message
// It will panic if interface{} is not a pointer to a struct
// Unmarshal will decode the entire message before trying to set values
// it will set the first matching segment / first matching field
// repeating segments and fields is not well suited to this
// for the moment all unmarshal target fields must be strings
func (m *Message) Unmarshal(it interface{}) error {
	st := reflect.ValueOf(it).Elem()
	stt := st.Type()
	for i := 0; i < st.NumField(); i++ {
		fld := stt.Field(i)
		r := fld.Tag.Get("hl7")
		if r != "" {
			if !st.Field(i).CanSet() {
				continue
			}

			//TODO support fields other than string and []string
			switch fld.Type.Kind() {
			case reflect.Slice:
				val, _ := m.FindAll(r)
				st.Field(i).Set(reflect.ValueOf(val))

			case reflect.String:
				if val, _ := m.Find(r); val != "" {
					st.Field(i).SetString(strings.TrimSpace(val))
				}
			}
		}
	}

	return nil
}

// Info returns the MsgInfo for the message
func (m *Message) Info() (MsgInfo, error) {
	mi := MsgInfo{}
	err := m.Unmarshal(&mi)
	return mi, err
}
