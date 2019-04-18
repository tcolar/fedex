package models

type ProcessShipmentRequest struct {
	Request
	RequestedShipment RequestedShipment `xml:"q0:RequestedShipment"`
}

type ShipResponseEnvelope struct {
	Reply ProcessShipmentReply `xml:"Body>ProcessShipmentReply"`
}

func (s *ShipResponseEnvelope) Error() error {
	return s.Reply.Error()
}

// ProcessShipReply : Process shipment reply root (`xml:"Body>ProcessShipmentReply"`)
type ProcessShipmentReply struct {
	Reply
	TransactionDetail       TransactionDetail
	CompletedShipmentDetail CompletedShipmentDetail
	Events                  []Event
}
