package fedex

// The original source code used string formatting to create requests. We are
// now marshalling go structs for our requests, so moving code using the old
// way of creating requests here.

import (
	"encoding/xml"
	"fmt"

	"github.com/happyreturns/fedex/models"
)

// TrackByShipperRef : Return tracking info for a specific shipper reference
// ShipperRef is usually an order ID or other unique identifier
// ShipperAccountNumber is the Fedex account number of the shipper
func (f Fedex) TrackByShipperRef(carrierCode string, shipperRef string,
	shipperAccountNumber string) (reply models.TrackReply, err error) {
	reqXML := soapRefTracking(f, carrierCode, shipperRef, shipperAccountNumber)
	content, err := f.postXML(f.FedexURL+"/trck", reqXML)
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
	content, err := f.postXML(f.FedexURL+"/trck", reqXML)
	if err != nil {
		return reply, err
	}
	return f.parseTrackReply(content)
}

// Unmarshal XML SOAP response into a TrackReply
func (f Fedex) parseTrackReply(xmlResp []byte) (reply models.TrackReply, err error) {
	data := struct {
		Reply models.TrackReply `xml:"Body>TrackReply"`
	}{}
	err = xml.Unmarshal(xmlResp, &data)
	return data.Reply, err
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

// Track by Tracking number
func trackRequest(fedex Fedex, body string) string {
	return fedex.wrapSoapRequest(fmt.Sprintf(`
		<q0:TrackRequest>
			%s
			%s
			<q0:ProcessingOptions>INCLUDE_DETAILED_SCANS</q0:ProcessingOptions>
		</q0:TrackRequest>
	`, fedex.soapCreds("trck", "16"), body), "http://fedex.com/ws/track/v16")
}

// Track by PO/Zip
func soapPoTracking(fedex Fedex, carrierCode string, po string,
	postalCode string, countryCode string) string {
	return trackRequest(fedex, fmt.Sprintf(`
		<q0:SelectionDetails>
			<q0:CarrierCode>%s</q0:CarrierCode>
			<q0:PackageIdentifier>
				<q0:Type>PURCHASE_ORDER</q0:Type>
				<q0:Value>%s</q0:Value>
			</q0:PackageIdentifier>
			<q0:Destination>
				<q0:PostalCode>%s</q0:PostalCode>
				<q0:CountryCode>%s</q0:CountryCode>
			</q0:Destination>
		</q0:SelectionDetails>
	`, carrierCode, po, postalCode, countryCode))
}

// Track by ShipperRef / ShipperAccount
func soapRefTracking(fedex Fedex, carrierCode string, ref string,
	shipAccount string) string {
	return trackRequest(fedex, fmt.Sprintf(`
		<q0:SelectionDetails>
			<q0:CarrierCode>%s</q0:CarrierCode>
			<q0:PackageIdentifier>
				<q0:Type>SHIPPER_REFERENCE</q0:Type>
				<q0:Value>%s</q0:Value>
			</q0:PackageIdentifier>
			<q0:ShipmentAccountNumber>%s</q0:ShipmentAccountNumber>
		</q0:SelectionDetails>
	`, carrierCode, ref, shipAccount))
}
