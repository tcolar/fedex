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
					ServiceType:   shipment.ServiceType(),
					PackagingType: "YOUR_PACKAGING",
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
						PaymentType: "SENDER",
						Payor: models.Payor{
							ResponsibleParty: models.Shipper{
								AccountNumber: a.Account,
							},
						},
					},
					SmartPostDetail:               a.SmartPostDetail(shipment),
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

func (a API) SmartPostDetail(shipment *models.Shipment) *models.SmartPostDetail {
	if shipment.ServiceType() == "SMART_POST" {
		return &models.SmartPostDetail{
			Indicia:              "PARCEL_RETURN",
			AncillaryEndorsement: "ADDRESS_CORRECTION",
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

	dutiesPayment := models.Payment{
		PaymentType: "SENDER",
		Payor: models.Payor{
			ResponsibleParty: models.Shipper{
				AccountNumber: a.Account,
			},
		},
	}
	if shipment.Importer != "" {
		dutiesPayment = models.Payment{
			PaymentType: "THIRD_PARTY",
			Payor: models.Payor{
				ResponsibleParty: models.Shipper{
					AccountNumber: a.Account,
					Contact: models.Contact{
						CompanyName: fmt.Sprintf("Importer - %s", shipment.Importer),
					},
				},
			},
		}
	}

	return &models.CustomsClearanceDetail{
		ImporterOfRecord: models.Shipper{
			Contact: models.Contact{
				CompanyName: shipment.Importer,
			},
		},
		DutiesPayment:                  dutiesPayment,
		CustomsValue:                   &customsValue,
		Commodities:                    shipment.Commodities,
		PartiesToTransactionAreRelated: false,
		CommercialInvoice: &models.CommercialInvoice{
			Purpose:        "REPAIR_AND_RETURN",
			OriginatorName: shipment.OriginatorName,
		},
	}, nil
}
