// History: Nov 20 13 tcolar Creation

// Package fedex provides access to (some) FedEx Soap API's and unmarshal answers into Go structures
package fedex

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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
	Key      string
	Password string
	Account  string
	Meter    string
	HubID    string // for SmartPost

	FedexURL string
}

// Shipment wraps all the Fedex API fields needed for creating a shipment
type Shipment struct {
	FromAddress       models.Address
	ToAddress         models.Address
	FromContact       models.Contact
	ToContact         models.Contact
	NotificationEmail string
	Reference         string
}

// TrackByNumber returns tracking info for a specific Fedex tracking number
func (f Fedex) TrackByNumber(carrierCode string, trackingNo string) (*models.TrackReply, error) {

	request := f.trackByNumberRequest(carrierCode, trackingNo)
	response := &models.TrackResponseEnvelope{}

	err := f.makeRequestAndUnmarshalResponse("/trck", request, response)
	if err != nil {
		return nil, fmt.Errorf("make track request and unmarshal: %s", err)
	}
	return &response.Reply, nil
}

// ShipGround creates a ground shipment
func (f Fedex) ShipGround(shipment *Shipment) (*models.ProcessShipmentReply, error) {

	request, err := f.shipmentRequest("FEDEX_GROUND", shipment)
	if err != nil {
		return nil, fmt.Errorf("create shipment request: %s", err)
	}

	response := &models.ShipResponseEnvelope{}
	err = f.makeRequestAndUnmarshalResponse("/ship/v23", request, response)
	if err != nil {
		return nil, fmt.Errorf("make ship ground request and unmarshal: %s", err)
	}
	return &response.Reply, nil
}

// ShipSmartPost creates a Smart Post return shipment
func (f Fedex) ShipSmartPost(shipment *Shipment) (*models.ProcessShipmentReply, error) {

	request, err := f.shipmentRequest("SMART_POST", shipment)
	if err != nil {
		return nil, fmt.Errorf("create shipment request: %s", err)
	}

	response := &models.ShipResponseEnvelope{}
	err = f.makeRequestAndUnmarshalResponse("/ship/v23", request, response)
	if err != nil {
		return nil, fmt.Errorf("make ship smart post request and unmarshal: %s", err)
	}

	return &response.Reply, nil
}

// Rate : Gets the estimated rates for a shipment
func (f Fedex) Rate(fromAddress models.Address, toAddress models.Address,
	fromContact models.Contact, toContact models.Contact) (*models.RateReply, error) {

	request := f.rateRequest(fromAddress, toAddress, fromContact, toContact)
	response := &models.RateResponseEnvelope{}

	err := f.makeRequestAndUnmarshalResponse("/rate/v24", request, response)
	if err != nil {
		return nil, fmt.Errorf("make rate request and unmarshal: %s", err)
	}

	return &response.Reply, nil
}

// CreatePickup creates a pickup
func (f Fedex) CreatePickup(pickupLocation models.PickupLocation, toAddress models.Address) (*models.CreatePickupReply, error) {

	request := f.createPickupRequest(pickupLocation, toAddress)
	response := &models.CreatePickupResponseEnvelope{}

	err := f.makeRequestAndUnmarshalResponse("/pickup/v17", request, response)
	if err != nil {
		return nil, fmt.Errorf("make create pickup request and unmarshal: %s", err)
	}

	return &response.Reply, nil
}

// SendNotifications gets notifications sent to an email
func (f Fedex) SendNotifications(trackingNo, email string) (*models.SendNotificationsReply, error) {

	request := f.notificationsRequest(trackingNo, email)
	response := &models.SendNotificationsResponseEnvelope{}

	err := f.makeRequestAndUnmarshalResponse("/track/v16", request, response)
	if err != nil {
		return nil, fmt.Errorf("make send notifications request: %s", err)
	}

	return &response.Reply, nil
}

func (f Fedex) makeRequestAndUnmarshalResponse(url string, request models.Envelope,
	response models.Response) error {
	// Create request body
	reqXML, err := xml.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshal request xml: %s", err)
	}

	// Post XML
	content, err := f.postXML(f.FedexURL+url, string(reqXML))
	if err != nil {
		return fmt.Errorf("post xml: %s", err)
	}

	// Parse response
	err = xml.Unmarshal(content, response)
	if err != nil {
		return fmt.Errorf("parse xml: %s", err)
	}

	// Check if reply failed (FedEx responds with 200 even though it failed)
	err = response.Error()
	if err != nil {
		return fmt.Errorf("response error: %s", err)
	}

	return nil
}

// postXML to Fedex API and return response
func (f Fedex) postXML(url string, xml string) (content []byte, err error) {
	resp, err := http.Post(url, "text/xml", strings.NewReader(xml))
	if err != nil {
		return content, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
