// History: Nov 20 13 tcolar Creation

package fedex

import (
	"fmt"

	"github.com/happyreturns/fedex/models"
)

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

func (f Fedex) trackByNumberSOAPRequest(carrierCode string, trackingNo string) models.Envelope {
	return models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: "http://fedex.com/ws/track/v16",
		Body: struct {
			TrackRequest models.TrackRequest `xml:"q0:TrackRequest"`
		}{
			TrackRequest: models.TrackRequest{
				Request: models.Request{
					WebAuthenticationDetail: models.WebAuthenticationDetail{
						UserCredential: models.UserCredential{
							Key:      f.Key,
							Password: f.Password,
						},
					},
					ClientDetail: models.ClientDetail{
						AccountNumber: f.Account,
						MeterNumber:   f.Meter,
					},
					Version: models.Version{
						ServiceID: "trck",
						Major:     16,
					},
				},
				ProcessingOptions: "INCLUDE_DETAILED_SCANS",
				SelectionDetails: models.SelectionDetails{
					CarrierCode: carrierCode,
					PackageIdentifier: models.PackageIdentifier{
						Type:  "TRACKING_NUMBER_OR_DOORTAG",
						Value: trackingNo,
					},
				},
			},
		},
	}
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
