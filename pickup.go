package fedex

import (
	"time"

	"github.com/happyreturns/fedex/models"
)

func (f Fedex) createPickupRequest(pickupLocation models.PickupLocation, toAddress models.Address) models.Envelope {
	return models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: "http://fedex.com/ws/pickup/v17",
		Body: struct {
			CreatePickupRequest models.CreatePickupRequest `xml:"q0:CreatePickupRequest"`
		}{
			CreatePickupRequest: models.CreatePickupRequest{
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
						ServiceID: "disp",
						Major:     17,
					},
				},
				OriginDetail: models.OriginDetail{
					UseAccountAddress:       false,
					PickupLocation:          pickupLocation,
					PackageLocation:         "NONE",     // TODO not necessarily true
					BuildingPart:            "BUILDING", // TODO not necessarily true
					BuildingPartDescription: "",
					ReadyTimestamp:          models.Timestamp(f.pickupTime()),
					CompanyCloseTime:        "16:00:00", // TODO not necessarily true
				},
				FreightPickupDetail: models.FreightPickupDetail{
					ApprovedBy:  pickupLocation.Contact,
					Payment:     "SENDER",
					Role:        "SHIPPER",
					SubmittedBy: models.Contact{},
					LineItems: []models.FreightPickupLineItem{
						{
							Service:        "INTERNATIONAL_ECONOMY_FREIGHT",
							SequenceNumber: 1,
							Destination:    toAddress,
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
				CarrierCode:          "FDXE",
				Remarks:              "",
				CommodityDescription: "",
			},
		},
	}
}

func (f Fedex) pickupTime() time.Time {
	now := time.Now()
	year, month, day := now.Date()

	// If it's past 9am, ship the next day, not today
	if now.Hour() > 9 {
		day++
	}

	return time.Date(year, month, day, 9, 0, 0, 0, now.Location())
}
