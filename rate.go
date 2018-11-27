package fedex

import (
	"time"

	"github.com/happyreturns/fedex/models"
)

func (f Fedex) rateSOAPRequest(fromLocation, toLocation models.Address, fromContact, toContact models.Contact) models.Envelope {
	return models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: "http://fedex.com/ws/rate/v24",
		Body: struct {
			RateRequest models.RateRequest `xml:"q0:RateRequest"`
		}{
			RateRequest: models.RateRequest{
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
						ServiceID: "crs",
						Major:     24,
					},
				},
				RequestedShipment: models.RequestedShipment{
					ShipTimestamp: models.Timestamp(time.Now()),
					DropoffType:   "REGULAR_PICKUP",
					ServiceType:   "FEDEX_GROUND",
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
					LabelSpecification: models.LabelSpecification{
						LabelFormatType: "COMMON2D",
						ImageType:       "PDF",
					},
					RateRequestTypes: "LIST",
					PackageCount:     1,
					RequestedPackageLineItems: []models.RequestedPackageLineItem{
						{
							SequenceNumber:    1,
							GroupPackageCount: 1,
							Weight: models.Weight{
								Units: "LB",
								Value: 40,
							},
							Dimensions: models.Dimensions{
								Length: 5,
								Width:  5,
								Height: 5,
								Units:  "IN",
							},
							PhysicalPackaging: "BAG",
							ItemDescription:   "Stuff",
							CustomerReferences: []models.CustomerReference{
								{
									CustomerReferenceType: "CUSTOMER_REFERENCE",
									Value: "NAFTA_COO",
								},
							},
						},
					},
				},
			},
		},
	}
}
