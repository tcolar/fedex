package fedex

import (
	"fmt"
	"strings"
	"time"

	"github.com/happyreturns/fedex/models"
)

var laTimeZone *time.Location

func init() {
	var err error
	laTimeZone, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}
}

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
					PackageLocation:         "NONE",  // TODO not necessarily true
					BuildingPart:            "SUITE", // TODO not necessarily true
					BuildingPartDescription: "",
					ReadyTimestamp:          models.Timestamp(f.pickupTime(pickupLocation.Address)),
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

func (f Fedex) pickupTime(pickupAddress models.Address) time.Time {
	location, err := toLocation(pickupAddress)
	if err != nil {
		location = laTimeZone
	}

	now := time.Now().In(location)
	year, month, day := now.Date()

	// If it's past 9am, ship the next day, not today
	if now.Hour() > 9 {
		day++
	}

	return time.Date(year, month, day, 9, 0, 0, 0, location)
}

// toLocation attempts to return the timezone based on state, returning los
// angeles if unable to
func toLocation(pickupAddress models.Address) (*time.Location, error) {
	tzDatabaseName := ""
	switch strings.ToUpper(pickupAddress.StateOrProvinceCode) {
	case "AK":
		tzDatabaseName = "America/Anchorage"
	case "HI":
		tzDatabaseName = "Pacific/Honolulu"
	case "AL", "AR", "IL", "IA", "KS", "KY", "LA", "MN", "MS", "MO", "NE", "ND", "OK", "SD", "TN", "TX", "WI":
		tzDatabaseName = "America/Chicago"
	case "AZ", "CO", "ID", "MT", "NM", "UT", "WY":
		tzDatabaseName = "America/Denver"
	case "CT", "DE", "FL", "GA", "IN", "ME", "MD", "MA", "MI", "NH", "NJ", "NY", "NC", "OH", "PA", "RI", "SC", "VT", "VA", "WV":
		tzDatabaseName = "America/New_York"
	default:
		return laTimeZone, nil
	}

	timeZone, err := time.LoadLocation(tzDatabaseName)
	if err != nil {
		return nil, fmt.Errorf("load location from time zone %s: %s", tzDatabaseName, err)
	}
	return timeZone, nil
}
