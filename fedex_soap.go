// History: Nov 20 13 tcolar Creation

package fedex

import (
	"fmt"
)

// SOAP monkey patching

func soapPoTracking(fedex Fedex, carrierCode string, po string,
	postalCode string, countryCode string) string {
	return soapHead(fedex) +
		fmt.Sprintf(`
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
</q0:SelectionDetails>`, carrierCode, po, postalCode, countryCode) +
		FEDEX_SOAP_TAIL
}

func soapRefTracking(fedex Fedex, carrierCode string, ref string,
	shipAccount string) string {
	return soapHead(fedex) +
		fmt.Sprintf(`
<q0:SelectionDetails>
<q0:CarrierCode>%s</q0:CarrierCode>
  <q0:PackageIdentifier>
    <q0:Type>SHIPPER_REFERENCE</q0:Type>
    <q0:Value>%s</q0:Value>
  </q0:PackageIdentifier>
  <q0:ShipmentAccountNumber>%s</q0:ShipmentAccountNumber>
</q0:SelectionDetails>`, carrierCode, ref, shipAccount) +
		FEDEX_SOAP_TAIL
}

func soapHead(f Fedex) string {
	return fmt.Sprintf(FEDEX_SOAP_HEAD,
		f.Key, f.Password, f.Account, f.Meter)
}

const (
	FEDEX_SOAP_HEAD = `
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
</q0:Version>`

	FEDEX_SOAP_TAIL = `
</q0:TrackRequest>
</soapenv:Body></soapenv:Envelope>`
)
