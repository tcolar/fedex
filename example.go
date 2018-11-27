// History: Nov 20 13 tcolar Creation

package fedex

import (
	"log"

	"github.com/happyreturns/fedex/models"
)

// Examples
func main() {
	// You will need to fill in all those with your actual Fedex web service data
	fedex := Fedex{
		FedexURL: FedexAPITestURL, // OR FedexAPIURL
		Key:      "replaceWithYourKey",
		Password: "replaceWithYourPass",
		Account:  "replaceWithYourAccount",
		Meter:    "replaceWithYourMeter",
	}

	trackByReference(fedex, "replaceWithYourOrderRef", "replaceWithYourAccount")
	trackByPo(fedex, "123", "99032", "us")
	trackByNumber(fedex, "366849311565474")
}

// Looking up some tracking info by Fedex trackig number
func trackByNumber(fedex Fedex, trackingNo string) {
	reply, err := fedex.TrackByNumber("FDXE", trackingNo)
	if err != nil {
		log.Fatal(err)
	}
	Dump(*reply)
}

// Looking up some tracking info by reference
func trackByReference(fedex Fedex, ref string, account string) {
	reply, err := fedex.TrackByShipperRef("FDXE", ref, account)
	if err != nil {
		log.Fatal(err)
	}
	Dump(reply)
}

// Looking up some tracking info by Shipper PO number + Destination Zip
func trackByPo(fedex Fedex, po string, postalCode string, countryCode string) {
	reply, err := fedex.TrackByPo("FDXE", po, postalCode, countryCode)
	if err != nil {
		log.Fatal(err)
	}
	Dump(reply)
}

// Dump : Dumps some of the query resuts for testing
func Dump(reply models.TrackReply) {
	log.Print(reply)
	// Dummy example of using the data
	log.Printf("Successs : %t", !reply.Failed())
	if !reply.Failed() {
		tracking := reply.CompletedTrackDetails[0].TrackDetails[0].TrackingNumber
		log.Printf("Tracking Number: %s", tracking)
		log.Print(reply.CompletedTrackDetails[0].TrackDetails[0].ActualDeliveryAddress)
	}
}
