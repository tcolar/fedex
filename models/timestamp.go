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

	// Try to unmarshal timestamps as time.RFC3339, then without the timezone
	// If both fail, just don't do anything
	timeFormats := []string{time.RFC3339, "2006-01-02T15:04:05"}
	for _, timeFormat := range timeFormats {
		t, err := time.Parse(timeFormat, s)
		if err != nil {
			continue
		}
		*ts = Timestamp(t)
		return nil
	}

	return nil
}
