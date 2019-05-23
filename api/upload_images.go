package api

import (
	"fmt"

	"github.com/happyreturns/fedex/models"
)

const (
	uploadVersion = "v11"
)

func (a API) UploadImages(images []models.Image) error {
	endpoint := fmt.Sprintf("/uploaddocument/%s", uploadVersion)
	request := a.uploadImagesRequest(images)
	response := &models.UploadImagesResponseEnvelope{}

	if err := a.makeRequestAndUnmarshalResponse(endpoint, request, response); err != nil {
		return fmt.Errorf("make upload images request and unmarshal: %s", err)
	}

	return nil
}

func (a API) uploadImagesRequest(images []models.Image) *models.Envelope {
	return &models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: fmt.Sprintf("http://fedex.com/ws/uploaddocument/%s", uploadVersion),
		Body: models.UploadImagesBody{
			UploadImagesRequest: models.UploadImagesRequest{
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
						ServiceID: "cdus",
						Major:     11,
					},
				},
				Images: images,
			},
		},
	}
}
