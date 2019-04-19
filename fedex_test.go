package fedex

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/happyreturns/fedex/models"
)

var (
	testFedex             Fedex
	prodFedex             Fedex
	laSmartPostFedex      Fedex
	blandonSmartPostFedex Fedex
)

func TestMain(m *testing.M) {
	credData, err := ioutil.ReadFile("creds.json")
	if err != nil {
		panic(err)
	}

	creds := map[string]Fedex{}
	if err := json.Unmarshal(credData, &creds); err != nil {
		panic(err)
	}

	testFedex = creds["test"]
	prodFedex = creds["prod"]
	laSmartPostFedex = creds["laSmartPost"]
	blandonSmartPostFedex = creds["blandonSmartPost"]

	os.Exit(m.Run())
}

func TestTrack(t *testing.T) {
	var (
		reply *models.TrackReply
		err   error
	)

	// Error case - invalid tracking
	_, err = testFedex.TrackByNumber(CarrierCodeExpress, "dkfjdkfj")
	checkErrorMatches(t, err, "make track request and unmarshal: response error: track detail error:")

	// Successful case
	reply, err = testFedex.TrackByNumber(CarrierCodeExpress, "123456789012")
	if err != nil {
		t.Fatal(err)
	}

	// Basic validation
	if reply.Error() != nil {
		t.Fatal("reply should not have failed")
	}
	if reply.HighestSeverity != "SUCCESS" ||
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
		reply.CompletedTrackDetails[0].TrackDetails[0].TrackingNumberUniqueIdentifier != "2458162000~123456789012~FX" ||
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

	// Successful case - smart post
	reply, err = testFedex.TrackByNumber(CarrierCodeSmartPost, "02396343485320033856")
	if err != nil {
		t.Fatal(err)
	}

	// Make sure smart post returned Events
	if len(reply.Events()) == 0 {
		t.Fatal("reply should have an event")
	}

	// Successful case - smart post
	reply, err = testFedex.TrackByNumber(CarrierCodeSmartPost, "02396343484520070272")
	if err != nil {
		t.Fatal(err)
	}

	// Make sure smart post returned Events
	if len(reply.Events()) == 0 {
		t.Fatal("reply should have an event")
	}
}

