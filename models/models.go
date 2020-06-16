package models

import (
	"fmt"
	"math"
)

type Address struct {
	StreetLines         []string `xml:"q0:StreetLines"`
	City                string   `xml:"q0:City"`
	StateOrProvinceCode string   `xml:"q0:StateOrProvinceCode"`
	PostalCode          string   `xml:"q0:PostalCode"`
	CountryCode         string   `xml:"q0:CountryCode"`
	// CountryName         string   `xml:"q0:CountryName"`
	Residential Bool `xml:"q0:Residential"`
}

func (a Address) ShipsOutWithInternationalEconomy() bool {
	switch a.CountryCode {
	case "CA", "US", "":
		return false
	default:
		return true
	}
}

type AddressReply struct {
	// When we use Address in a request, we want to xml marshall the fields to
	// "q0:fieldName" but on unmarshalling an Address in a response, we want to
	// unmarshal the fields as just "fieldName"
	// TODO instead of two structs, we could merge this with Address, add
	// omitemptys to each of the below fields, so that whenever we marshall, we
	// marshall the q0 fields, but on unmarshalling a response, we check the
	// below fields. i'm not sure what's best
	StreetLines         []string `xml:"StreetLines"`
	City                string   `xml:"City"`
	StateOrProvinceCode string   `xml:"StateOrProvinceCode"`
	PostalCode          string   `xml:"PostalCode"`
	CountryCode         string   `xml:"CountryCode"`
	CountryName         string   `xml:"CountryName"`
}

type AdvanceNotificationDetail struct {
	EstimatedTimeOfArrival Timestamp
	Reason                 string
	Status                 string
	StatusDescription      string
	StatusTime             Timestamp
}

type AncillaryDetail struct {
	Reason            string
	ReasonDescription string
}

type Broker struct {
	Type   string  `xml:"q0:Type"`
	Broker Shipper `xml:"q0:Broker"`
}

type Charge struct {
	Currency string
	Amount   float64
}

type Commodity struct {
	Name           string `xml:"q0:Name"`
	NumberOfPieces int    `xml:"q0:NumberOfPieces"`
	Description    string `xml:"q0:Description"`
	// Purpose *string
	CountryOfManufacture string  `xml:"q0:CountryOfManufacture"`
	HarmonizedCode       *string `xml:"q0:HarmonizedCode"`
	Weight               Weight  `xml:"q0:Weight"`
	Quantity             int     `xml:"q0:Quantity"`
	QuantityUnits        string  `xml:"q0:QuantityUnits"`
	// AdditionalMeasure *int
	UnitPrice                   *Money   `xml:"q0:UnitPrice"`
	CustomsValue                *Money   `xml:"q0:CustomsValue"`
	ExportLicenseExpirationDate *string  `xml:"q0:ExportLicenseExpirationDate"`
	CIMarksAndNumbers           []string `xml:"q0:CIMarksAndNumbers"`
}

type Commodities []Commodity

func (c Commodities) Weight() Weight {
	if len(c) == 0 {
		return Weight{Units: "LB", Value: 0.0}
	}

	// Assume all the units are the same
	weight := Weight{
		Units: c[0].Weight.Units,
		Value: 0.0,
	}
	for _, commodity := range c {
		weight.Value += commodity.Weight.Value
	}
	return weight
}

func (c Commodities) CustomsValue() (Money, error) {
	total := Money{Currency: "USD"}

	if len(c) == 0 {
		return total, nil
	}

	// Set the currency
	for _, commodity := range c {
		if commodity.CustomsValue != nil {
			total.Currency = commodity.CustomsValue.Currency
			break
		}
	}

	for _, commodity := range c {
		if commodity.CustomsValue == nil {
			continue
		}
		if commodity.CustomsValue.Currency != total.Currency {
			return total, fmt.Errorf("mismatching customs currencies: %s %s", commodity.CustomsValue.Currency, total.Currency)
		}
		total.Amount += commodity.CustomsValue.Amount
	}

	total.Amount = math.Ceil(total.Amount*100.0) / 100.0
	return total, nil
}

type CompletedPackageDetails struct {
	SequenceNumber string
	TrackingIds    []TrackingID
	Label          Label
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
	ShipmentDocuments       []ShipmentDocument
	CompletedPackageDetails CompletedPackageDetails
}

type CompletedTrackDetail struct {
	HighestSeverity  string
	Notifications    []Notification
	DuplicateWaybill bool
	MoreData         bool
	TrackDetails     []TrackDetail
}

type CommercialInvoice struct {
	Purpose        string `xml:"q0:Purpose"`
	OriginatorName string `xml:"q0:OriginatorName"`
}

