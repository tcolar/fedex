// History: Nov 20 13 tcolar Creation

// Package fedex provides access to () FedEx Soap API's and unmarshal answers into Go structures
package fedex

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/happyreturns/fedex/api"
	"github.com/happyreturns/fedex/models"
	log "github.com/sirupsen/logrus"
)

// Convenience constants for standard Fedex API URLs
const (
	FedexAPIURL               = "https://ws.fedex.com:443/web-services"
	FedexAPITestURL           = "https://wsbeta.fedex.com:443/web-services"
	CarrierCodeExpress        = "FDXE"
	CarrierCodeGround         = "FDXG"
	CarrierCodeFreight        = "FXFR"
	CarrierCodeSmartPost      = "FXSP"
	CarrierCodeCustomCritical = "FXCC"
)

// Fedex : Utility to retrieve data from Fedex API
// Bypassing painful proper SOAP implementation and just crafting minimal XML messages to get the data we need.
// Fedex WSDL docs here: http://images.fedex.com/us/developer/product/WebServices/MyWebHelp/DeveloperGuide2012.pdf
type Fedex struct {
	api.API
}

var laTimeZone *time.Location

func init() {
	var err error
	laTimeZone, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

// CreatePickup creates a pickup with retry logic to try pickups on the following days
func (f Fedex) CreatePickup(pickup *models.Pickup) (*models.PickupSuccess, error) {
	var (
		reply *models.CreatePickupReply
		err   error
	)

	for delay := 0; delay <= 5; delay++ {
		fields := log.Fields{"pickup": pickup}

		// Calculate pickup window, but just try the next window in case of error
		window, err := pickupTimeWindow(pickup.PickupLocation.Address, delay)
		if err != nil {
			log.WithFields(fields).Error("calculate pickup time", err)
			continue
		}
		fields["window"] = window

		reply, err = f.API.CreatePickup(pickup, window)
		switch err.(type) {
		case nil:
			fields["reply"] = reply
			log.WithFields(fields).Info("made pickup")
			return &models.PickupSuccess{
				ConfirmationNumber: reply.PickupConfirmationNumber,
				Window:             *window,
			}, nil

		case models.PickupAlreadyExistsError:
			fields["reply"] = reply
			log.WithFields(fields).Info("pickup already exists")
			return &models.PickupSuccess{
				Window: *window,
			}, nil

		default:
			fields["err"] = err
			log.WithFields(fields).Info("failed pickup")
		}
	}

	return nil, fmt.Errorf("fedex create pickup: %s", err)
}

func pickupTimeWindow(pickupAddress models.Address, numDaysToDelay int) (*models.PickupTimeWindow, error) {
	location, err := toLocation(pickupAddress)
	if err != nil {
		location = laTimeZone
	}

	readyTime := time.Now().In(location).Add(time.Duration(numDaysToDelay*24) * time.Hour)

	// If it's past the ready time of the current day, ship the next day, not today
	if readyTime.After(timeForReadyPickup(readyTime)) {
		readyTime = readyTime.Add(24 * time.Hour)
	}
	readyTime = timeForReadyPickup(readyTime)

	// Don't schedule pickups for Saturday or Sunday
	if readyTime.Weekday() == time.Saturday || readyTime.Weekday() == time.Sunday {
		return nil, fmt.Errorf("no pickups on saturday or sunday %d", numDaysToDelay)
	}

	return &models.PickupTimeWindow{
		ReadyTime: readyTime,
		CloseTime: readyTime.Add(8 * time.Hour),
	}, nil
}

func timeForReadyPickup(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 10, 45, 0, 0, t.Location())
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

func (f Fedex) Ship(shipment *models.Shipment) (*models.ProcessShipmentReply, error) {
	if f.isSmartPost() && shipment.IsInternational() {
		return nil, errors.New("do not ship internationally with smartpost")
	}

	// Don't use non-smartpost accounts for returns
	if !f.isSmartPost() {
		shipment.Service = "default"
	}

	reply, err := f.API.ProcessShipment(shipment)
	if err != nil {
		return nil, fmt.Errorf("api process shipment: %s", err)
	}

	return reply, nil
}

func (f Fedex) isSmartPost() bool {
	return f.API.HubID != ""
}
