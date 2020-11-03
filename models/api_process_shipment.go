package models

import (
	"errors"
	"regexp"
	"time"
)

// Shipment wraps all the Fedex API fields needed for creating a shipment
type Shipment struct {
	FromAndTo

	NotificationEmail string
	References        []string
	Service           string
	Dimensions        Dimensions
	InvoiceNumber     string

	// Only used for international ground shipments
	OriginatorName    string
	Commodities       Commodities
	LetterheadImageID string
}

var (
	nonAlphanumericRegex *regexp.Regexp
)

func init() {
	nonAlphanumericRegex = regexp.MustCompile("[^a-zA-Z0-9]+")
}

func (s *Shipment) ServiceType() string {
	return ServiceType(s.FromAndTo, s.Service)
}

func (s *Shipment) Broker() string {
	switch s.ServiceType() {
	case ServiceTypeInternationalEconomy:
		return "FedEx Express"
	default:
		return "FedEx Logistics"
	}
}

func (s *Shipment) ShipTime() time.Time {
	t := time.Now()
	if s.IsInternational() {
		t = t.Add(9 * 24 * time.Hour)
	}

	return t
}

func (s *Shipment) ShippingDocumentSpecification() *ShippingDocumentSpecification {
	if s.ServiceType() == ServiceTypeSmartPost || !s.IsInternational() {
		return nil
	}

	letterheadImageID := s.LetterheadImageID
	if s.LetterheadImageID == "" {
		letterheadImageID = "IMAGE_1"
	}

	return &ShippingDocumentSpecification{
		ShippingDocumentTypes: []string{DocumentTypeCommercialInvoice},
		CommercialInvoiceDetail: []CommercialInvoiceDetail{
			{
				Format: Format{
					ImageType: ImageTypePDF,
					StockType: StockTypePaperLetter,
				},
				CustomerImageUsages: []CustomerImageUsage{
					{
						Type: CustomerImageUsageTypeLetterHead,
						ID:   letterheadImageID,
					},
					{
						Type: CustomerImageUsageTypeSignature,
						ID:   "IMAGE_2",
					},
				},
			},
		},
	}
}

func (s *Shipment) LabelSpecification() *LabelSpecification {
	if s.IsInternational() {
		stockType := StockTypePaper4x6
		return &LabelSpecification{
			LabelFormatType: LabelFormatTypeCommon2D,
			ImageType:       ImageTypePDF,
			LabelStockType:  &stockType,
		}

	}
	return &LabelSpecification{
		LabelFormatType: LabelFormatTypeCommon2D,
		ImageType:       ImageTypePNG,
	}
}

func (s *Shipment) DropoffType() string {
	if s.IsInternational() {
		return "BUSINESS_SERVICE_CENTER"
	}
	return "REGULAR_PICKUP"
}

func (s *Shipment) Weight() Weight {
	commoditiesWeight := s.Commodities.Weight()
	if !commoditiesWeight.IsZero() && s.IsInternational() {
		// Add a little extra weight to the entire shipment weight, since adding floats
		// in golang sometimes results in a float that is a little less than the actual
		// sum, and then the FedEx API will return an error
		commoditiesWeight.Value += 0.5
		return commoditiesWeight
	}

	switch s.ServiceType() {
	case ServiceTypeSmartPost:
		return Weight{Units: WeightUnitsLB, Value: 0.99}
	default:
		return Weight{Units: WeightUnitsLB, Value: 2}
	}
}

func (s *Shipment) ValidatedDimensions() Dimensions {
	if s.Dimensions.IsValid() {
		return s.Dimensions
	}

	switch s.ServiceType() {
	case ServiceTypeSmartPost:
		return Dimensions{Length: 6, Width: 5, Height: 5, Units: DimensionsUnitsIn}
	default:
		return Dimensions{Length: 10, Width: 5, Height: 5, Units: DimensionsUnitsIn}
	}
}

func (s *Shipment) SpecialServicesRequested() *SpecialServicesRequested {
	var (
		specialServiceTypes []string

		etdDetail               *EtdDetail
		eventNotificationDetail *EventNotificationDetail
		returnShipmentDetail    *ReturnShipmentDetail
	)

	if s.ServiceType() == ServiceTypeSmartPost {
		specialServiceTypes = append(specialServiceTypes, SpecialServiceTypeReturnShipment)
		returnShipmentDetail = &ReturnShipmentDetail{
			ReturnType: ReturnTypePrintReturnLabel,
		}
	}

	if s.IsInternational() {
		specialServiceTypes = append(specialServiceTypes, SpecialServiceTypeElectronicTradeDocuments)
		etdDetail = &EtdDetail{
			RequestedDocumentCopies: DocumentTypeCommercialInvoice,
		}
	}

	if s.NotificationEmail != "" {
		specialServiceTypes = append(specialServiceTypes, "EVENT_NOTIFICATION")
		eventNotificationDetail = defaultEventNotificationDetail(s.NotificationEmail)
	}

	if len(specialServiceTypes) == 0 {
		return nil
	}
	return &SpecialServicesRequested{
		SpecialServiceTypes: specialServiceTypes,

		EtdDetail:               etdDetail,
		EventNotificationDetail: eventNotificationDetail,
		ReturnShipmentDetail:    returnShipmentDetail,
	}
}

