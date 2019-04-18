package models

import "errors"

type RateRequest struct {
	Request
	RequestedShipment RequestedShipment `xml:"q0:RequestedShipment"`
}

type RateResponseEnvelope struct {
	Reply RateReply `xml:"Body>RateReply"`
}

func (r *RateResponseEnvelope) Error() error {
	return r.Reply.Error()
}

// RateReply : Process shipment reply root (`xml:"Body>RateReply"`)
type RateReply struct {
	Reply
	TransactionDetail TransactionDetail
	RateReplyDetails  []RateReplyDetail
}

// TotalCost returns the first TotalNetChargeWithDutiesAndTaxes in the reply
func (rr *RateReply) TotalCost() (Charge, error) {
	// TotalNetChargeWithDutiesAndTaxes
	for _, rateReplyDetail := range rr.RateReplyDetails {
		for _, ratedShipmentDetail := range rateReplyDetail.RatedShipmentDetails {
			totalNetCharge := ratedShipmentDetail.ShipmentRateDetail.TotalNetChargeWithDutiesAndTaxes
			if totalNetCharge.Currency != "" && totalNetCharge.Amount != "" {
				return totalNetCharge, nil
			}
		}
	}
	return Charge{}, errors.New("no total net charge found on reply")
}
