package models

// SendNotificationsRequest
type SendNotificationsRequest struct {
	Request
	TrackingNumber         string `xml:"q0:TrackingNumber"`
	TrackingNumberUniqueID string `xml:"q0:TrackingNumberUniqueId"`
	// ShipDateRangeBegin      DateOrTimestamp          `xml:"q0:ShipDateRangeBegin"` // Don't bother with these for now
	// ShipDateRangeEnd        DateOrTimestamp          `xml:"q0:ShipDateRangeEnd"`
	SenderEmailAddress      string                  `xml:"q0:SenderEMailAddress"`
	SenderContactName       string                  `xml:"q0:SenderContactName"`
	EventNotificationDetail EventNotificationDetail `xml:"q0:EventNotificationDetail"`
}

type SendNotificationsResponseEnvelope struct {
	Reply SendNotificationsReply `xml:"Body>SendNotificationsReply"`
}

func (s *SendNotificationsResponseEnvelope) Error() error {
	return s.Reply.Error()
}

// SendNotificationsReply : CreatePickup reply root (`xml:"Body>SendNotificationsReply"`)
type SendNotificationsReply struct {
	Reply
	DuplicateWaybill  bool
	MoreDataAvailable bool
	PagingToken       string
	Packages          []Package
}