func TestRate(t *testing.T) {
	// Error case - invalid request
	_, err := prodFedex.Rate(models.Address{}, models.Address{}, models.Contact{}, models.Contact{})
	checkErrorMatches(t, err, "make rate request and unmarshal: response error: reply got error:")

	// Successful case
	reply, err := prodFedex.Rate(models.Address{
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

	// Basic validation
	if reply.Error() != nil {
		t.Fatal("reply should not have failed")
	}
	if reply.HighestSeverity != "SUCCESS" ||
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
	// Error case - nil shipment
	_, err := prodFedex.ShipGround(nil)
	checkErrorMatches(t, err, "create shipment request: empty shipment")

	// Error case - invalid shipment
	_, err = prodFedex.ShipGround(&Shipment{})
	checkErrorMatches(t, err, "make ship ground request and unmarshal: response error: reply got error:")

	// Successful case
	exampleShipment := &Shipment{
		FromAddress: models.Address{
			StreetLines:         []string{"1517 Lincoln Blvd"},
			City:                "Santa Monica",
			StateOrProvinceCode: "CA",
			PostalCode:          "90401",
			CountryCode:         "US",
		},
		ToAddress: models.Address{
			StreetLines:         []string{"1106 Broadway"},
			City:                "Santa Monica",
			StateOrProvinceCode: "CA",
			PostalCode:          "90401",
			CountryCode:         "US",
		},
		FromContact: models.Contact{
			PersonName:   "Jenny",
			PhoneNumber:  "213 867 5309",
			EmailAddress: "jenny@jenny.com",
		},
		ToContact: models.Contact{
			CompanyName:  "Some Company",
			PhoneNumber:  "214 867 5309",
			EmailAddress: "somecompany@somecompany.com",
		},
		NotificationEmail: "dev-notifications@happyreturns.com",
		Reference:         "My ship ground reference",
	}
	reply, err := prodFedex.ShipGround(exampleShipment)
	if err != nil {
		t.Fatal(err)
	}

	if reply.Error() != nil {
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
		reply.CompletedShipmentDetail.CompletedPackageDetails.Label.ImageType != "PNG" ||
		len(reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts) != 1 ||
		len(reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts[0].Image) == 0 {
		t.Fatal("output not correct")
	}

	// Decode png bytes from base64 data
	pngBytes, err := base64.StdEncoding.DecodeString(string(reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts[0].Image))
	if err != nil {
		t.Fatal(err)
	}

	// Write label as png, and manually check it
	err = ioutil.WriteFile(fmt.Sprintf("output-ground-%s.png", prodFedex.Key), pngBytes, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// it also works with no email
	exampleShipment.NotificationEmail = ""
	reply, err = prodFedex.ShipGround(exampleShipment)
	if err != nil {
		t.Fatal(err)
	}

	if reply.Error() != nil {
		fmt.Println(reply)
		t.Fatal("reply should not have failed")
	}
}

func TestShipSmartPost(t *testing.T) {
	// Error case - nil shipment
	_, err := laSmartPostFedex.ShipSmartPost(nil)
	if err == nil || !strings.HasPrefix(err.Error(), "create shipment request: empty shipment") {
		t.Fatal("error did not match", err)
	}

	// Error case - invalid shipment
	_, err = laSmartPostFedex.ShipSmartPost(&Shipment{})
	if err == nil || !strings.HasPrefix(err.Error(), "make ship smart post request and unmarshal: response error: reply got error:") {
		t.Fatal("error did not match", err)
	}

	// Successful cases
	testShipSmartPostSuccess(t, laSmartPostFedex)
	testShipSmartPostSuccess(t, blandonSmartPostFedex)
}

func TestCreatePickup(t *testing.T) {
	reply, err := prodFedex.CreatePickup(models.PickupLocation{
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

	if reply.Error() != nil {
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

func TestSendNotifications(t *testing.T) {
	// Error case - invalid tracking number
	_, err := prodFedex.SendNotifications("123", "dev-notifications@happyreturns.com")
	checkErrorMatches(t, err, "make send notifications request: response error: reply got error: Invalid tracking numbers.")

	// Successful case
	reply, err := prodFedex.SendNotifications("02396343485320152281", "dev-notifications@happyreturns.com")
	if err != nil {
		t.Fatal("error did not match", err)
	}
	if reply.Error() != nil {
		t.Fatal("reply should not have failed")
	}

	// Basic validation
	if reply.HighestSeverity != "SUCCESS" ||
		len(reply.Notifications) != 1 ||
		reply.Notifications[0].Source != "trck" ||
		reply.Notifications[0].Code != "0" ||
		reply.Notifications[0].Message != "Request was successfully processed." ||
		reply.Notifications[0].LocalizedMessage != "Request was successfully processed." ||
		reply.DuplicateWaybill ||
		reply.MoreDataAvailable ||
		len(reply.Packages) != 1 ||
		reply.Packages[0].TrackingNumber != "02396343485320152281" ||
		len(reply.Packages[0].TrackingNumberUniqueIdentifiers) != 1 ||
		reply.Packages[0].CarrierCode != "FXSP" ||
		reply.Packages[0].ShipDate != "2019-03-11" ||
		reply.Packages[0].Destination.City != "BLANDON" ||
		reply.Packages[0].Destination.StateOrProvinceCode != "PA" ||
		reply.Packages[0].Destination.CountryCode != "US" ||
		bool(reply.Packages[0].Destination.Residential) ||
		len(reply.Packages[0].RecipientDetails) != 1 ||
		len(reply.Packages[0].RecipientDetails[0].NotificationEventsAvailable) != 1 {
		t.Fatal("output not correct")
	}
}

func testShipSmartPostSuccess(t *testing.T, fedexAccount Fedex) {
	exampleShipment := &Shipment{
		FromAddress: models.Address{
			StreetLines:         []string{"1517 Lincoln Blvd"},
			City:                "Santa Monica",
			StateOrProvinceCode: "CA",
			PostalCode:          "90401",
			CountryCode:         "US",
		},
		ToAddress: models.Address{},
		FromContact: models.Contact{
			PersonName:   "Jenny",
			PhoneNumber:  "213 867 5309",
			EmailAddress: "jenny@jenny.com",
		},
		ToContact: models.Contact{
			CompanyName:  "Some Company",
			PhoneNumber:  "214 867 5309",
			EmailAddress: "somecompany@somecompany.com",
		},
		NotificationEmail: "dev-notifications@happyreturns.com",
		Reference:         "My reference",
	}
	reply, err := fedexAccount.ShipSmartPost(exampleShipment)
	if err != nil {
		t.Fatal(err)
	}

	// Basic validation
	if reply.Error() != nil {
		fmt.Println(reply)
		t.Fatal("reply should not have failed")
	}
	if reply.HighestSeverity != "SUCCESS" ||
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
		reply.CompletedShipmentDetail.CompletedPackageDetails.Label.ImageType != "PNG" ||
		len(reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts) != 1 ||
		len(reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts[0].Image) == 0 {
		t.Fatal("output not correct")
	}

	// Decode pdf bytes from base64 data
	pngBytes, err := base64.StdEncoding.DecodeString(string(reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts[0].Image))
	if err != nil {
		t.Fatal(err)
	}

	// Write label as png, and manually check it
	err = ioutil.WriteFile(fmt.Sprintf("output-smart-post-%s.png", fedexAccount.Key), pngBytes, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// It also works with no email
	exampleShipment.NotificationEmail = ""
	reply, err = fedexAccount.ShipSmartPost(exampleShipment)
	if err != nil {
		t.Fatal(err)
	}

	if reply.Error() != nil {
		fmt.Println(reply)
		t.Fatal("reply should not have failed")
	}
}

func checkErrorMatches(t *testing.T, err error, expectedText string) {
	if err == nil || !strings.HasPrefix(err.Error(), expectedText) {
		t.Fatal("error", err, "doesn't match", expectedText)
	}
}