type CommercialInvoiceDetail struct {
	Format              Format               `xml:"q0:Format"`
	CustomerImageUsages []CustomerImageUsage `xml:"q0:CustomerImageUsages"`
}

type Contact struct {
	PersonName   string `xml:"q0:PersonName"`
	CompanyName  string `xml:"q0:CompanyName"`
	PhoneNumber  string `xml:"q0:PhoneNumber"`
	EmailAddress string `xml:"q0:EMailAddress"`
}

type ContactAndAddress struct {
	Contact Contact `xml:"q0:Contact"`
	Address Address `xml:"q0:Address"`
}

type ContentRecord struct {
	PartNumber       string
	ItemNumber       string
	ReceivedQuantity int
	Description      string
}

type CustomerImageUsage struct {
	Type string `xml:"q0:Type"`
	ID   string `xml:"q0:Id"`
}

type CustomerReference struct {
	CustomerReferenceType string `xml:"q0:CustomerReferenceType"`
	Value                 string `xml:"q0:Value"`
}

type CustomsClearanceDetail struct {
	Brokers                        []Broker           `xml:"q0:Brokers"`
	ImporterOfRecord               Shipper            `xml:"q0:ImporterOfRecord"`
	DutiesPayment                  Payment            `xml:"q0:DutiesPayment"`
	DocumentContent                *string            `xml:"q0:DocumentContent"`
	CustomsValue                   *Money             `xml:"q0:CustomsValue"`
	PartiesToTransactionAreRelated bool               `xml:"q0:PartiesToTransactionAreRelated"`
	CommercialInvoice              *CommercialInvoice `xml:"q0:CommercialInvoice"`
	Commodities                    Commodities        `xml:"q0:Commodities"`
}

type DateOrTimestamp struct {
	Type            string
	DateOrTimestamp Timestamp
}

type Destination struct {
	StreetLines         []string
	City                string
	StateOrProvinceCode string
	PostalCode          string
	CountryCode         string
	CountryName         string
	Residential         bool
}

type Dimensions struct {
	Length int    `xml:"q0:Length"`
	Width  int    `xml:"q0:Width"`
	Height int    `xml:"q0:Height"`
	Units  string `xml:"q0:Units"`
}

func (d Dimensions) IsValid() bool {
	valuesAreValid := d.Length > 0 && d.Width > 0 && d.Height > 0
	unitsIsValid := d.Units == DimensionsUnitsIn || d.Units == DimensionsUnitsCm

	return valuesAreValid && unitsIsValid
}

type EmailDetail struct {
	EmailAddress string `xml:"q0:EmailAddress"`
	Name         string `xml:"q0:Name"`
}

type EtdDetail struct {
	RequestedDocumentCopies string `xml:"q0:RequestedDocumentCopies"`
}

type EdtCommodityTax struct {
	HarmonizedCode string
	Taxes          []EdtTaxDetail
}

type EdtTaxDetail struct {
	TaxType      string
	Name         string
	TaxableValue Charge
	Description  string
	Formula      string
	Amount       Charge
}

type Event struct {
	Timestamp                  Timestamp
	EventType                  string
	EventDescription           string
	StatusExceptionCode        string
	StatusExceptionDescription string
	Address                    Address
	ArrivalLocation            string
}

type EventNotification struct {
	Role                string              `xml:"q0:Role"`
	Events              []string            `xml:"q0:Events"`
	NotificationDetail  NotificationDetail  `xml:"q0:NotificationDetail"`
	FormatSpecification FormatSpecification `xml:"q0:FormatSpecification"`
}

type EventNotificationDetail struct {
	AggregationType    string              `xml:"q0:AggregationType"`
	PersonalMessage    string              `xml:"q0:PersonalMessage"`
	EventNotifications []EventNotification `xml:"q0:EventNotifications"`
}

type Format struct {
	ImageType string `xml:"q0:ImageType"`
	StockType string `xml:"q0:StockType"`
}

