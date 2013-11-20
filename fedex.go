// History: Nov 20 13 tcolar Creation

package fedex

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Utility to retrieve data from Fedex API
// Bypassing painful proper SOAP implementation and just crafting minimal XML messages to get the data we need.
type Fedex struct {
	Key, Password, Account, Meter string
	FedexUrl                      string
}

// Return tracking info for a specific shipper reference
// ShipperRef is usually an order ID or other unique identifier
// ShipperAccountNumber is the Fedex account number of the shipper
func (f Fedex) TrackByShipperRef(shipperRef string, shipperAccountNumber string) (reply TrackReply, err error) {
	reqXml := fmt.Sprintf(FEDEX_XML_BY_REFERENCE,
		f.Key, f.Password, f.Account, f.Meter, shipperRef, shipperAccountNumber)
	content, err := f.PostXml(f.FedexUrl+"/trck", reqXml)
	if err != nil {
		return reply, err
	}
	return f.ParseTrackReply(content)
}

// Return tracking info for a specific Purchase Order (often the OrderId)
// Note that Fedex requires the Destination Postal Code & country
//   to match when making PO queries
func (f Fedex) TrackByPo(po string, postalCode string, countryCode string) (reply TrackReply, err error) {
	reqXml := fmt.Sprintf(FEDEX_XML_BY_PO,
		f.Key, f.Password, f.Account, f.Meter, po, postalCode, countryCode)
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

const (
	// Convenience constants for standard Fedex API url's
	FEDEX_API_URL       = "https://ws.fedex.com:443/web-services"
	FEDEX_API_TEST_URL  = "https://wsbeta.fedex.com:443/web-services"
	FEDEX_TEST_TRACKING = "123456789012"

	// XML templates
	FEDEX_XML_BY_REFERENCE = `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:q0="http://fedex.com/ws/track/v7" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
<soapenv:Body>
<q0:TrackRequest>
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
<q0:ServiceId>trck</q0:ServiceId>
<q0:Major>7</q0:Major>
<q0:Intermediate>0</q0:Intermediate>
<q0:Minor>0</q0:Minor>
</q0:Version>
<q0:SelectionDetails>
<q0:CarrierCode>FDXE</q0:CarrierCode>
  <q0:PackageIdentifier>
    <q0:Type>SHIPPER_REFERENCE</q0:Type>
    <q0:Value>%s</q0:Value>
  </q0:PackageIdentifier>
  <q0:ShipmentAccountNumber>%s</q0:ShipmentAccountNumber>
</q0:SelectionDetails>
</q0:TrackRequest>
</soapenv:Body></soapenv:Envelope>`

	FEDEX_XML_BY_PO = `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:q0="http://fedex.com/ws/track/v7" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
<soapenv:Body>
<q0:TrackRequest>
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
<q0:ServiceId>trck</q0:ServiceId>
<q0:Major>7</q0:Major>
<q0:Intermediate>0</q0:Intermediate>
<q0:Minor>0</q0:Minor>
</q0:Version>
<q0:SelectionDetails>
<q0:CarrierCode>FDXE</q0:CarrierCode>
  <q0:PackageIdentifier>
    <q0:Type>PURCHASE_ORDER</q0:Type>
    <q0:Value>%s</q0:Value>
  </q0:PackageIdentifier>
  <q0:Destination>
    <q0:PostalCode>%s</q0:PostalCode>
    <q0:CountryCode>%s</q0:CountryCode>
  </q0:Destination>
</q0:SelectionDetails>
</q0:TrackRequest>
</soapenv:Body></soapenv:Envelope>`
)
