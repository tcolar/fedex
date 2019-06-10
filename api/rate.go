package api

import (
	"fmt"
	"time"

	"github.com/happyreturns/fedex/models"
)

const (
	rateVersion = "v24"
)

func (a API) Rate(rate *models.Rate) (*models.RateReply, error) {

	endpoint := fmt.Sprintf("/rate/%s", rateVersion)
	request := a.rateRequest(rate)
	response := &models.RateResponseEnvelope{}

	err := a.makeRequestAndUnmarshalResponse(endpoint, request, response)
	if err != nil {
		return nil, fmt.Errorf("make rate request and unmarshal: %s", err)
	}

	return &response.Reply, nil
}

func (a API) rateRequest(rate *models.Rate) *models.Envelope {
	rateRequestTypes := "LIST"
	packageCount := 1
	serviceType := "FEDEX_GROUND"
	if rate.FromAddress.ShipsOutWithInternationalEconomy() {
		serviceType = "INTERNATIONAL_ECONOMY"
	}

	return &models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: fmt.Sprintf("http://fedex.com/ws/rate/%s", rateVersion),
		Body: models.RateBody{
			RateRequest: models.RateRequest{
				Request: models.Request{
					WebAuthenticationDetail: models.WebAuthenticationDetail{
						UserCredential: models.UserCredential{
							Key:      a.Key,
							Password: a.Password,
						},
					},
					ClientDetail: models.ClientDetail{
						AccountNumber: a.Account,
						MeterNumber:   a.Meter,
					},
					TransactionDetail: &models.TransactionDetail{
						CustomerTransactionID: "Rate Request",
					},
					Version: models.Version{
						ServiceID: "crs",
						Major:     24,
					},
				},
				RequestedShipment: models.RequestedShipment{
					ShipTimestamp: models.Timestamp(time.Now()),
					DropoffType:   "REGULAR_PICKUP",
					ServiceType:   serviceType,
					PackagingType: "YOUR_PACKAGING",
					Shipper: models.Shipper{
						AccountNumber: a.Account,
						Address:       rate.FromAndTo.FromAddress,
						Contact:       rate.FromAndTo.FromContact,
					},
					Recipient: models.Shipper{
						AccountNumber: a.Account,
						Address:       rate.FromAndTo.ToAddress,
						Contact:       rate.FromAndTo.ToContact,
					},
					ShippingChargesPayment: &models.Payment{
						PaymentType: "SENDER",
						Payor: models.Payor{
							ResponsibleParty: models.Shipper{
								AccountNumber: a.Account,
							},
						},
					},
					LabelSpecification: &models.LabelSpecification{
						LabelFormatType: "COMMON2D",
						ImageType:       "PDF",
					},
					RateRequestTypes: &rateRequestTypes,
					PackageCount:     &packageCount,
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
									Value:                 "NAFTA_COO",
								},
							},
						},
					},
				},
			},
		},
	}
}
