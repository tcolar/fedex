package fedex

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/happyreturns/fedex/models"
)

var f Fedex = Fedex{
	// TODO fill in your creds
	// Key:      "",
	// Password: "",
	// Account:  "",
	// Meter:    "",
	// HubID:    "", // for SmartPost
}

func TestTrack(t *testing.T) {
	if f.HubID != "" || f.FedexURL != FedexAPITestURL {
		t.SkipNow()
	}
	reply, err := f.TrackByNumber(CarrierCodeExpress, "123456789012")
	if err != nil {
		t.Fatal(err)
	}
	if reply.Failed() {
		t.Fatal("reply should not have failed")
	}
	if reply.HighestSeverity != "SUCCESS" ||
		// Basic validation
		len(reply.Notifications) != 1 ||
		reply.Notifications[0].Source != "trck" ||
		reply.Notifications[0].Code != "0" ||
		reply.Notifications[0].Message != "Request was successfully processed." ||
		reply.Notifications[0].LocalizedMessage != "Request was successfully processed." ||
		reply.Version.ServiceID != "trck" ||
		reply.Version.Major != 16 ||
		reply.Version.Intermediate != 0 ||
		reply.Version.Minor != 0 ||
		len(reply.CompletedTrackDetails) != 1 ||
		!reply.CompletedTrackDetails[0].DuplicateWaybill ||
		reply.CompletedTrackDetails[0].MoreData ||
		len(reply.CompletedTrackDetails[0].TrackDetails) < 1 ||
		reply.CompletedTrackDetails[0].TrackDetails[0].OperatingCompanyOrCarrierDescription != "FedEx Express" ||
		reply.CompletedTrackDetails[0].TrackDetails[0].TrackingNumber != "123456789012" ||
		reply.CompletedTrackDetails[0].TrackDetails[0].TrackingNumberUniqueIdentifier != "2458115001~123456789012~FX" ||
		reply.CompletedTrackDetails[0].TrackDetails[0].CarrierCode != "FDXE" {
		t.Fatal("output not correct")
	}

	// Check the timestamps make sense
	if len(reply.CompletedTrackDetails[0].TrackDetails[0].DatesOrTimes) == 0 {
		t.Fatal("should get at least one dateortime")
	}
	if reply.ActualDelivery().IsZero() {
		t.Fatal("actual delivery should be set")
	}
	if reply.EstimatedDelivery().IsZero() {
		t.Fatal("actual delivery should be set")
	}
}

