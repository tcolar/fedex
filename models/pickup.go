package models

type CreatePickupRequest struct {
	Request
	OriginDetail         OriginDetail        `xml:"q0:OriginDetail"`
	FreightPickupDetail  FreightPickupDetail `xml:"q0:FreightPickupDetail"`
	PackageCount         int                 `xml:"q0:PackageCount"`
	CarrierCode          string              `xml:"q0:CarrierCode"`
	Remarks              string              `xml:"q0:Remarks"`
	CommodityDescription string              `xml:"q0:CommodityDescription"`
}

type CreatePickupResponseEnvelope struct {
	Reply CreatePickupReply `xml:"Body>CreatePickupReply"`
}

func (c *CreatePickupResponseEnvelope) Error() error {
	return c.Reply.Error()
}

// CreatePickupReply : CreatePickup reply root (`xml:"Body>CreatePickupReply"`)
type CreatePickupReply struct {
	Reply
	PickupConfirmationNumber string
	Location                 string
}
