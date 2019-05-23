package api

import (
	"fmt"

	"github.com/happyreturns/fedex/models"
)

func (a API) TrackByNumber(carrierCode, trackingNo string) (*models.TrackReply, error) {
	request := a.trackByNumberRequest(carrierCode, trackingNo)
	response := &models.TrackResponseEnvelope{}

	err := a.makeRequestAndUnmarshalResponse("/trck", request, response)
	if err != nil {
		return nil, fmt.Errorf("make track request and unmarshal: %s", err)
	}
	return &response.Reply, nil
}

func (a API) trackByNumberRequest(carrierCode string, trackingNo string) *models.Envelope {
	return &models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: "http://fedex.com/ws/track/v16",
		Body: models.TrackBody{
			TrackRequest: models.TrackRequest{
				Request: models.Request{
					WebAuthenticationDetail: models.WebAuthenticationDetail{
						UserCredential: models.UserCredential{
							Key:      a.Key,
							Password: a.Password,
						},
					},
					ClientDetail: models.ClientDetail{
						AccountNumber: a.Account,
						MeterNumber:   a.Meter,
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
