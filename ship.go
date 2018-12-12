package fedex

import (
	"time"

	"github.com/happyreturns/fedex/models"
)

func (f Fedex) shipmentEnvelope(shipmentType string, fromLocation, toLocation models.Address, fromContact, toContact models.Contact, smartPostKey string) (models.Envelope, error) {
	var serviceType string
	var weight models.Weight
	var dimensions models.Dimensions
	var smartPostDetail *models.SmartPostDetail
	var specialServicesRequested *models.SpecialServicesRequested

	switch shipmentType {
	case "SMART_POST":
		serviceType = "SMART_POST"
		weight = models.Weight{
			Units: "LB",
			Value: 0.99,
		}
		dimensions = models.Dimensions{
			Length: 6,
			Width:  5,
			Height: 5,
			Units:  "IN",
		}

		smartPostDetail = &models.SmartPostDetail{
			Indicia:              "PARCEL_RETURN",
			AncillaryEndorsement: "ADDRESS_CORRECTION",
			HubID:                f.HubID,
		}
		specialServicesRequested = &models.SpecialServicesRequested{
			SpecialServiceTypes: []string{"RETURN_SHIPMENT"},
			ReturnShipmentDetail: models.ReturnShipmentDetail{
				ReturnType: "PRINT_RETURN_LABEL",
			},
		}
	default:
		serviceType = "FEDEX_GROUND"
		weight = models.Weight{
			Units: "LB",
			Value: 13,
		}
		dimensions = models.Dimensions{
			Length: 13,
			Width:  13,
			Height: 13,
			Units:  "IN",
		}
	}

	req := models.ProcessShipmentRequest{
		Request: models.Request{
			WebAuthenticationDetail: models.WebAuthenticationDetail{
				UserCredential: models.UserCredential{
					Key:      f.Key,
					Password: f.Password,
				},
			},
			ClientDetail: models.ClientDetail{
				AccountNumber: f.Account,
				MeterNumber:   f.Meter,
			},
			Version: models.Version{
				ServiceID: "ship",
				Major:     23,
			},
		},
		RequestedShipment: models.RequestedShipment{
			ShipTimestamp: models.Timestamp(time.Now()),
			DropoffType:   "REGULAR_PICKUP",
			ServiceType:   serviceType,
			PackagingType: "YOUR_PACKAGING",
			Shipper: models.Shipper{
				AccountNumber: f.Account,
				Address:       fromLocation,
				Contact:       fromContact,
			},
			Recipient: models.Shipper{
				AccountNumber: f.Account,
				Address:       toLocation,
				Contact:       toContact,
			},
			ShippingChargesPayment: models.Payment{
				PaymentType: "SENDER",
				Payor: models.Payor{
					ResponsibleParty: models.ResponsibleParty{
						AccountNumber: f.Account,
					},
				},
			},
			SmartPostDetail:          smartPostDetail,
			SpecialServicesRequested: specialServicesRequested,
			LabelSpecification: models.LabelSpecification{
				LabelFormatType: "COMMON2D",
				ImageType:       "PDF",
			},
			RateRequestTypes: "LIST",
			PackageCount:     1,
			RequestedPackageLineItems: []models.RequestedPackageLineItem{
				{
					SequenceNumber:    1,
					PhysicalPackaging: "BAG",
					ItemDescription:   "Stuff",
					CustomerReferences: []models.CustomerReference{
						{
							CustomerReferenceType: "CUSTOMER_REFERENCE",
							Value: "NAFTA_COO",
						},
					},
					Weight:     weight,
					Dimensions: dimensions,
				},
			},
		},
	}

	return models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: "http://fedex.com/ws/ship/v23",
		Body: struct {
			ProcessShipmentRequest models.ProcessShipmentRequest `xml:"q0:ProcessShipmentRequest"`
		}{
			ProcessShipmentRequest: req,
		},
	}, nil
}
