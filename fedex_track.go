// History: Nov 20 13 tcolar Creation

package fedex

import (
	"fmt"
)

// Track by Tracking number
func trackRequest(fedex Fedex, body string) string {
	return fedex.wrapSoapRequest(fmt.Sprintf(`
		<q0:TrackRequest>
			%s
			%s
			<q0:ProcessingOptions>INCLUDE_DETAILED_SCANS</q0:ProcessingOptions>
		</q0:TrackRequest>
	`, fedex.soapCreds(), body))
}

func soapNumberTracking(fedex Fedex, carrierCode string, trackingNo string) string {
	return trackRequest(fedex, fmt.Sprintf(`
		<q0:SelectionDetails>
			<q0:CarrierCode>%s</q0:CarrierCode>
			<q0:PackageIdentifier>
				<q0:Type>TRACKING_NUMBER_OR_DOORTAG</q0:Type>
				<q0:Value>%s</q0:Value>
			</q0:PackageIdentifier>
		</q0:SelectionDetails>
	`, carrierCode, trackingNo))
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
