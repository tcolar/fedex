package api

import (
	"fmt"
	"strings"

	"github.com/happyreturns/fedex/models"
)

const (
	createPickupVersion = "v17"
)

func (a API) CreatePickup(pickup *models.Pickup, window *models.PickupTimeWindow) (*models.CreatePickupReply, error) {
	request, err := a.createPickupRequest(pickup, window)
	if err != nil {
		return nil, fmt.Errorf("create pickup request: %s", err)
	}

	endpoint := fmt.Sprintf("/pickup/%s", createPickupVersion)
	response := &models.CreatePickupResponseEnvelope{}
	err = a.makeRequestAndUnmarshalResponse(endpoint, request, response)

	switch {
	case err != nil && strings.Contains(err.Error(), "pickup already exists"):
		return nil, models.PickupAlreadyExistsError{}
	case err != nil:
		return nil, fmt.Errorf("make create pickup request and unmarshal: %s", err)
	default:
		return &response.Reply, nil
	}
}

func (a API) createPickupRequest(pickup *models.Pickup, window *models.PickupTimeWindow) (*models.Envelope, error) {
	return &models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: fmt.Sprintf("http://fedex.com/ws/pickup/%s", createPickupVersion),
		Body: models.CreatePickupBody{
			CreatePickupRequest: models.CreatePickupRequest{
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
						ServiceID: "disp",
						Major:     17,
					},
				},
				OriginDetail: models.OriginDetail{
					UseAccountAddress:       false,
					PickupLocation:          pickup.PickupLocation,
					PackageLocation:         "NONE",
					BuildingPart:            "SUITE",
					BuildingPartDescription: "",
					ReadyTimestamp:          models.Timestamp(window.ReadyTime),
					CompanyCloseTime:        window.CloseTime.Format("15:04:05-07:00"),
				},
				FreightPickupDetail: models.FreightPickupDetail{
					ApprovedBy:  pickup.PickupLocation.Contact,
					Payment:     "SENDER",
					Role:        "SHIPPER",
					SubmittedBy: models.Contact{},
					LineItems: []models.FreightPickupLineItem{
						{
							Service:        "INTERNATIONAL_ECONOMY_FREIGHT",
							SequenceNumber: 1,
							Destination:    pickup.ToAddress,
							Packaging:      "BAG",
							Pieces:         1,
							Weight: models.Weight{
								Units: "LB",
								Value: 1,
							},
							TotalHandlingUnits: 1,
							JustOneMore:        false,
							Description:        "",
						},
					},
				},
				PackageCount:         1,
				CarrierCode:          "FDXG",
				Remarks:              "",
				CommodityDescription: "",
			},
		},
	}, nil
}
