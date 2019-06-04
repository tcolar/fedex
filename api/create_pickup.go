package api

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/happyreturns/fedex/models"
)

const (
	createPickupVersion = "v17"
)

var laTimeZone *time.Location

func init() {
	var err error
	laTimeZone, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}
}

func (a API) CreatePickup(pickup *models.Pickup, numDaysToDelay int) (*models.CreatePickupReply, error) {
	request, err := a.createPickupRequest(pickup, numDaysToDelay)
	if err != nil {
		return nil, fmt.Errorf("create pickup request: %s", err)
	}

	endpoint := fmt.Sprintf("/pickup/%s", createPickupVersion)
	response := &models.CreatePickupResponseEnvelope{}
	err = a.makeRequestAndUnmarshalResponse(endpoint, request, response)
	if err != nil {
		return nil, fmt.Errorf("make create pickup request and unmarshal: %s", err)
	}

	return &response.Reply, nil
}

func (a API) createPickupRequest(pickup *models.Pickup, numDaysToDelay int) (*models.Envelope, error) {

	pickupTime, err := calculatePickupTime(pickup.PickupLocation.Address, numDaysToDelay)
	if err != nil {
		return nil, fmt.Errorf("calculate pickup time: %s", err)
	}
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
					ReadyTimestamp:          models.Timestamp(pickupTime),
					CompanyCloseTime:        "16:00:00", // TODO not necessarily true
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

func calculatePickupTime(pickupAddress models.Address, numDaysToDelay int) (time.Time, error) {
	location, err := toLocation(pickupAddress)
	if err != nil {
		location = laTimeZone
	}

	pickupTime := time.Now().In(location).Add(time.Duration(numDaysToDelay*24) * time.Hour)

	// If it's past 12pm, ship the next day, not today
	if pickupTime.Hour() >= 12 {
		pickupTime = pickupTime.Add(24 * time.Hour)
	}

	// Don't schedule pickups for Saturday or Sunday
	if pickupTime.Weekday() == time.Saturday || pickupTime.Weekday() == time.Sunday {
		return time.Time{}, errors.New("no pickups on saturday or sunday")
	}

	year, month, day := pickupTime.Date()
	return time.Date(year, month, day, 12, 0, 0, 0, location), nil
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