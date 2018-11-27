package models

import (
	"encoding/xml"
	"time"
)

// Timestamp marshals time.Times using RFC3339
type Timestamp time.Time

// MarshalXML marshals time.Times to strings using time.RFC3339 format
func (ts Timestamp) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(time.Time(ts).Format(time.RFC3339), start)
}

// UnmarshalXML unmarshals strings in time.RFC3339 format to time.Times
func (ts *Timestamp) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	*ts = Timestamp(t)
	return nil
}
