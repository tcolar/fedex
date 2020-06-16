package api

import (
	"fmt"

	"github.com/happyreturns/fedex/models"
)

const (
	processShipmentVersion = "v23"
)

func (a API) ProcessShipment(shipment *models.Shipment) (*models.ProcessShipmentReply, error) {
	request, err := a.processShipmentRequest(shipment)
	if err != nil {
		return nil, fmt.Errorf("create process shipment request: %s", err)
	}

	endpoint := fmt.Sprintf("/ship/%s", processShipmentVersion)
	response := &models.ShipResponseEnvelope{}
	if err := a.makeRequestAndUnmarshalResponse(endpoint, request, response); err != nil {
		return nil, fmt.Errorf("make process shipment request and unmarshal: %s", err)
	}

	return &response.Reply, nil
}

func (a API) processShipmentRequest(shipment *models.Shipment) (*models.Envelope, error) {
	customsClearanceDetail, err := a.customsClearanceDetail(shipment)
	if err != nil {
		return nil, fmt.Errorf("customs clearance detail: %s", err)
	}

	packageCount := 1
	serviceType := shipment.ServiceType()

	return &models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: fmt.Sprintf("http://fedex.com/ws/ship/%s", processShipmentVersion),
		Body: models.ProcessShipmentBody{
			ProcessShipmentRequest: models.ProcessShipmentRequest{
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
					Version: models.Version{
						ServiceID: "ship",
						Major:     23,
					},
				},
				RequestedShipment: models.RequestedShipment{
					ShipTimestamp: models.Timestamp(shipment.ShipTime()),
					DropoffType:   shipment.DropoffType(),
					ServiceType:   serviceType,
					PackagingType: models.PackagingTypeYourPackaging,
					Shipper: models.Shipper{
						AccountNumber: a.Account,
						Address:       shipment.FromAddress,
						Contact:       shipment.FromContact,
					},
					Recipient: models.Shipper{
						AccountNumber: a.Account,
						Address:       shipment.ToAddress,
						Contact:       shipment.ToContact,
					},
					ShippingChargesPayment: &models.Payment{
						PaymentType: models.PaymentTypeSender,
						Payor: models.Payor{
							ResponsibleParty: models.Shipper{
								AccountNumber: a.Account,
							},
						},
					},
					SmartPostDetail:               a.SmartPostDetail(serviceType),
					SpecialServicesRequested:      shipment.SpecialServicesRequested(),
					CustomsClearanceDetail:        customsClearanceDetail,
					LabelSpecification:            shipment.LabelSpecification(),
					ShippingDocumentSpecification: shipment.ShippingDocumentSpecification(),
					PackageCount:                  &packageCount,
					RequestedPackageLineItems:     shipment.RequestedPackageLineItems(),
				},
			},
		},
	}, nil
}

func (a API) SmartPostDetail(serviceType string) *models.SmartPostDetail {
	if serviceType == models.ServiceTypeSmartPost {
		return &models.SmartPostDetail{
			Indicia:              models.IndiciaParcelReturn,
			AncillaryEndorsement: models.AncillaryEndorsementAddressCorrection,
			HubID:                a.HubID,
		}
	}
	return nil
}

func (a API) customsClearanceDetail(shipment *models.Shipment) (*models.CustomsClearanceDetail, error) {
	if !shipment.IsInternational() {
		return nil, nil
	}

	customsValue, err := shipment.Commodities.CustomsValue()
	if err != nil {
		return nil, fmt.Errorf("commodities customs value: %s", err)
	}

	importerOfRecord := models.Shipper{
		AccountNumber: a.Account,
		Contact: models.Contact{
			CompanyName: "Happy Returns",
			PhoneNumber: "424 325 9510",
		},
		Address: models.Address{
			StreetLines:         []string{"1106 Broadway"},
			City:                "Santa Monica",
			StateOrProvinceCode: "CA",
			PostalCode:          "90401",
			CountryCode:         "US",
		},
	}
	dutiesPayment := models.Payment{
		PaymentType: models.PaymentTypeRecipient,
		Payor: models.Payor{
			ResponsibleParty: importerOfRecord,
		},
	}

	return &models.CustomsClearanceDetail{
		Brokers: []models.Broker{{
			Type: models.BrokerTypeImport,
			Broker: models.Shipper{
				AccountNumber: a.Account,
				Contact: models.Contact{
					CompanyName: shipment.Broker(),
				},
			},
		}},
		ImporterOfRecord:               importerOfRecord,
		DutiesPayment:                  dutiesPayment,
		CustomsValue:                   &customsValue,
		Commodities:                    shipment.Commodities,
		PartiesToTransactionAreRelated: false,
		CommercialInvoice: &models.CommercialInvoice{
			Purpose:        models.CommercialInvoicePurposeRepairAndReturn,
			OriginatorName: shipment.OriginatorName,
		},
	}, nil
}
