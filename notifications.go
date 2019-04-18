package fedex

import "github.com/happyreturns/fedex/models"

func (f Fedex) notificationsRequest(trackingNo, email string) models.Envelope {
	return models.Envelope{
		Soapenv:   "http://schemas.xmlsoap.org/soap/envelope/",
		Namespace: "http://fedex.com/ws/track/v16",
		Body: struct {
			SendNotificationsRequest models.SendNotificationsRequest `xml:"q0:SendNotificationsRequest"`
		}{
			SendNotificationsRequest: models.SendNotificationsRequest{
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
