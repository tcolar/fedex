package models

import (
	"encoding/xml"
	"errors"
	"fmt"
)

// Bool marshals false and true to 0 and 1
type Bool bool

// MarshalXML converts {false, true} to {0, 1}
func (b Bool) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if b {
		return e.EncodeElement(1, start)
	}
	return e.EncodeElement(0, start)
}

// UnmarshalXML converts {0, 1} to {false, true}
func (b *Bool) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var i int
	err := d.DecodeElement(&i, &start)
	if err != nil {
		return fmt.Errorf("decode element as int: %s", err)
	}
	switch i {
	case 0:
		*b = Bool(false)
	case 1:
		*b = Bool(true)
	default:
		return errors.New("int must be 0 or 1")
	}
	return nil
}
