// History: Nov 20 13 tcolar Creation

// Package fedex provides access to () FedEx Soap API's and unmarshal answers into Go structures
package fedex

import (
	"errors"
	"fmt"

	"github.com/happyreturns/fedex/api"
	"github.com/happyreturns/fedex/models"
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

// CreatePickup creates a pickup
func (f Fedex) CreatePickup(pickup *models.Pickup) (*models.CreatePickupReply, error) {
	var (
		reply *models.CreatePickupReply
		err   error
	)

	for delay := 0; delay <= 5; delay++ {
		reply, err = f.API.CreatePickup(pickup, delay)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("api create pickup: %s", err)
	}
	return reply, nil
}

func (f Fedex) Ship(shipment *models.Shipment) (*models.ProcessShipmentReply, error) {
	if f.API.HubID != "" && shipment.IsInternational() {
		return nil, errors.New("do not ship internationally with smartpost")
	}

	reply, err := f.API.ProcessShipment(shipment)
	if err != nil {
		return nil, fmt.Errorf("api process shipment: %s", err)
	}

	return reply, nil
}
