// History: Nov 20 13 tcolar Creation

package fedex

import (
	"log"
)

// Just some example
func main() {
	// You will need to fill in all those with your actual Fedex web service data
	fedex := Fedex{
		FedexUrl: FEDEX_API_TEST_URL, // OR FEDEX_API_URL
		Key:      "replaceWithYourKey",
		Password: "replaceWithYourPass",
		Account:  "replaceWithYourAccount",
		Meter:    "replaceWithYourMeter",
	}

	trackByReference(fedex, "replaceWithYourOrderRef", "replaceWithYourAccount")
	trackByPo(fedex, "123", "99032", "us")
}

//Looking up some tracking info by reference
func trackByReference(fedex Fedex, ref string, account string) {
	reply, err := fedex.TrackByShipperRef("FDXE", ref, account)
	if err != nil {
		log.Fatal(err)
	}
	dump(reply)
}

//Looking up some tracking info by reference
func trackByPo(fedex Fedex, po string, postalCode string, countryCode string) {
	reply, err := fedex.TrackByPo("FDXE", po, postalCode, countryCode)
	if err != nil {
		log.Fatal(err)
	}
	dump(reply)
}

// Dump some of the query resuts as an example
func dump(reply TrackReply) {
	log.Print(reply)
	// Dummy example of using the data
	log.Printf("Successs : %t", !reply.Failed())
	if !reply.Failed() {
		tracking := reply.CompletedTrackDetails[0].TrackDetails[0].TrackingNumber
		log.Printf("Tracking Number: %s", tracking)
		log.Print(reply.CompletedTrackDetails[0].TrackDetails[0].ActualDeliveryAddress)
	}
}
