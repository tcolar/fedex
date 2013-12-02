// History: Nov 20 13 tcolar Creation

package fedex

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	// Convenience constants for standard Fedex API url's
	FEDEX_API_URL       = "https://ws.fedex.com:443/web-services"
	FEDEX_API_TEST_URL  = "https://wsbeta.fedex.com:443/web-services"
	FEDEX_TEST_TRACKING = "123456789012"
)

// Utility to retrieve data from Fedex API
// Bypassing painful proper SOAP implementation and just crafting minimal XML messages to get the data we need.
// Fedex WSDL docs here: http://images.fedex.com/us/developer/product/WebServices/MyWebHelp/DeveloperGuide2012.pdf
type Fedex struct {
	Key, Password, Account, Meter string
	FedexUrl                      string
}

// Return tracking info for a specific Fedex tracking number
func (f Fedex) TrackByNumber(carrierCode string, trackingNo string) (reply TrackReply, err error) {
	reqXml := soapNumberTracking(f, carrierCode, trackingNo)
	content, err := f.PostXml(f.FedexUrl+"/trck", reqXml)
	if err != nil {
		return reply, err
	}
	return f.ParseTrackReply(content)
}

// Return tracking info for a specific shipper reference
// ShipperRef is usually an order ID or other unique identifier
// ShipperAccountNumber is the Fedex account number of the shipper
func (f Fedex) TrackByShipperRef(carrierCode string, shipperRef string,
	shipperAccountNumber string) (reply TrackReply, err error) {
	reqXml := soapRefTracking(f, carrierCode, shipperRef, shipperAccountNumber)
	content, err := f.PostXml(f.FedexUrl+"/trck", reqXml)
	if err != nil {
		return reply, err
	}
	return f.ParseTrackReply(content)
}

// Return tracking info for a specific Purchase Order (often the OrderId)
// Note that Fedex requires the Destination Postal Code & country
//   to match when making PO queries
func (f Fedex) TrackByPo(carrierCode string, po string, postalCode string,
	countryCode string) (reply TrackReply, err error) {
	reqXml := soapPoTracking(f, carrierCode, po, postalCode, countryCode)
	content, err := f.PostXml(f.FedexUrl+"/trck", reqXml)
	if err != nil {
		return reply, err
	}
	return f.ParseTrackReply(content)
}

// Unmarshal XML SOAP response into a TrackReply
func (f Fedex) ParseTrackReply(xmlResp []byte) (reply TrackReply, err error) {
	data := struct {
		Reply TrackReply `xml:"Body>TrackReply"`
	}{}
	err = xml.Unmarshal(xmlResp, &data)
	return data.Reply, err
}

// Post Xml to Fedex API and return response
func (f Fedex) PostXml(url string, xml string) (content []byte, err error) {
	resp, err := http.Post(f.FedexUrl, "text/xml", strings.NewReader(xml))
	if err != nil {
		return content, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// Dump some of the query resuts as an example
func Dump(reply TrackReply) {
	log.Print(reply)
	// Dummy example of using the data
	log.Printf("Successs : %t", !reply.Failed())
	if !reply.Failed() {
		tracking := reply.CompletedTrackDetails[0].TrackDetails[0].TrackingNumber
		log.Printf("Tracking Number: %s", tracking)
		log.Print(reply.CompletedTrackDetails[0].TrackDetails[0].ActualDeliveryAddress)
	}
}
