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
	rateRequestTypes := models.RequestTypePreferred
	packageCount := 1

	// When the service type is smartpost, getting rates from FedEx API doesn't
	// work
	serviceType := rate.ServiceType()
	weight := rate.Weight()

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
					ShipTimestamp:     models.Timestamp(time.Now()),
					DropoffType:       models.DropoffTypeRegularPickup,
					ServiceType:       serviceType,
					PackagingType:     models.PackagingTypeYourPackaging,
					PreferredCurrency: models.PreferredCurrencyUSD,
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
						PaymentType: models.PaymentTypeSender,
						Payor: models.Payor{
							ResponsibleParty: models.Shipper{
								AccountNumber: a.Account,
							},
						},
					},
					SpecialServicesRequested: rate.SpecialServicesRequested(),
					SmartPostDetail:          a.SmartPostDetail(serviceType),
					LabelSpecification: &models.LabelSpecification{
						LabelFormatType: models.LabelFormatTypeCommon2D,
						ImageType:       models.ImageTypePDF,
					},
					RateRequestTypes: &rateRequestTypes,
					PackageCount:     &packageCount,
					RequestedPackageLineItems: []models.RequestedPackageLineItem{
						{
							SequenceNumber:    1,
							GroupPackageCount: 1,
							Weight:            weight,
							Dimensions: models.Dimensions{
								Length: 5,
								Width:  5,
								Height: 5,
								Units:  models.DimensionsUnitsIn,
							},
							PhysicalPackaging: models.PackagingBag,
							ItemDescription:   "Stuff",
							CustomerReferences: []models.CustomerReference{
								{
									CustomerReferenceType: models.CustomerReferenceTypeCustomerReference,
									Value:                 models.CustomerReferenceValueNaftaCoo,
								},
							},
						},
					},
				},
			},
		},
	}
}
