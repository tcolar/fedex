package api

import (
	"fmt"

	"github.com/happyreturns/fedex/models"
)

const (
	sendNotificationsVersion = "v16"
)

// SendNotifications gets notifications sent to an email
func (a API) SendNotifications(trackingNo, email string) (*models.SendNotificationsReply, error) {

	endpoint := fmt.Sprintf("/track/%s", sendNotificationsVersion)
	request := a.sendNotificationsRequest(trackingNo, email)
	response := &models.SendNotificationsResponseEnvelope{}

	err := a.makeRequestAndUnmarshalResponse(endpoint, request, response)
	if err != nil {
		return nil, fmt.Errorf("make send notifications request: %s", err)
	}
	return &response.Reply, nil
}

func (a API) sendNotificationsRequest(trackingNo, email string) *models.Envelope {
	return &models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: fmt.Sprintf("http://fedex.com/ws/track/%s", sendNotificationsVersion),
		Body: models.SendNotificationsBody{
			SendNotificationsRequest: models.SendNotificationsRequest{
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
				TrackingNumber:     trackingNo,
				SenderEmailAddress: email,
				SenderContactName:  "Customer",
				EventNotificationDetail: models.EventNotificationDetail{
					AggregationType: "PER_PACKAGE",
					PersonalMessage: "Message",
					EventNotifications: []models.EventNotification{{
						Role: "SHIPPER",
						Events: []string{
							"ON_DELIVERY",
							"ON_ESTIMATED_DELIVERY",
							"ON_EXCEPTION",
							"ON_SHIPMENT",
							"ON_TENDER",
						},
						NotificationDetail: models.NotificationDetail{
							NotificationType: "EMAIL",
							EmailDetail: models.EmailDetail{
								EmailAddress: "joachim@happyreturns.com",
								Name:         "joachim@happyreturns.com",
							},
							Localization: models.Localization{
								LanguageCode: "en",
							},
						},
						FormatSpecification: models.FormatSpecification{
							Type: "HTML",
						},
					}},
				},
			},
		},
	}
}