func TestRate(t *testing.T) {
	reply, err := f.Rate(models.Address{
		StreetLines:         []string{"1517 Lincoln Blvd"},
		City:                "Santa Monica",
		StateOrProvinceCode: "CA",
		PostalCode:          "90401",
		CountryCode:         "US",
	}, models.Address{
		StreetLines:         []string{"1106 Broadway"},
		City:                "Santa Monica",
		StateOrProvinceCode: "CA",
		PostalCode:          "90401",
		CountryCode:         "US",
	}, models.Contact{
		PersonName:   "Jenny",
		PhoneNumber:  "213 867 5309",
		EmailAddress: "jenny@jenny.com",
	}, models.Contact{
		CompanyName:  "Some Company",
		PhoneNumber:  "214 867 5309",
		EmailAddress: "somecompany@somecompany.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	if reply.Failed() {
		t.Fatal("reply should not have failed")
	}
	if reply.HighestSeverity != "SUCCESS" ||
		// Basic validation
		len(reply.Notifications) != 1 ||
		reply.Notifications[0].Source != "crs" ||
		reply.Notifications[0].Code != "0" ||
		reply.Version.ServiceID != "crs" ||
		reply.Version.Major != 24 ||
		reply.Version.Intermediate != 0 ||
		reply.Version.Minor != 0 ||
		len(reply.RateReplyDetails) != 1 ||
		reply.RateReplyDetails[0].ServiceType != "FEDEX_GROUND" ||
		reply.RateReplyDetails[0].ServiceDescription.ServiceType != "FEDEX_GROUND" ||
		reply.RateReplyDetails[0].ServiceDescription.Code != "92" ||
		reply.RateReplyDetails[0].ServiceDescription.AstraDescription != "FXG" ||
		reply.RateReplyDetails[0].PackagingType != "YOUR_PACKAGING" ||
		reply.RateReplyDetails[0].DestinationAirportID != "LAX" ||
		reply.RateReplyDetails[0].IneligibleForMoneyBackGuarantee ||
		reply.RateReplyDetails[0].SignatureOption != "SERVICE_DEFAULT" ||
		reply.RateReplyDetails[0].ActualRateType != "PAYOR_ACCOUNT_PACKAGE" ||
		len(reply.RateReplyDetails[0].RatedShipmentDetails) != 2 ||
		reply.RateReplyDetails[0].RatedShipmentDetails[0].EffectiveNetDiscount.Amount == "" ||
		len(reply.RateReplyDetails[0].RatedShipmentDetails[0].RatedPackages) != 1 ||
		reply.RateReplyDetails[0].RatedShipmentDetails[0].RatedPackages[0].PackageRateDetail.NetCharge.Amount == "0.0" ||
		len(reply.RateReplyDetails[0].RatedShipmentDetails[1].RatedPackages) != 1 ||
		reply.RateReplyDetails[0].RatedShipmentDetails[1].RatedPackages[0].PackageRateDetail.NetCharge.Amount == "0.0" {
		t.Fatal("output not correct")
	}

	charge, err := reply.TotalCost()
	if err != nil {
		t.Fatal(err)
	}

	if charge.Currency != "USD" || charge.Amount == "0.00" {
		t.Fatal("totalCost should be non-zero, USD")
	}
}

func TestShipGround(t *testing.T) {
	if f.HubID != "" {
		t.SkipNow()
	}
	reply, err := f.ShipGround(models.Address{
		StreetLines:         []string{"1517 Lincoln Blvd"},
		City:                "Santa Monica",
		StateOrProvinceCode: "CA",
		PostalCode:          "90401",
		CountryCode:         "US",
	}, models.Address{
		StreetLines:         []string{"1106 Broadway"},
		City:                "Santa Monica",
		StateOrProvinceCode: "CA",
		PostalCode:          "90401",
		CountryCode:         "US",
	}, models.Contact{
		PersonName:   "Jenny",
		PhoneNumber:  "213 867 5309",
		EmailAddress: "jenny@jenny.com",
	}, models.Contact{
		CompanyName:  "Some Company",
		PhoneNumber:  "214 867 5309",
		EmailAddress: "somecompany@somecompany.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	if reply.Failed() {
		t.Fatal("reply should not have failed")
	}
	if reply.HighestSeverity != "SUCCESS" ||
		// Basic validation
		len(reply.Notifications) != 1 ||
		reply.Notifications[0].Source != "ship" ||
		reply.Notifications[0].Code != "0000" ||
		reply.Notifications[0].Message != "Success" ||
		reply.Notifications[0].LocalizedMessage != "Success" ||
		reply.Version.ServiceID != "ship" ||
		reply.Version.Major != 23 ||
		reply.Version.Intermediate != 0 ||
		reply.Version.Minor != 0 ||
		reply.JobID == "" ||
		reply.CompletedShipmentDetail.UsDomestic != "true" ||
		reply.CompletedShipmentDetail.CarrierCode != "FDXG" ||
		reply.CompletedShipmentDetail.MasterTrackingId.TrackingIdType != "FEDEX" ||
		reply.CompletedShipmentDetail.MasterTrackingId.TrackingNumber == "" ||
		reply.CompletedShipmentDetail.ServiceTypeDescription != "FXG" ||
		reply.CompletedShipmentDetail.ServiceDescription.ServiceType != "FEDEX_GROUND" ||
		reply.CompletedShipmentDetail.ServiceDescription.Code != "92" ||
		// skip ServiceDescription.Names
		reply.CompletedShipmentDetail.PackagingDescription != "YOUR_PACKAGING" ||
		reply.CompletedShipmentDetail.OperationalDetail.OriginLocationNumber != "901" ||
		reply.CompletedShipmentDetail.OperationalDetail.DestinationLocationNumber != "901" ||
		reply.CompletedShipmentDetail.OperationalDetail.TransitTime != "ONE_DAY" ||
		reply.CompletedShipmentDetail.OperationalDetail.IneligibleForMoneyBackGuarantee != "false" ||
		reply.CompletedShipmentDetail.OperationalDetail.DeliveryEligibilities != "SATURDAY_DELIVERY" ||
		reply.CompletedShipmentDetail.OperationalDetail.ServiceCode != "92" ||
		reply.CompletedShipmentDetail.OperationalDetail.PackagingCode != "01" ||
		reply.CompletedShipmentDetail.ShipmentRating.ActualRateType != "PAYOR_ACCOUNT_PACKAGE" ||
		reply.CompletedShipmentDetail.ShipmentRating.EffectiveNetDiscount.Currency != "USD" ||
		len(reply.CompletedShipmentDetail.ShipmentRating.ShipmentRateDetails) != 2 ||
		// // skip most ShipmentRateDetails fields
		reply.CompletedShipmentDetail.ShipmentRating.ShipmentRateDetails[0].RateType != "PAYOR_ACCOUNT_PACKAGE" ||
		reply.CompletedShipmentDetail.ShipmentRating.ShipmentRateDetails[1].RateType != "PAYOR_LIST_PACKAGE" ||
		len(reply.CompletedShipmentDetail.CompletedPackageDetails.TrackingIds) != 1 ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.TrackingIds[0].TrackingIdType != "FEDEX" ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Type != "OUTBOUND_LABEL" ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.Label.ImageType != "PDF" ||
		len(reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts) != 1 ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts[0].Image == "" {
		t.Fatal("output not correct")
	}
}

func TestShipSmartPost(t *testing.T) {
	if f.HubID == "" {
		t.SkipNow()
	}

	reply, err := f.ShipSmartPost(
		models.Address{
			StreetLines:         []string{"1517 Lincoln Blvd"},
			City:                "Santa Monica",
			StateOrProvinceCode: "CA",
			PostalCode:          "90401",
			CountryCode:         "US",
		},
		models.Address{},
		models.Contact{
			PersonName:   "Jenny",
			PhoneNumber:  "213 867 5309",
			EmailAddress: "jenny@jenny.com",
		}, models.Contact{
			CompanyName:  "Some Company",
			PhoneNumber:  "214 867 5309",
			EmailAddress: "somecompany@somecompany.com",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if reply.Failed() {
		fmt.Println(reply)
		t.Fatal("reply should not have failed")
	}
	if reply.HighestSeverity != "SUCCESS" ||
		// Basic validation
		len(reply.Notifications) != 1 ||
		reply.Notifications[0].Source != "ship" ||
		reply.Notifications[0].Code != "0000" ||
		reply.Notifications[0].Message != "Success" ||
		reply.Notifications[0].LocalizedMessage != "Success" ||
		reply.Version.ServiceID != "ship" ||
		reply.Version.Major != 23 ||
		reply.Version.Intermediate != 0 ||
		reply.Version.Minor != 0 ||
		reply.JobID == "" ||
		reply.CompletedShipmentDetail.UsDomestic != "true" ||
		reply.CompletedShipmentDetail.CarrierCode != "FXSP" ||
		reply.CompletedShipmentDetail.MasterTrackingId.TrackingIdType != "USPS" ||
		reply.CompletedShipmentDetail.MasterTrackingId.TrackingNumber == "" ||
		reply.CompletedShipmentDetail.ServiceTypeDescription != "SMART POST" ||
		reply.CompletedShipmentDetail.ServiceDescription.ServiceType != "SMART_POST" ||
		reply.CompletedShipmentDetail.PackagingDescription != "YOUR_PACKAGING" ||
		reply.CompletedShipmentDetail.OperationalDetail.TransitTime != "TWO_DAYS" ||
		reply.CompletedShipmentDetail.OperationalDetail.IneligibleForMoneyBackGuarantee != "false" ||
		len(reply.CompletedShipmentDetail.CompletedPackageDetails.TrackingIds) != 1 ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.TrackingIds[0].TrackingIdType != "USPS" ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Type != "OUTBOUND_LABEL" ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.Label.ImageType != "PDF" ||
		len(reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts) != 1 ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts[0].Image == "" {
		t.Fatal("output not correct")
	}

	// Decode pdf bytes from base64 data
	pdfBytes, err := base64.StdEncoding.DecodeString(reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts[0].Image)
	if err != nil {
		t.Fatal(err)
	}

	// Write label as pdf, and manually check it
	err = ioutil.WriteFile("output-smart-post.pdf", pdfBytes, 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreatePickup(t *testing.T) {
	// Fill in prod creds to run this test, as this test only works in prod
	t.SkipNow()

	reply, err := f.CreatePickup(models.PickupLocation{
		Address: models.Address{
			StreetLines:         []string{"1517 Lincoln Blvd"},
			City:                "Santa Monica",
			StateOrProvinceCode: "CA",
			PostalCode:          "90401",
			CountryCode:         "US",
		},
		Contact: models.Contact{
			PersonName:   "Jenny",
			PhoneNumber:  "213 867 5309",
			EmailAddress: "jenny@jenny.com",
		}},
		models.Address{
			StreetLines:         []string{"1106 Broadway"},
			City:                "Santa Monica",
			StateOrProvinceCode: "CA",
			PostalCode:          "90401",
			CountryCode:         "US",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if reply.Failed() {
		t.Fatal("reply should not have failed")
	}

	if reply.HighestSeverity != "SUCCESS" ||
		// Basic validation
		len(reply.Notifications) != 1 ||
		reply.Notifications[0].Source != "disp" ||
		reply.Notifications[0].Code != "0000" ||
		reply.Notifications[0].Message != "Success" ||
		reply.Version.ServiceID != "disp" ||
		reply.Version.Major != 17 ||
		reply.Version.Intermediate != 0 ||
		reply.Version.Minor != 0 ||
		reply.PickupConfirmationNumber == "" ||
		reply.PickupConfirmationNumber == "0" ||
		reply.Location != "SMOA" {
		t.Fatal("output not correct")
	}
}
