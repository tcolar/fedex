// History: Nov 20 13 tcolar Creation

// fedex provides access to (some) FedEx Soap API's and unmarshall answers into Go structures
package fedex

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/happyreturns/fedex/models"
)

const (
	// Convenience constants for standard Fedex API url's
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

	SmartPostKey      string
	SmartPostPassword string
	SmartPostAccount  string
	SmartPostMeter    string
	SmartPostHubID    string

	FedexURL string
}

func (f Fedex) wrapSoapRequest(body string, namespace string) string {
	return fmt.Sprintf(`
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:q0="%s">
		<soapenv:Body>
			%s
		</soapenv:Body>
	</soapenv:Envelope>
	`, namespace, body)
}

func (f Fedex) soapCreds(serviceID, majorVersion string) string {
	return fmt.Sprintf(`
		<q0:WebAuthenticationDetail>
			<q0:UserCredential>
				<q0:Key>%s</q0:Key>
				<q0:Password>%s</q0:Password>
			</q0:UserCredential>
		</q0:WebAuthenticationDetail>
		<q0:ClientDetail>
			<q0:AccountNumber>%s</q0:AccountNumber>
			<q0:MeterNumber>%s</q0:MeterNumber>
		</q0:ClientDetail>
		<q0:Version>
			<q0:ServiceId>%s</q0:ServiceId>
			<q0:Major>%s</q0:Major>
			<q0:Intermediate>0</q0:Intermediate>
			<q0:Minor>0</q0:Minor>
		</q0:Version>
	`, f.Key, f.Password, f.Account, f.Meter, serviceID, majorVersion)
}

// TrackByShipperRef : Return tracking info for a specific shipper reference
// ShipperRef is usually an order ID or other unique identifier
// ShipperAccountNumber is the Fedex account number of the shipper
func (f Fedex) TrackByShipperRef(carrierCode string, shipperRef string,
	shipperAccountNumber string) (reply models.TrackReply, err error) {
	reqXML := soapRefTracking(f, carrierCode, shipperRef, shipperAccountNumber)
	content, err := f.PostXML(f.FedexURL+"/trck", reqXML)
	if err != nil {
		return reply, err
	}
	return f.parseTrackReply(content)
}

// TrackByPo : Returns tracking info for a specific Purchase Order (often the OrderId)
// Note that Fedex requires the Destination Postal Code & country
//   to match when making PO queries
func (f Fedex) TrackByPo(carrierCode string, po string, postalCode string,
	countryCode string) (reply models.TrackReply, err error) {
	reqXML := soapPoTracking(f, carrierCode, po, postalCode, countryCode)
	content, err := f.PostXML(f.FedexURL+"/trck", reqXML)
	if err != nil {
		return reply, err
	}
	return f.parseTrackReply(content)
}

func (f Fedex) makeRequestAndUnmarshal(url string, request models.Envelope,
	response interface{}) error {
	// Create request body
	reqXML, err := xml.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshal request xml: %s", err)
	}

	// Post XML
	content, err := f.PostXML(f.FedexURL+url, string(reqXML))
	if err != nil {
		return fmt.Errorf("post xml: %s", err)
	}

	// Parse response
	err = xml.Unmarshal(content, response)
	if err != nil {
		return fmt.Errorf("parse xml: %s", err)
	}
	return nil
}

// TrackByNumber : Returns tracking info for a specific Fedex tracking number
func (f Fedex) TrackByNumber(carrierCode string, trackingNo string) (*models.TrackReply, error) {
	request := f.trackByNumberSOAPRequest(carrierCode, trackingNo)
	response := &models.TrackResponseEnvelope{}

	err := f.makeRequestAndUnmarshal("/trck", request, response)
	if err != nil {
		return nil, fmt.Errorf("make track request and unmarshal: %s", err)
	}
	return &response.Reply, nil
}

// ShipGround : Creates a ground shipment
func (f Fedex) ShipGround(fromAddress models.Address, toAddress models.Address,
	fromContact models.Contact, toContact models.Contact) (*models.ProcessShipmentReply, error) {

	request := f.shipGroundSOAPRequest(fromAddress, toAddress, fromContact, toContact)
	response := &models.ShipResponseEnvelope{}

	err := f.makeRequestAndUnmarshal("/ship/v23", request, response)
	if err != nil {
		return nil, fmt.Errorf("make ship ground request and unmarshal: %s", err)
	}
	return &response.Reply, nil
}

// ShipSmartPost : Creates a Smart Post return shipment
func (f Fedex) ShipSmartPost(fromAddress models.Address, toAddress models.Address,
	fromContact models.Contact, toContact models.Contact) (*models.ProcessShipmentReply, error) {

	request := f.shipSmartPostSOAPRequest(fromAddress, toAddress, fromContact, toContact)
	response := &models.ShipResponseEnvelope{}

	err := f.makeRequestAndUnmarshal("/ship/v23", request, response)
	if err != nil {
		return nil, fmt.Errorf("make ship smart post request and unmarshal: %s", err)
	}
	return &response.Reply, nil
}

// Rate : Gets the estimated rates for a shipment
func (f Fedex) Rate(fromAddress models.Address, toAddress models.Address,
	fromContact models.Contact, toContact models.Contact) (*models.RateReply, error) {

	request := f.rateSOAPRequest(fromAddress, toAddress, fromContact, toContact)
	response := &models.RateResponseEnvelope{}

	err := f.makeRequestAndUnmarshal("/rate/v24", request, response)
	if err != nil {
		return nil, fmt.Errorf("make rate request and unmarshal: %s", err)
	}
	return &response.Reply, nil
}

// TODO
func (f Fedex) CreatePickup(pickupLocation models.PickupLocation, toAddress models.Address) (*models.CreatePickupReply, error) {

	request := f.createPickupRequest(pickupLocation, toAddress)
	response := &models.CreatePickupResponseEnvelope{}

	err := f.makeRequestAndUnmarshal("/pickup/v17", request, response)
	if err != nil {
		return nil, fmt.Errorf("make create pickup request and unmarshal: %s", err)
	}
	return &response.Reply, nil
}

// Unmarshal XML SOAP response into a TrackReply
func (f Fedex) parseTrackReply(xmlResp []byte) (reply models.TrackReply, err error) {
	data := struct {
		Reply models.TrackReply `xml:"Body>TrackReply"`
	}{}
	err = xml.Unmarshal(xmlResp, &data)
	return data.Reply, err
}

// Post Xml to Fedex API and return response
func (f Fedex) PostXML(url string, xml string) (content []byte, err error) {
	resp, err := http.Post(url, "text/xml", strings.NewReader(xml))
	if err != nil {
		return content, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
