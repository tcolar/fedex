package models

import (
	"errors"
	"fmt"
	"time"
)

const notificationSeverityError = "ERROR"
const notificationSeveritySuccess = "SUCCESS"

// Shipment is convenience struct that has fields for creating a shipment (not part of FedEx API)
type Shipment struct {
	FromAddress       Address
	ToAddress         Address
	FromContact       Contact
	ToContact         Contact
	NotificationEmail string
}

// Envelope is the soap wrapper for all requests
type Envelope struct {
	XMLName   string      `xml:"soapenv:Envelope"`
	Body      interface{} `xml:"soapenv:Body"`
	Soapenv   string      `xml:"xmlns:soapenv,attr"`
	Namespace string      `xml:"xmlns:q0,attr"`
}

type Response interface {
	Error() error
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

type ShipResponseEnvelope struct {
	Reply ProcessShipmentReply `xml:"Body>ProcessShipmentReply"`
}

func (s *ShipResponseEnvelope) Error() error {
	return s.Reply.Error()
}

type RateResponseEnvelope struct {
	Reply RateReply `xml:"Body>RateReply"`
}

func (r *RateResponseEnvelope) Error() error {
	return r.Reply.Error()
}

type CreatePickupResponseEnvelope struct {
	Reply CreatePickupReply `xml:"Body>CreatePickupReply"`
}

func (c *CreatePickupResponseEnvelope) Error() error {
	return c.Reply.Error()
}

// Request has just the default auth fields on all requests
type Request struct {
	WebAuthenticationDetail WebAuthenticationDetail `xml:"q0:WebAuthenticationDetail"`
	ClientDetail            ClientDetail            `xml:"q0:ClientDetail"`
	// TransitionDetail        *TransactionDetail      `xml:"q0:TransitionDetail,omitempty"`
	Version Version `xml:"q0:Version"`
}

type WebAuthenticationDetail struct {
	UserCredential UserCredential `xml:"q0:UserCredential"`
}

type ClientDetail struct {
	AccountNumber string `xml:"q0:AccountNumber"`
	MeterNumber   string `xml:"q0:MeterNumber"`
}

type UserCredential struct {
	Key      string `xml:"q0:Key"`
	Password string `xml:"q0:Password"`
}

type RateRequest struct {
	Request
	RequestedShipment RequestedShipment `xml:"q0:RequestedShipment"`
}

type TrackRequest struct {
	Request
	SelectionDetails  SelectionDetails `xml:"q0:SelectionDetails"`
	ProcessingOptions string           `xml:"q0:ProcessingOptions"`
}

type CreatePickupRequest struct {
	Request
	OriginDetail         OriginDetail        `xml:"q0:OriginDetail"`
	FreightPickupDetail  FreightPickupDetail `xml:"q0:FreightPickupDetail"`
	PackageCount         int                 `xml:"q0:PackageCount"`
	CarrierCode          string              `xml:"q0:CarrierCode"`
	Remarks              string              `xml:"q0:Remarks"`
	CommodityDescription string              `xml:"q0:CommodityDescription"`
}

type FreightPickupDetail struct {
	ApprovedBy  Contact                 `xml:"q0:ApprovedBy"`
	Payment     string                  `xml:"q0:Payment"`
	Role        string                  `xml:"q0:Role"`
	SubmittedBy Contact                 `xml:"q0:SubmittedBy"`
	LineItems   []FreightPickupLineItem `xml:"q0:LineItems"`
}

type FreightPickupLineItem struct {
	Service            string  `xml:"q0:Service"`
	SequenceNumber     int     `xml:"q0:SequenceNumber"`
	Destination        Address `xml:"q0:Destination"`
	Packaging          string  `xml:"q0:Packaging"`
	Pieces             int     `xml:"q0:Pieces"`
	Weight             Weight  `xml:"q0:Weight"`
	TotalHandlingUnits int     `xml:"q0:TotalHandlingUnits"`
	JustOneMore        bool    `xml:"q0:JustOneMore"`
	Description        string  `xml:"q0:Description"`
}

type OriginDetail struct {
	UseAccountAddress       Bool           `xml:"q0:UseAccountAddress"`
	PickupLocation          PickupLocation `xml:"q0:PickupLocation"`
	PackageLocation         string         `xml:"q0:PackageLocation"`
	BuildingPart            string         `xml:"q0:BuildingPart"`
	BuildingPartDescription string         `xml:"q0:BuildingPartDescription"`
	ReadyTimestamp          Timestamp      `xml:"q0:ReadyTimestamp"`
	CompanyCloseTime        string         `xml:"q0:CompanyCloseTime"`
}

type PickupLocation struct {
	Contact Contact `xml:"q0:Contact"`
	Address Address `xml:"q0:Address"`
}

type SelectionDetails struct {
	CarrierCode       string            `xml:"q0:CarrierCode"`
	PackageIdentifier PackageIdentifier `xml:"q0:PackageIdentifier"`
	// Destination           Destination
	// ShipmentAccountNumber string
}

type PackageIdentifier struct {
	Type  string `xml:"q0:Type"`
	Value string `xml:"q0:Value"`
}

type Version struct {
	ServiceID    string `xml:"q0:ServiceId"`
	Major        int    `xml:"q0:Major"`
	Intermediate int    `xml:"q0:Intermediate"`
	Minor        int    `xml:"q0:Minor"`
}

type ProcessShipmentRequest struct {
	Request
	RequestedShipment RequestedShipment `xml:"q0:RequestedShipment"`
}

type RequestedShipment struct {
	ShipTimestamp Timestamp `xml:"q0:ShipTimestamp"`
	DropoffType   string    `xml:"q0:DropoffType"`
	ServiceType   string    `xml:"q0:ServiceType"`
	PackagingType string    `xml:"q0:PackagingType"`

	// We don't use these, but may do so later
	// ShipmentManifestDetail      *ShipmentManifestDetail      `xml:"q0:ShipmentManifestDetail,omitempty"`
	// TotalWeight                 *Weight                      `xml:"q0:TotalWeight,omitempty"`
	// TotalInsuredValue           *Money                       `xml:"q0:TotalInsuredValue,omitempty"`
	// PreferredCurrency           string                       `xml:"q0:PreferredCurrency,omitempty"`
	// ShipmentAuthorizationDetail *ShipmentAuthorizationDetail `xml:"q0:ShipmentAuthorizationDetail,omitempty"`

	Shipper   Shipper `xml:"q0:Shipper"`
	Recipient Shipper `xml:"q0:Recipient"`

	ShippingChargesPayment    Payment                    `xml:"q0:ShippingChargesPayment"`
	SpecialServicesRequested  *SpecialServicesRequested  `xml:"q0:SpecialServicesRequested,omitempty"`
	SmartPostDetail           *SmartPostDetail           `xml:"q0:SmartPostDetail,omitempty"`
	LabelSpecification        LabelSpecification         `xml:"q0:LabelSpecification"`
	RateRequestTypes          string                     `xml:"q0:RateRequestTypes"`
	PackageCount              int                        `xml:"q0:PackageCount"`
	RequestedPackageLineItems []RequestedPackageLineItem `xml:"q0:RequestedPackageLineItems"`
}

type SpecialServicesRequested struct {
	SpecialServiceTypes     []string                 `xml:"q0:SpecialServiceTypes,omitempty"`
	EventNotificationDetail *EventNotificationDetail `xml:"q0:EventNotificationDetail,omitempty"`
	ReturnShipmentDetail    *ReturnShipmentDetail    `xml:"q0:ReturnShipmentDetail,omitempty"`
}

type EventNotificationDetail struct {
	AggregationType    string              `xml:"q0:AggregationType"`
	PersonalMessage    string              `xml:"q0:PersonalMessage"`
	EventNotifications []EventNotification `xml:"q0:EventNotifications"`
}

type EventNotification struct {
	Role                string              `xml:"q0:Role"`
	Events              []string            `xml:"q0:Events"`
	NotificationDetail  NotificationDetail  `xml:"q0:NotificationDetail"`
	FormatSpecification FormatSpecification `xml:"q0:FormatSpecification"`
}

type NotificationDetail struct {
	NotificationType string       `xml:"q0:NotificationType"`
	EmailDetail      EmailDetail  `xml:"q0:EmailDetail"`
	Localization     Localization `xml:"q0:Localization"`
}

type Localization struct {
	LanguageCode string `xml:"q0:LanguageCode"`
}

type EmailDetail struct {
	EmailAddress string `xml:"q0:EmailAddress"`
	Name         string `xml:"q0:Name"`
}

type FormatSpecification struct {
	Type string `xml:"q0:Type"`
}

type NotificationFormatType struct {
}

type ReturnShipmentDetail struct {
	ReturnType string `xml:"q0:ReturnType"`
}

type ShipmentManifestDetail struct {
	ManifestReferenceType string `xml:"q0:ManifestReferenceType,omitempty"`
}

type SmartPostDetail struct {
	Indicia              string `xml:"q0:Indicia"`
	AncillaryEndorsement string `xml:"q0:AncillaryEndorsement"`
	HubID                string `xml:"q0:HubId"`
}
type RequestedPackageLineItem struct {
	SequenceNumber     int                 `xml:"q0:SequenceNumber"`
	GroupPackageCount  int                 `xml:"q0:GroupPackageCount,omitempty"`
	Weight             Weight              `xml:"q0:Weight"`
	Dimensions         Dimensions          `xml:"q0:Dimensions"`
	PhysicalPackaging  string              `xml:"q0:PhysicalPackaging"`
	ItemDescription    string              `xml:"q0:ItemDescription"`
	CustomerReferences []CustomerReference `xml:"q0:CustomerReferences"`
}

type CustomerReference struct {
	CustomerReferenceType string `xml:"q0:CustomerReferenceType"`
	Value                 string `xml:"q0:Value"`
}

type Weight struct {
	Units string  `xml:"q0:Units"`
	Value float64 `xml:"q0:Value"`
}

type Contact struct {
	PersonName   string `xml:"q0:PersonName"`
	CompanyName  string `xml:"q0:CompanyName"`
	PhoneNumber  string `xml:"q0:PhoneNumber"`
	EmailAddress string `xml:"q0:EMailAddress"`
}

type Dimensions struct {
	Length int    `xml:"q0:Length"`
	Width  int    `xml:"q0:Width"`
	Height int    `xml:"q0:Height"`
	Units  string `xml:"q0:Units"`
}

type Payment struct {
	PaymentType string `xml:"q0:PaymentType"`
	Payor       Payor  `xml:"q0:Payor"`
}

type Payor struct {
	ResponsibleParty ResponsibleParty `xml:"q0:ResponsibleParty"`
}

type ResponsibleParty struct {
	AccountNumber string `xml:"q0:AccountNumber"`
}

type LabelSpecification struct {
	LabelFormatType string `xml:"q0:LabelFormatType"`
	ImageType       string `xml:"q0:ImageType"`
}

type Shipper struct {
	AccountNumber string  `xml:"q0:AccountNumber"`
	Contact       Contact `xml:"q0:Contact"`
	Address       Address `xml:"q0:Address"`
}

// Reply has common stuff on all responses from FedEx API
type Reply struct {
	HighestSeverity string
	Notifications   []Notification
	Version         VersionResponse
	JobID           string `xml:"JobId"`
}

func (r Reply) Error() error {
	if r.HighestSeverity == notificationSeveritySuccess {
		return nil
	}

	for _, notification := range r.Notifications {
		if notification.Severity == r.HighestSeverity {
			return fmt.Errorf("reply got error: %s", notification.Message)
		}
	}
	return fmt.Errorf("reply got status: %s", r.HighestSeverity)

}

// TrackReply : Track reply root (`xml:"Body>TrackReply"`)
type TrackReply struct {
	Reply
	CompletedTrackDetails []CompletedTrackDetail
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

// ProcessShipReply : Process shipment reply root (`xml:"Body>ProcessShipmentReply"`)
type ProcessShipmentReply struct {
	Reply
	TransactionDetail       TransactionDetail
	CompletedShipmentDetail CompletedShipmentDetail
}

// RateReply : Process shipment reply root (`xml:"Body>RateReply"`)
type RateReply struct {
	Reply
	TransactionDetail TransactionDetail
	RateReplyDetails  []RateReplyDetail
}

// TotalCost returns the first TotalNetChargeWithDutiesAndTaxes in the reply
func (rr *RateReply) TotalCost() (Charge, error) {
	// TotalNetChargeWithDutiesAndTaxes
	for _, rateReplyDetail := range rr.RateReplyDetails {
		for _, ratedShipmentDetail := range rateReplyDetail.RatedShipmentDetails {
			totalNetCharge := ratedShipmentDetail.ShipmentRateDetail.TotalNetChargeWithDutiesAndTaxes
			if totalNetCharge.Currency != "" && totalNetCharge.Amount != "" {
				return totalNetCharge, nil
			}
		}
	}
	return Charge{}, errors.New("no total net charge found on reply")
}

// CreatePickupReply : CreatePickup reply root (`xml:"Body>CreatePickupReply"`)
type CreatePickupReply struct {
	Reply
	PickupConfirmationNumber string
	Location                 string
}

type RateReplyDetail struct {
	ServiceType                     string
	ServiceDescription              ServiceDescription
	PackagingType                   string
	DestinationAirportID            string `xml:"DestinationAirportId"`
	IneligibleForMoneyBackGuarantee bool
	SignatureOption                 string
	ActualRateType                  string
	RatedShipmentDetails            []Rating // TODO
}

type TransactionDetail struct {
	CustomerTransactionID string `xml:"q0:CustomerTransactionId,omitempty"`
}

type CompletedShipmentDetail struct {
	UsDomestic              string
	CarrierCode             string
	MasterTrackingId        TrackingID
	ServiceTypeDescription  string
	ServiceDescription      ServiceDescription
	PackagingDescription    string
	OperationalDetail       OperationalDetail
	ShipmentRating          Rating
	CompletedPackageDetails CompletedPackageDetails
}

type Part struct {
	DocumentPartSequenceNumber string
	Image                      []byte
}

type Label struct {
	Type                        string
	ShippingDocumentDisposition string
	ImageType                   string
	Resolution                  string
	CopiesToPrint               string
	Parts                       []Part
}

type CompletedPackageDetails struct {
	SequenceNumber string
	TrackingIds    []TrackingID
	Label          Label
}

type TrackingID struct {
	TrackingIdType string
	TrackingNumber string
}

type Name struct {
	Type     string
	Encoding string
	Value    string
}

type ServiceDescription struct {
	ServiceType      string
	Code             string
	Names            []Name
	Description      string
	AstraDescription string
}

type Surcharge struct {
	SurchargeType string
	Level         string
	Description   string
	Amount        Charge
}

type OperationalDetail struct {
	OriginLocationNumber            string
	DestinationLocationNumber       string
	TransitTime                     string
	IneligibleForMoneyBackGuarantee string
	DeliveryEligibilities           string
	ServiceCode                     string
	PackagingCode                   string
}

type RatedShipmentDetail struct {
	EffectiveNetDiscount Charge
	ShipmentRateDetail   RateDetail
	RatedPackages        []RatedPackage
}

type Rating struct {
	ActualRateType       string
	GroupNumber          string
	EffectiveNetDiscount Charge

	// For the shipping service, the rate details is an array, but for the rate service, it is not
	ShipmentRateDetails []RateDetail
	ShipmentRateDetail  RateDetail

	RatedPackages []RatedPackage
}

type Charge struct {
	Currency string
	Amount   string
}

type RatedPackage struct {
	GroupNumber          string
	EffectiveNetDiscount Charge
	PackageRateDetail    RateDetail
}

type RateDetail struct {
	RateType                         string
	RateZone                         string
	RatedWeightMethod                string
	DimDivisor                       string
	FuelSurchargePercent             string
	TotalBillingWeight               Weight
	TotalBaseCharge                  Charge
	TotalFreightDiscounts            Charge
	TotalNetFreight                  Charge
	TotalSurcharges                  Charge
	TotalNetFedExCharge              Charge
	TotalTaxes                       Charge
	TotalNetCharge                   Charge
	NetCharge                        Charge
	TotalRebates                     Charge
	TotalDutiesAndTaxes              Charge
	TotalAncillaryFeesAndTaxes       Charge
	TotalDutiesTaxesAndFees          Charge
	TotalNetChargeWithDutiesAndTaxes Charge
	Surcharges                       []Surcharge
}

type VersionResponse struct {
	ServiceID    string `xml:"ServiceId"`
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
	Notification                   Notification
	TrackingNumber                 string
	Barcode                        StringBarcode
	TrackingNumberUniqueIdentifier string
	StatusDetail                   StatusDetail
	InformationNotes               []InformationNoteDetail

	// Not gonna bother with all of these fields until we need them
	// Most of the fields in this block are not important
	CustomerExceptionRequests            []InformationNoteDetail
	Reconciliations                      []Reconciliation
	ServiceCommitMessage                 string
	DestinationServiceArea               string
	DestinationServiceAreaDescription    string
	CarrierCode                          string
	OperatingCompanyType                 string
	OperatingCompanyOrCarrierDescription string
	CartageAgentCompanyName              string
	ProductionLocationContactAndAddress  ContactAndAddress
	ContentRecord                        ContentRecord
	// ... more

	Service               Service
	PackageWeight         Weight
	ShipmentWeight        Weight
	Packaging             string
	PackagingType         string
	PhysicalPackagingType string
	PackageSequenceNumber int
	PackageCount          int
	Charges               Charge
	NickName              string
	Notes                 string
	Attributes            []string
	ShipmentContents      []ContentRecord
	PackageContents       string

	TrackAdvanceNotificationDetail AdvanceNotificationDetail
	Shipper                        Contact
	ShipperAddress                 Address
	OriginLocationAddress          Address

	// DatesOrTimes contains estimated arrivals, departures, etc.
	DatesOrTimes []DateOrTimestamp

	Recipient                              Contact
	DestinationAddress                     Address
	ActualDeliveryAddress                  Address
	SpecialHandlings                       []SpecialHandling
	DeliveryLocationType                   string
	DeliveryLocationDescription            string
	DeliveryAttempts                       int
	DeliverySignatureName                  string
	TotalUniqueAddressCountInConsolidation int
	NotificationEventsAvailable            string
}

type DateOrTimestamp struct {
	Type            string
	DateOrTimestamp Timestamp
}

type AdvanceNotificationDetail struct {
	EstimatedTimeOfArrival Timestamp
	Reason                 string
	Status                 string
	StatusDescription      string
	StatusTime             Timestamp
}

type ContentRecord struct {
	PartNumber       string
	ItemNumber       string
	ReceivedQuantity int
	Description      string
}

type ContactAndAddress struct {
	Contact Contact `xml:"q0:Contact"`
	Address Address `xml:"q0:Address"`
}

type Reconciliation struct {
	Status      string
	Description string
}

type InformationNoteDetail struct {
	Code        string
	Description string
}

type StringBarcode struct {
	Type  string
	Value string
}

type Notification struct {
	Severity         string
	Source           string
	Code             string
	Message          string
	LocalizedMessage string
}

type StatusDetail struct {
	CreationTime     Timestamp
	Code             string
	Description      string
	Location         Address
	AncillaryDetails []AncillaryDetail
}

type Address struct {
	StreetLines         []string `xml:"q0:StreetLines"`
	City                string   `xml:"q0:City"`
	StateOrProvinceCode string   `xml:"q0:StateOrProvinceCode"`
	PostalCode          string   `xml:"q0:PostalCode"`
	CountryCode         string   `xml:"q0:CountryCode"`
	// CountryName         string   `xml:"q0:CountryName"`
	Residential Bool `xml:"q0:Residential"`
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
	Address                    Address
	ArrivalLocation            string
}
