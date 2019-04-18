package models

import (
	"fmt"
	"time"
)

type TrackRequest struct {
	Request
	SelectionDetails  SelectionDetails `xml:"q0:SelectionDetails"`
	ProcessingOptions string           `xml:"q0:ProcessingOptions"`
}

type TrackResponseEnvelope struct {
	Reply TrackReply `xml:"Body>TrackReply"`
}

func (t *TrackResponseEnvelope) Error() error {
	// TrackResponses are odd in that for invalid tracking numbers, the Reply
	// doesn't say it errored, even though the Reply.CompletedTrackDetails does

	// Error if Reply has error
	err := t.Reply.Error()
	if err != nil {
		return fmt.Errorf("track reply error: %s", err)
	}

	// Error if CompletedTrackDetails has error
	for _, completedTrackDetail := range t.Reply.CompletedTrackDetails {
		for _, trackDetail := range completedTrackDetail.TrackDetails {
			if trackDetail.Notification.Severity == notificationSeverityError {
				return fmt.Errorf("track detail error: %s", trackDetail.Notification.Message)
			}
		}
	}

	return nil
}

// TrackReply : Track reply root (`xml:"Body>TrackReply"`)
type TrackReply struct {
	Reply
	CompletedTrackDetails []CompletedTrackDetail
}

// ActualDelivery returns the first ACTUAL_DELIVERY timestamp
func (tr *TrackReply) ActualDelivery() *time.Time {
	return tr.searchDatesOrTimes("ACTUAL_DELIVERY")
}

// EstimatedDelivery returns the first ESTIMATED_DELIVERY timestamp
func (tr *TrackReply) EstimatedDelivery() *time.Time {
	return tr.searchDatesOrTimes("ESTIMATED_DELIVERY")
}

// Ship returns the first SHIP timestamp
func (tr *TrackReply) Ship() *time.Time {
	return tr.searchDatesOrTimes("SHIP")
}

func (tr *TrackReply) Events() []Event {
	events := []Event{}
	for _, completedTrackDetail := range tr.CompletedTrackDetails {
		for _, trackDetail := range completedTrackDetail.TrackDetails {
			events = append(events, trackDetail.Events...)
		}
	}
	return events
}

func (tr *TrackReply) searchDatesOrTimes(dateOrTimeType string) *time.Time {
	for _, completedTrackDetail := range tr.CompletedTrackDetails {
		for _, trackDetail := range completedTrackDetail.TrackDetails {
			for _, dateOrTime := range trackDetail.DatesOrTimes {
				if dateOrTime.Type == dateOrTimeType {
					ts := time.Time(dateOrTime.DateOrTimestamp)
					return &ts
				}
			}
		}
	}

	return nil
}
