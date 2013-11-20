// History: Nov 20 13 tcolar Creation
package fedex

// Structures to unmarshall the Fedex SOAP answer into

// Track reply root (`xml:"Body>TrackReply"`)
type TrackReply struct {
	HighestSeverity       string
	Notifications         []Notification
	Version               Version
	CompletedTrackDetails []CompletedTrackDetail
}

func (r TrackReply) Failed() bool {
	return r.HighestSeverity != "SUCCESS"
}

type Version struct {
	ServiceId    string
	Major        int
	Intermediate int
	Minor        int
}

type CompletedTrackDetail struct {
	HighestSeverity  string
	Notifications    []Notification
	DuplicateWaybill bool
	MoreData         bool
	TrackDetails     []TrackDetail
}

type TrackDetail struct {
	TrackingNumber                         string
	TrackingNumberUniqueIdentifier         string
	Notification                           Notification
	StatusDetail                           StatusDetail
	CarrierCode                            string
	OperatingCompanyOrCarrierDescription   string
	OtherIdentifiers                       []OtherIdentifier
	Service                                Service
	PackageWeight                          Weight
	ShipmentWeight                         Weight
	Packaging                              string
	PackagingType                          string
	PackageSequenceNumber                  int
	PackageCount                           int
	SpecialHandlings                       []SpecialHandling
	ShipTimestamp                          string
	ActualDeliveryTimestamp                string
	DestinationAddress                     Location
	ActualDeliveryAddress                  Location
	DeliveryLocationType                   string
	DeliveryLocationDescription            string
	DeliveryAttempts                       int
	DeliverySignatureName                  string
	TotalUniqueAddressCountInConsolidation int
	NotificationEventsAvailable            string
	RedirectToHoldEligibility              string
	Events                                 []Event
}

type Notification struct {
	Severity         string
	Source           string
	Code             int
	Message          string
	LocalizedMessage string
}

type StatusDetail struct {
	CreationTime     string
	Code             string
	Description      string
	Location         Location
	AncillaryDetails []AncillaryDetail
}

type Location struct {
	StreetLines         string
	City                string
	StateOrProvinceCode string
	CountryCode         string
	CountryName         string
	Residential         bool
}

type AncillaryDetail struct {
	Reason            string
	ReasonDescription string
}

type OtherIdentifier struct {
	PackageIdentifier Identifier
}

type Service struct {
	Type             string
	Description      string
	ShortDescription string
}

type Weight struct {
	Units string
	Value float64
}

type Identifier struct {
	Type  string
	Value string
}

type SpecialHandling struct {
	Type        string
	Description string
	PaymentType string
}

type Event struct {
	Timestamp                  string
	EventType                  string
	EventDescription           string
	StatusExceptionCode        string
	StatusExceptionDescription string
	Address                    Location
	ArrivalLocation            string
}
