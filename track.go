// History: Nov 20 13 tcolar Creation

package fedex

import (
	"github.com/happyreturns/fedex/models"
)

func (f Fedex) trackByNumberRequest(carrierCode string, trackingNo string) models.Envelope {
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