func (s *Shipment) CustomerReferences() []CustomerReference {
	customerReferences := make([]CustomerReference, len(s.References))
	for idx, reference := range s.References {
		switch s.ServiceType() {
		case ServiceTypeSmartPost:
			customerReferences[idx] = CustomerReference{
				CustomerReferenceType: CustomerReferenceTypeRMAAssociation,
				Value:                 sanitizeReferenceForFedexAPI(reference),
			}
		default:
			customerReferences[idx] = CustomerReference{
				CustomerReferenceType: CustomerReferenceTypeCustomerReference,
				Value:                 sanitizeReferenceForFedexAPI(reference),
			}
		}
	}

	if s.InvoiceNumber != "" {
		customerReferences = append(customerReferences, CustomerReference{
			CustomerReferenceType: CustomerReferenceTypeInvoice,
			Value:                 s.InvoiceNumber,
		})
	}
	return customerReferences
}

func sanitizeReferenceForFedexAPI(reference string) string {

	// Remove non-alphanumeric chars
	validatedReference := nonAlphanumericRegex.ReplaceAllString(reference, "")

	// Trim length
	if len(validatedReference) > 20 {
		validatedReference = validatedReference[0:20]
	}

	return validatedReference
}

func defaultEventNotificationDetail(notificationEmail string) *EventNotificationDetail {
	return &EventNotificationDetail{
		AggregationType: AggregationTypePerShipment,
		EventNotifications: []EventNotification{{
			Role: RoleShipper,
			Events: []string{
				"ON_DELIVERY",
				"ON_ESTIMATED_DELIVERY",
				"ON_EXCEPTION",
				"ON_SHIPMENT",
				"ON_TENDER",
			},
			NotificationDetail: NotificationDetail{
				NotificationType: "EMAIL",
				EmailDetail: EmailDetail{
					EmailAddress: notificationEmail,
					Name:         "Happy Returns dev team",
				},
				Localization: Localization{
					LanguageCode: "en",
				},
			},
			FormatSpecification: FormatSpecification{
				Type: "HTML",
			},
		}},
	}
}

func (s *Shipment) RequestedPackageLineItems() []RequestedPackageLineItem {
	return []RequestedPackageLineItem{{
		SequenceNumber:     1,
		PhysicalPackaging:  PackagingBag,
		ItemDescription:    "ItemDescription",
		CustomerReferences: s.CustomerReferences(),
		Weight:             s.Weight(),
		Dimensions:         s.ValidatedDimensions(),
	}}
}

type ProcessShipmentBody struct {
	ProcessShipmentRequest ProcessShipmentRequest `xml:"q0:ProcessShipmentRequest"`
}

type ProcessShipmentRequest struct {
	Request
	RequestedShipment RequestedShipment `xml:"q0:RequestedShipment"`
}

type ShipResponseEnvelope struct {
	Reply ProcessShipmentReply `xml:"Body>ProcessShipmentReply"`
}

func (s *ShipResponseEnvelope) Error() error {
	return s.Reply.Error()
}

// ProcessShipReply : Process shipment reply root (`xml:"Body>ProcessShipmentReply"`)
type ProcessShipmentReply struct {
	Reply
	TransactionDetail       TransactionDetail
	CompletedShipmentDetail CompletedShipmentDetail
	Events                  []Event
}

func (p *ProcessShipmentReply) LabelDataAndImageType() ([]byte, string, error) {
	if label := p.CompletedShipmentDetail.CompletedPackageDetails.Label; len(label.Parts) > 0 {
		return []byte(label.Parts[0].Image), label.ImageType, nil
	}
	return nil, "", errors.New("no label")
}

func (p *ProcessShipmentReply) CommercialInvoiceDataAndImageType() ([]byte, string, error) {
	for _, document := range p.CompletedShipmentDetail.ShipmentDocuments {
		if document.Type == DocumentTypeCommercialInvoice && len(document.Parts) > 0 {
			return []byte(document.Parts[0].Image), document.ImageType, nil
		}
	}
	return nil, "", errors.New("no commercial invoice")
}