type FormatSpecification struct {
	Type string `xml:"q0:Type"`
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

type FromAndTo struct {
	FromAddress Address
	ToAddress   Address
	FromContact Contact
	ToContact   Contact
}

func (ft FromAndTo) IsInternational() bool {
	fromCountryCode := ft.FromAddress.CountryCode
	if fromCountryCode == "" {
		fromCountryCode = "US"
	}

	toCountryCode := ft.ToAddress.CountryCode
	if toCountryCode == "" {
		toCountryCode = "US"
	}

	return fromCountryCode != toCountryCode
}

type Identifier struct {
	Type  string
	Value string
}

type Image struct {
	ID    string `xml:"q0:Id"`
	Image string `xml:"q0:Image"`
}

type InformationNoteDetail struct {
	Code        string
	Description string
}

type Label struct {
	Type                        string
	ShippingDocumentDisposition string
	ImageType                   string
	Resolution                  string
	CopiesToPrint               string
	Parts                       []Part
}

type LabelSpecification struct {
	LabelFormatType string  `xml:"q0:LabelFormatType"`
	ImageType       string  `xml:"q0:ImageType"`
	LabelStockType  *string `xml:"q0:LabelStockType"`
}

type Localization struct {
	LanguageCode string `xml:"q0:LanguageCode"`
}

type Money struct {
	Currency string  `xml:"q0:Currency"`
	Amount   float64 `xml:"q0:Amount"`
}

type Name struct {
	Type     string
	Encoding string
	Value    string
}

type NotificationDetail struct {
	NotificationType string       `xml:"q0:NotificationType"`
	EmailDetail      EmailDetail  `xml:"q0:EmailDetail"`
	Localization     Localization `xml:"q0:Localization"`
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

type OriginDetail struct {
	UseAccountAddress       Bool           `xml:"q0:UseAccountAddress"`
	PickupLocation          PickupLocation `xml:"q0:PickupLocation"`
	PackageLocation         string         `xml:"q0:PackageLocation"`
	BuildingPart            string         `xml:"q0:BuildingPart"`
	BuildingPartDescription string         `xml:"q0:BuildingPartDescription"`
	ReadyTimestamp          Timestamp      `xml:"q0:ReadyTimestamp"`
	CompanyCloseTime        string         `xml:"q0:CompanyCloseTime"`
}

type PackageIdentifier struct {
	Type  string `xml:"q0:Type"`
	Value string `xml:"q0:Value"`
}

type Part struct {
	DocumentPartSequenceNumber string
	Image                      []byte
}

type Payment struct {
	PaymentType string `xml:"q0:PaymentType"`
	Payor       Payor  `xml:"q0:Payor"`
}

type Payor struct {
	ResponsibleParty Shipper `xml:"q0:ResponsibleParty"`
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
	DutiesAndTaxes                   []EdtCommodityTax
}

type RateReplyDetail struct {
	ServiceType                     string
	ServiceDescription              ServiceDescription
	PackagingType                   string
	DestinationAirportID            string `xml:"DestinationAirportId"`
	IneligibleForMoneyBackGuarantee bool
	SignatureOption                 string
	ActualRateType                  string
	RatedShipmentDetails            []Rating
}

type RatedPackage struct {
	GroupNumber          string
	EffectiveNetDiscount Charge
	PackageRateDetail    RateDetail
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

type RecipientDetail struct {
	NotificationEventsAvailable []string
}

type Reconciliation struct {
	Status      string
	Description string
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

type RequestedShipment struct {
	ShipTimestamp     Timestamp `xml:"q0:ShipTimestamp"`
	DropoffType       string    `xml:"q0:DropoffType"`
	ServiceType       string    `xml:"q0:ServiceType,omitempty"`
	PackagingType     string    `xml:"q0:PackagingType"`
	PreferredCurrency string    `xml:"q0:PreferredCurrency,omitempty"`

	// We don't use these, but may do so later
	// ShipmentManifestDetail      *ShipmentManifestDetail      `xml:"q0:ShipmentManifestDetail,omitempty"`
	// TotalWeight                 *Weight                      `xml:"q0:TotalWeight,omitempty"`
	// TotalInsuredValue           *Money                       `xml:"q0:TotalInsuredValue,omitempty"`
	// ShipmentAuthorizationDetail *ShipmentAuthorizationDetail `xml:"q0:ShipmentAuthorizationDetail,omitempty"`

	Shipper   Shipper `xml:"q0:Shipper"`
	Recipient Shipper `xml:"q0:Recipient"`

	ShippingChargesPayment        *Payment                       `xml:"q0:ShippingChargesPayment"`
	SpecialServicesRequested      *SpecialServicesRequested      `xml:"q0:SpecialServicesRequested,omitempty"`
	SmartPostDetail               *SmartPostDetail               `xml:"q0:SmartPostDetail,omitempty"`
	CustomsClearanceDetail        *CustomsClearanceDetail        `xml:"q0:CustomsClearanceDetail,omitempty"`
	LabelSpecification            *LabelSpecification            `xml:"q0:LabelSpecification"`
	ShippingDocumentSpecification *ShippingDocumentSpecification `xml:"q0:ShippingDocumentSpecification"`
	RateRequestTypes              *string                        `xml:"q0:RateRequestTypes"`
	EdtRequestType                *string                        `xml:"q0:EdtRequestType"`
	PackageCount                  *int                           `xml:"q0:PackageCount"`
	RequestedPackageLineItems     []RequestedPackageLineItem     `xml:"q0:RequestedPackageLineItems"`
}

type ReturnShipmentDetail struct {
	ReturnType string `xml:"q0:ReturnType"`
}

type Package struct {
	TrackingNumber                  string
	TrackingNumberUniqueIdentifiers []string
	CarrierCode                     string
	ShipDate                        string
	Destination                     Destination
	RecipientDetails                []RecipientDetail
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

type Service struct {
	Type             string
	Description      string
	ShortDescription string
}

type ServiceDescription struct {
	ServiceType      string
	Code             string
	Names            []Name
	Description      string
	AstraDescription string
}

type ShipmentDocument struct {
	Type                        string
	ShippingDocumentDisposition string
	ImageType                   string
	Resolution                  string
	CopiesToPrint               string
	Parts                       []Part
}

type ShipmentManifestDetail struct {
	ManifestReferenceType string `xml:"q0:ManifestReferenceType,omitempty"`
}

type Shipper struct {
	AccountNumber string  `xml:"q0:AccountNumber"`
	Contact       Contact `xml:"q0:Contact"`
	Address       Address `xml:"q0:Address"`
}

type ShippingDocumentSpecification struct {
	ShippingDocumentTypes []string `xml:"q0:ShippingDocumentTypes"`
	// CertificateOfOrigin                     []CertificateOfOrigin
	CommercialInvoiceDetail []CommercialInvoiceDetail `xml:"q0:CommercialInvoiceDetail"`
	// CustomPackageDocumentDetail             []CustomPackageDocumentDetail
	// CustomShipmentDocumentDetail            []CustomShipmentDocumentDetail
	// ExportDeclarationDetail                 []ExportDeclarationDetail
	// GeneralAgencyAgreementDetail            []GeneralAgencyAgreementDetail
	// NaftaCertificateOfOriginDetail          []NaftaCertificateOfOriginDetail
	// Op900Detail                             []Op900Detail
	// DangerousGoodsShippersDeclarationDetail []DangerousGoodsShippersDeclarationDetail
	// FreightAddressLabelDetail               []FreightAddressLabelDetail
	// FreightBillOfLadingDetail               []FreightBillOfLadingDetail
	// ReturnInstructionsDetail                []ReturnInstructionsDetail
}

type SmartPostDetail struct {
	Indicia              string `xml:"q0:Indicia"`
	AncillaryEndorsement string `xml:"q0:AncillaryEndorsement"`
	HubID                string `xml:"q0:HubId"`
}

type SpecialHandling struct {
	Type        string
	Description string
	PaymentType string
}

type SpecialServicesRequested struct {
	SpecialServiceTypes     []string                 `xml:"q0:SpecialServiceTypes,omitempty"`
	EventNotificationDetail *EventNotificationDetail `xml:"q0:EventNotificationDetail,omitempty"`
	ReturnShipmentDetail    *ReturnShipmentDetail    `xml:"q0:ReturnShipmentDetail,omitempty"`
	EtdDetail               *EtdDetail               `xml:"q0:EtdDetail,omitempty"`
}

type StatusDetail struct {
	CreationTime     Timestamp
	Code             string
	Description      string
	Location         Address
	AncillaryDetails []AncillaryDetail
}

type StringBarcode struct {
	Type  string
	Value string
}

type Surcharge struct {
	SurchargeType string
	Level         string
	Description   string
	Amount        Charge
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
	ShipperAddress                 AddressReply
	OriginLocationAddress          AddressReply

	// DatesOrTimes contains estimated arrivals, departures, etc.
	DatesOrTimes []DateOrTimestamp

	Recipient                              Contact
	DestinationAddress                     AddressReply
	ActualDeliveryAddress                  AddressReply
	SpecialHandlings                       []SpecialHandling
	DeliveryLocationType                   string
	DeliveryLocationDescription            string
	DeliveryAttempts                       int
	DeliverySignatureName                  string
	TotalUniqueAddressCountInConsolidation int
	NotificationEventsAvailable            string
	Events                                 []Event
}

type TrackingID struct {
	TrackingIdType string
	TrackingNumber string
}

type TransactionDetail struct {
	CustomerTransactionID string `xml:"q0:CustomerTransactionId,omitempty"`
}

type Weight struct {
	Units string  `xml:"q0:Units"`
	Value float64 `xml:"q0:Value"`
}

func (w Weight) IsZero() bool {
	return w.Value == 0.0
}
