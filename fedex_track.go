// History: Nov 20 13 tcolar Creation

package fedex

import (
	"fmt"
)

// Track by Tracking number
func trackRequest(fedex Fedex, body string) string {
	return fedex.wrapSoapRequest(fmt.Sprintf(`
		<v16:TrackRequest>
			%s
			%s
			<v16:ProcessingOptions>INCLUDE_DETAILED_SCANS</v16:ProcessingOptions>
		</v16:TrackRequest>
	`, fedex.soapCreds(), body))
}

func soapNumberTracking(fedex Fedex, carrierCode string, trackingNo string) string {
	return trackRequest(fedex, fmt.Sprintf(`
		<v16:SelectionDetails>
			<v16:CarrierCode>%s</v16:CarrierCode>
			<v16:PackageIdentifier>
				<v16:Type>TRACKING_NUMBER_OR_DOORTAG</v16:Type>
				<v16:Value>%s</v16:Value>
			</v16:PackageIdentifier>
		</v16:SelectionDetails>
	`, carrierCode, trackingNo))
}

// Track by PO/Zip
func soapPoTracking(fedex Fedex, carrierCode string, po string,
	postalCode string, countryCode string) string {
	return trackRequest(fedex, fmt.Sprintf(`
		<v16:SelectionDetails>
			<v16:CarrierCode>%s</v16:CarrierCode>
			<v16:PackageIdentifier>
				<v16:Type>PURCHASE_ORDER</v16:Type>
				<v16:Value>%s</v16:Value>
			</v16:PackageIdentifier>
			<v16:Destination>
				<v16:PostalCode>%s</v16:PostalCode>
				<v16:CountryCode>%s</v16:CountryCode>
			</v16:Destination>
		</v16:SelectionDetails>
	`, carrierCode, po, postalCode, countryCode))
}

// Track by ShipperRef / ShipperAccount
func soapRefTracking(fedex Fedex, carrierCode string, ref string,
	shipAccount string) string {
	return trackRequest(fedex, fmt.Sprintf(`
		<v16:SelectionDetails>
			<v16:CarrierCode>%s</v16:CarrierCode>
			<v16:PackageIdentifier>
				<v16:Type>SHIPPER_REFERENCE</v16:Type>
				<v16:Value>%s</v16:Value>
			</v16:PackageIdentifier>
			<v16:ShipmentAccountNumber>%s</v16:ShipmentAccountNumber>
		</v16:SelectionDetails>
	`, carrierCode, ref, shipAccount))
}
